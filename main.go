package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ic-plus-backend/config"
	"ic-plus-backend/database"
	"ic-plus-backend/handlers"
	"ic-plus-backend/middleware"
	"ic-plus-backend/repositories"
	"ic-plus-backend/routes"
	"ic-plus-backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load config
	cfg := config.GetConfig()
	log.Println("Starting IC+ Backend...")
	log.Printf("Environment: %s", cfg.AppEnv)
	log.Printf("Port: %s", cfg.AppPort)

	// 2. Connect to database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer database.Close()
	db := database.GetDB()

	// 3. Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	pasienRepo := repositories.NewPasienRepository(db)
	antrianRepo := repositories.NewAntrianRepository(db)
	rmRepo := repositories.NewRekamMedisRepository(db)
	jadwalRepo := repositories.NewJadwalRepository(db)
	notifRepo := repositories.NewNotifikasiRepository(db)
	invRepo := repositories.NewInventoriRepository(db)

	// 4. Initialize services
	authService := services.NewAuthService(userRepo, pasienRepo)
	antrianService := services.NewAntrianService(antrianRepo, db)
	rmService := services.NewRekamMedisService(rmRepo, antrianRepo, db)
	jadwalService := services.NewJadwalService(jadwalRepo)
	notifService := services.NewNotifikasiService(notifRepo)
	invService := services.NewInventoriService(invRepo)
	_ = notifService

	// 5. Initialize WA gateway (real Fonnte or stub for development)
	var waGateway services.WAGateway
	if cfg.WAAPIToken != "" && cfg.WAAPIToken != "your-fonnte-api-token" {
		waGateway = services.NewFonnteWAGateway(cfg.WAGatewayURL, cfg.WAAPIToken)
		log.Println("📱 WhatsApp Gateway: Fonnte (LIVE)")
	} else {
		waGateway = &services.StubWAGateway{}
		log.Println("📱 WhatsApp Gateway: Stub (console only — set WA_API_TOKEN to enable)")
	}
	scheduler := services.NewSchedulerService(jadwalRepo, notifRepo, waGateway)

	// 6. Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	pasienHandler := handlers.NewPasienHandler(pasienRepo)
	bidanHandler := handlers.NewBidanHandler(pasienRepo, userRepo, db)
	klinikHandler := handlers.NewKlinikHandler(db)
	antrianHandler := handlers.NewAntrianHandler(antrianService, antrianRepo, pasienRepo)
	rmHandler := handlers.NewRekamMedisHandler(rmService, rmRepo, pasienRepo, db)
	resepHandler := handlers.NewResepHandler(rmRepo)
	jadwalHandler := handlers.NewJadwalHandler(jadwalService, jadwalRepo, pasienRepo, db)
	notifHandler := handlers.NewNotifikasiHandler(notifRepo, pasienRepo)
	invHandler := handlers.NewInventoriHandler(invService, invRepo, db)
	dashHandler := handlers.NewDashboardHandler(antrianRepo, invRepo, db)
	monitorHandler := handlers.NewMonitorHandler(db)

	// 7. Setup Gin router
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// 8. Register routes
	routes.RegisterRoutes(r, authHandler, pasienHandler, bidanHandler, klinikHandler, antrianHandler, rmHandler, resepHandler, jadwalHandler, notifHandler, invHandler, dashHandler, monitorHandler)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format(time.RFC3339)})
	})

	// Test WhatsApp sending endpoint
	r.POST("/api/v1/test-wa", func(c *gin.Context) {
		var req struct {
			Target  string `json:"target" binding:"required"`
			Message string `json:"message" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Target dan message wajib diisi"})
			return
		}
		err := waGateway.SendMessage(req.Target, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "WhatsApp berhasil dikirim!"})
	})

	// 9. Start scheduler
	scheduler.Start()
	defer scheduler.Stop()

	// 10. Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		log.Printf("Server running on http://localhost:%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}
	log.Println("Server stopped gracefully")
}
