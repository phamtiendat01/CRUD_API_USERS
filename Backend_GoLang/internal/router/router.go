package router

import (
	"net/url"
	"os"
	"strings"
	"time"

	_ "crud_api_us/docs"

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

func parseCORSOrigins() (exact []string, suffixes []string) {
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		// fallback để dev local không bị chặn
		raw = "http://localhost:5173"
	}
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		// Cho phép pattern dạng "*.domain.com"
		if strings.HasPrefix(s, "*.") {
			suffixes = append(suffixes, strings.TrimPrefix(s, "*."))
		} else {
			exact = append(exact, s)
		}
	}
	return
}

func New() *gin.Engine {
	r := gin.Default()

	// ===== CORS =====
	exact, suffixes := parseCORSOrigins()

	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,           // nếu bạn dùng cookie/refresh token
		MaxAge:           12 * time.Hour, // cache preflight
	}

	if len(suffixes) == 0 {
		// Không có wildcard -> dùng danh sách exact
		cfg.AllowOrigins = exact
	} else {
		// Có wildcard (*.domain) -> dùng AllowOriginFunc
		cfg.AllowOriginFunc = func(origin string) bool {
			u, err := url.Parse(origin)
			if err != nil {
				return false
			}
			// So khớp exact (bao gồm schema+host+port)
			o := u.Scheme + "://" + u.Host
			for _, e := range exact {
				if e == o {
					return true
				}
			}
			// So khớp wildcard theo suffix host
			host := u.Hostname()
			for _, suf := range suffixes {
				if strings.HasSuffix(host, suf) {
					return true
				}
			}
			return false
		}
	}

	r.Use(cors.New(cfg))
	// ===== End CORS =====

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// DB, migrate, seed admin, DI, routes ... (giữ như bạn đã có)
	// ...
	// (phần dưới không đổi)

	// DB
	cfgDB := database.LoadConfigFromEnv()
	db, err := database.Open(cfgDB)
	if err != nil {
		panic("cannot connect MySQL: " + err.Error())
	}
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
			PasswordHash: hash(adminPass),
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
		v1.POST("/auth/register", a.Register)
		v1.POST("/auth/login", a.Login)
		v1.POST("/auth/logout", a.Logout)

		v1.GET("/auth/me", authMW, a.Me)

		admin := v1.Group("/admin", authMW, middleware.RequireRoles("admin"))
		{
			admin.GET("/users", u.ListUsers)
			admin.GET("/users/:id", u.GetUser)
			admin.POST("/users", u.CreateUser)
			admin.PUT("/users/:id", u.UpdateUser)
			admin.DELETE("/users/:id", u.DeleteUser)
		}
	}

	return r
}
