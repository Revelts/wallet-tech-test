package main

import (
	"log"
	"wallet-service/internal/handler"
	"wallet-service/internal/infrastructure/config"
	"wallet-service/internal/infrastructure/database"
	"wallet-service/internal/infrastructure/jwt"
	"wallet-service/internal/middleware"
	"wallet-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewMySQLConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := database.NewUserRepositoryImpl(db)
	walletRepo := database.NewWalletRepositoryImpl(db)
	transactionRepo := database.NewTransactionRepositoryImpl(db)
	apiLogRepo := database.NewAPILogRepositoryImpl(db)
	txManager := database.NewTransactionManagerImpl(db)

	registerUsecase := usecase.NewRegisterUsecase(userRepo, walletRepo)
	loginUsecase := usecase.NewLoginUsecase(userRepo)
	getBalanceUsecase := usecase.NewGetBalanceUsecase(walletRepo)
	withdrawUsecase := usecase.NewWithdrawUsecase(userRepo, walletRepo, transactionRepo, txManager)

	jwtService := jwt.NewJWTService(cfg.JWTSecret)

	registerHandler := handler.NewRegisterHandler(registerUsecase)
	loginHandler := handler.NewLoginHandler(loginUsecase, jwtService)
	balanceHandler := handler.NewBalanceHandler(getBalanceUsecase)
	withdrawHandler := handler.NewWithdrawHandler(withdrawUsecase)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware(apiLogRepo))

	api := router.Group("/api")
	{
		api.POST("/register", registerHandler.Handle)
		api.POST("/login", loginHandler.Handle)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			protected.GET("/balance", balanceHandler.Handle)
			protected.POST("/withdraw", withdrawHandler.Handle)
		}
	}

	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server running on %s", serverAddr)
	
	err = router.Run(serverAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
