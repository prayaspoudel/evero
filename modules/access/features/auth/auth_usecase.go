package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/prayaspoudel/modules/access/entity"
	"github.com/prayaspoudel/modules/access/model"
	"github.com/prayaspoudel/modules/access/model/converter"
	"github.com/prayaspoudel/modules/access/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Viper             *viper.Viper
	UserRepository    *repository.UserRepository
	SessionRepository *repository.SessionRepository
	RefreshTokenRepo  *repository.RefreshTokenRepository
	CompanyRepository *repository.CompanyRepository
}

func NewAuthUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	viper *viper.Viper,
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	companyRepo *repository.CompanyRepository,
) *AuthUseCase {
	return &AuthUseCase{
		DB:                db,
		Log:               log,
		Viper:             viper,
		UserRepository:    userRepo,
		SessionRepository: sessionRepo,
		RefreshTokenRepo:  refreshTokenRepo,
		CompanyRepository: companyRepo,
	}
}

func (uc *AuthUseCase) Register(req *model.RegisterUserRequest) (*model.UserResponse, error) {
	// Check if user exists
	var existingUser entity.User
	err := uc.UserRepository.FindByEmail(uc.DB, &existingUser, req.Email)
	if err == nil {
		return nil, fiber.NewError(fiber.StatusConflict, "email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.Log.WithError(err).Error("error checking existing user")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.Log.WithError(err).Error("error hashing password")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Create user
	user := &entity.User{
		ID:            uuid.New().String(),
		Email:         req.Email,
		PasswordHash:  string(hashedPassword),
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		IsActive:      true,
		IsVerified:    false,
		EmailVerified: false,
	}

	if err := uc.UserRepository.Create(uc.DB, user); err != nil {
		uc.Log.WithError(err).Error("error creating user")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return converter.UserToResponse(user), nil
}

func (uc *AuthUseCase) Login(req *model.LoginUserRequest, ipAddress string) (*model.LoginResponse, error) {
	// Find user
	var user entity.User
	err := uc.UserRepository.FindByEmail(uc.DB, &user, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
		}
		uc.Log.WithError(err).Error("error finding user")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fiber.NewError(fiber.StatusForbidden, "account is inactive")
	}

	// Generate access token
	accessToken, expiresIn, err := uc.generateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenStr := uuid.New().String()
	refreshToken := &entity.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // 30 days
	}

	if err := uc.RefreshTokenRepo.Create(uc.DB, refreshToken); err != nil {
		uc.Log.WithError(err).Error("error creating refresh token")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Create session
	session := &entity.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		SessionToken: accessToken,
		IPAddress:    ipAddress,
		UserAgent:    "", // TODO: Get from request
		ExpiresAt:    time.Now().Add(time.Duration(expiresIn) * time.Second),
	}

	if err := uc.SessionRepository.Create(uc.DB, session); err != nil {
		uc.Log.WithError(err).Error("error creating session")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Update last login
	if err := uc.UserRepository.UpdateLastLogin(uc.DB, user.ID, ipAddress); err != nil {
		uc.Log.WithError(err).Error("error updating last login")
	}

	// Get user companies
	var companies []entity.Company
	if err := uc.CompanyRepository.FindByUserID(uc.DB, &companies, user.ID); err != nil {
		uc.Log.WithError(err).Error("error fetching companies")
		// Don't fail the login, just return empty companies
	}

	companyResponses := make([]model.CompanyResponse, len(companies))
	for i, company := range companies {
		companyResponses[i] = *converter.CompanyToResponse(&company)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		User:         converter.UserToResponse(&user),
		Companies:    companyResponses,
	}, nil
}

func (uc *AuthUseCase) Logout(userID string, token string) error {
	// Delete session
	if err := uc.SessionRepository.DeleteByToken(uc.DB, token); err != nil {
		uc.Log.WithError(err).Error("error deleting session")
	}

	// Revoke all refresh tokens for user
	if err := uc.RefreshTokenRepo.RevokeByUserID(uc.DB, userID); err != nil {
		uc.Log.WithError(err).Error("error revoking refresh tokens")
		return fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return nil
}

func (uc *AuthUseCase) RefreshToken(req *model.RefreshTokenRequest) (*model.LoginResponse, error) {
	// Find refresh token
	var refreshToken entity.RefreshToken
	err := uc.RefreshTokenRepo.FindByToken(uc.DB, &refreshToken, req.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid refresh token")
		}
		uc.Log.WithError(err).Error("error finding refresh token")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Check if token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "refresh token expired")
	}

	// Get user
	var user entity.User
	if err := uc.UserRepository.FindByID(uc.DB, &user, refreshToken.UserID); err != nil {
		uc.Log.WithError(err).Error("error finding user")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Generate new access token
	accessToken, expiresIn, err := uc.generateAccessToken(&user)
	if err != nil {
		return nil, err
	}

	// Rotate refresh token
	newRefreshTokenStr := uuid.New().String()
	newRefreshToken := &entity.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     newRefreshTokenStr,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}

	if err := uc.RefreshTokenRepo.Create(uc.DB, newRefreshToken); err != nil {
		uc.Log.WithError(err).Error("error creating refresh token")
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	// Revoke old refresh token
	if err := uc.RefreshTokenRepo.RevokeByToken(uc.DB, req.RefreshToken); err != nil {
		uc.Log.WithError(err).Error("error revoking old refresh token")
	}

	// Get user companies
	var companies []entity.Company
	if err := uc.CompanyRepository.FindByUserID(uc.DB, &companies, user.ID); err != nil {
		uc.Log.WithError(err).Error("error fetching companies")
	}

	companyResponses := make([]model.CompanyResponse, len(companies))
	for i, company := range companies {
		companyResponses[i] = *converter.CompanyToResponse(&company)
	}

	return &model.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		User:         converter.UserToResponse(&user),
		Companies:    companyResponses,
	}, nil
}

func (uc *AuthUseCase) generateAccessToken(user *entity.User) (string, int, error) {
	expiresIn := uc.Viper.GetInt("jwt.expiration")
	if expiresIn == 0 {
		expiresIn = 3600 // Default 1 hour
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := uc.Viper.GetString("jwt.secret")
	if secret == "" {
		uc.Log.Error("JWT secret not configured")
		return "", 0, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		uc.Log.WithError(err).Error("error signing token")
		return "", 0, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	return tokenString, expiresIn, nil
}

func (uc *AuthUseCase) VerifyAccessToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
		}
		return []byte(uc.Viper.GetString("jwt.secret")), nil
	})

	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
}
