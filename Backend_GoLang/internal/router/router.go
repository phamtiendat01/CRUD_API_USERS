package router

import (
	"os"
	"time"

	_ "crud_api_us/docs" // nạp swagger docs

	"crud_api_us/internal/database"
	"crud_api_us/internal/handlers"
	"crud_api_us/internal/middleware"
	"crud_api_us/internal/models"
	"crud_api_us/internal/repository"
	"crud_api_us/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/crypto/bcrypt"
)

func New() *gin.Engine {
	r := gin.Default()
	// CORS cho FE (Vite port 5173)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// DB
	cfg := database.LoadConfigFromEnv()
	db, err := database.Open(cfg)
	if err != nil {
		panic("cannot connect MySQL: " + err.Error())
	}

	// Migrate & ensure admin (idempotent)
	if err := db.AutoMigrate(&models.User{}); err != nil {
		panic("migrate failed: " + err.Error())
	}
	hash := func(p string) string {
		b, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		return string(b)
	}
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@example.com"
	}
	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminPass == "" {
		adminPass = "Admin@123"
	}
	var admin models.User
	if err := db.
		Where("email = ?", adminEmail).
		Attrs(models.User{
			Username:     "admin",
			FullName:     "Administrator",
			Role:         "admin",
			PasswordHash: hash(adminPass), // chỉ dùng khi create mới
		}).
		FirstOrCreate(&admin).Error; err != nil {
		panic("ensure admin failed: " + err.Error())
	}
	if admin.Role != "admin" {
		_ = db.Model(&admin).Update("role", "admin").Error
	}

	// DI
	userRepo := repository.NewMySQLUserRepo(db)
	authRepo := repository.NewMySQLAuthRepo(db)
	jwtCfg := services.LoadJWTConfigFromEnv()

	// Handlers
	u := handlers.NewUserHandler(userRepo)
	a := handlers.NewAuthHandler(userRepo, authRepo, jwtCfg)

	// Middleware
	authMW := middleware.WithAuth(jwtCfg.Secret)

	// Routes
	v1 := r.Group("/api/v1")
	{
		// --- Public Auth ---
		v1.POST("/auth/register", a.Register)
		v1.POST("/auth/login", a.Login)
		v1.POST("/auth/logout", a.Logout)

		// --- Protected (đã đăng nhập) ---
		v1.GET("/auth/me", authMW, a.Me)

		// --- ADMIN ONLY ---
		admin := v1.Group("/admin", authMW, middleware.RequireRoles("admin"))
		{
			admin.GET("/users", u.ListUsers)         // @Security BearerAuth (khai báo trong handler)
			admin.GET("/users/:id", u.GetUser)       // @Security BearerAuth
			admin.POST("/users", u.CreateUser)       // @Security BearerAuth
			admin.PUT("/users/:id", u.UpdateUser)    // @Security BearerAuth
			admin.DELETE("/users/:id", u.DeleteUser) // @Security BearerAuth
		}
	}

	return r
}
