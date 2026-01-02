package access

import (
	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/prayaspoudel/modules/access/delivery/http"
	"github.com/prayaspoudel/modules/access/delivery/http/route"
	"github.com/prayaspoudel/modules/access/features/auth"
	"github.com/prayaspoudel/modules/access/middleware"
	"github.com/prayaspoudel/modules/access/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func Bootstrap(config *BootstrapConfig) {
	// Setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	sessionRepository := repository.NewSessionRepository(config.Log)
	tokenRepository := repository.NewRefreshTokenRepository(config.Log)
	companyRepository := repository.NewCompanyRepository(config.Log)

	// Setup use cases
	authUseCase := auth.NewAuthUseCase(
		config.DB,
		config.Log,
		config.Config,
		userRepository,
		sessionRepository,
		tokenRepository,
		companyRepository,
	)

	// Setup controllers
	authController := http.NewAuthController(config.Log, authUseCase, config.Validate)

	// Setup middleware
	authMiddleware := middleware.NewAuthMiddleware(authUseCase)

	// Setup routes
	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}
