package main

import (
	"agnos-gin/config"
	"agnos-gin/internal/domain"
	"agnos-gin/internal/handler"
	"agnos-gin/internal/infrastructure"
	"agnos-gin/internal/middleware"
	"agnos-gin/internal/repository"
	"agnos-gin/internal/service"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// 2. Connect DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	var db *gorm.DB

	// Retry connection logic
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Database connected successfully.")
			break
		}
		log.Printf("Waiting for database... (%d/10) Error: %v\n", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Cannot connect to DB after retries:", err)
	}

	// 3. Migrate (Hospital must be first due to FK dependencies)
	if err := db.AutoMigrate(&domain.Hospital{}, &domain.Staff{}, &domain.Patient{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	// 4. Seed default hospitals
	hospitals := []domain.Hospital{
		{Name: "Hospital A"},
		{Name: "Hospital B"},
	}
	for _, h := range hospitals {
		// Only create if not exists
		db.FirstOrCreate(&h, domain.Hospital{Name: h.Name})
	}
	log.Println("Default hospitals seeded.")

	// 5. Init Layers
	// Repos
	hospitalRepo := repository.NewHospitalRepository(db)
	staffRepo := repository.NewStaffRepository(db)
	patientRepo := repository.NewPatientRepository(db)

	// Infrastructure
	hospitalClient := infrastructure.NewHospitalApiClient(cfg.HospitalAPIURL)

	// Services
	authService := service.NewAuthService(staffRepo, hospitalRepo, cfg.JWTSecret)
	patientService := service.NewPatientService(patientRepo, hospitalClient)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	patientHandler := handler.NewPatientHandler(patientService)

	// 6. Router
	r := gin.Default()

	// Serve static files (CSS, JS)
	r.Static("/static", "./web")

	// Serve the frontend page
	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	// Routes
	r.POST("/staff/create", authHandler.Register)
	r.POST("/staff/login", authHandler.Login)

	protected := r.Group("/patient")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/search", patientHandler.Search)
	}

	// Run
	log.Println("Server starting on :8080")
	r.Run(":8080")
}
