package model

// UserResponse represents user data in API responses
type UserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	CompanyID     string `json:"companyId,omitempty"`
	Role          string `json:"role"`
	IsActive      bool   `json:"isActive"`
	IsVerified    bool   `json:"isVerified"`
	EmailVerified bool   `json:"emailVerified"`
	CreatedAt     int64  `json:"createdAt"`
	UpdatedAt     int64  `json:"updatedAt"`
}

// RegisterUserRequest represents a registration request
type RegisterUserRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
	FirstName string `json:"firstName" validate:"required,max=100"`
	LastName  string `json:"lastName" validate:"required,max=100"`
}

// LoginUserRequest represents a login request
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	ClientID string `json:"clientId,omitempty"`
}

// LoginResponse represents the authentication response
type LoginResponse struct {
	AccessToken  string            `json:"accessToken"`
	RefreshToken string            `json:"refreshToken"`
	ExpiresIn    int               `json:"expiresIn"`
	TokenType    string            `json:"tokenType"`
	User         *UserResponse     `json:"user"`
	Companies    []CompanyResponse `json:"companies,omitempty"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

// ChangePasswordRequest represents a password change request
type ChangePasswordRequest struct {
	UserID      string `json:"-" validate:"required"`
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// ForgotPasswordRequest represents a forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest represents a password reset request
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

// VerifyEmailRequest represents an email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// GetUserRequest represents a get user request
type GetUserRequest struct {
	UserID string `json:"-" validate:"required"`
}

// UpdateUserRequest represents an update user request
type UpdateUserRequest struct {
	UserID    string `json:"-" validate:"required"`
	FirstName string `json:"firstName" validate:"max=100"`
	LastName  string `json:"lastName" validate:"max=100"`
	Email     string `json:"email" validate:"email,max=255"`
}

// LogoutRequest represents a logout request
type LogoutRequest struct {
	UserID string `json:"-" validate:"required"`
}
