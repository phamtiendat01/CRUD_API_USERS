package router

import (
	"fmt"
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
		raw = "http://localhost:5173"
	}
	for _, s := range strings.Split(raw, ",") {
		s = strings.TrimSpace(s)
		s = strings.TrimRight(s, "/") // QUAN TRỌNG: bỏ dấu "/" cuối
		if s == "" {
			continue
		}
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
	debugCORS := os.Getenv("DEBUG_CORS") == "1"

	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Dùng AllowOriginFunc để hỗ trợ wildcard
	cfg.AllowOriginFunc = func(origin string) bool {
		u, err := url.Parse(origin)
		if err != nil {
			if debugCORS {
				fmt.Println("[CORS] invalid origin:", origin)
			}
			return false
		}
		o := u.Scheme + "://" + u.Host
		host := u.Hostname()

		for _, e := range exact {
			if o == e {
				if debugCORS {
					fmt.Println("[CORS] ALLOW exact:", origin)
				}
				return true
			}
		}
		for _, suf := range suffixes {
			if strings.HasSuffix(host, suf) {
				if debugCORS {
					fmt.Println("[CORS] ALLOW wildcard:", origin, "matches *.", suf)
				}
				return true
			}
		}
		if debugCORS {
			fmt.Println("[CORS] DENY:", origin)
		}
		return false
	}

	r.Use(cors.New(cfg))

	// Bảo đảm mọi preflight đều được trả về 204 (phòng khi middleware bị lách)
	r.OPTIONS("/*path", func(c *gin.Context) { c.Status(204) })
	// ===== End CORS =====

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ===== DB & DI (giữ nguyên phần còn lại của bạn) =====
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

	userRepo := repository.NewMySQLUserRepo(db)
	authRepo := repository.NewMySQLAuthRepo(db)
	jwtCfg := services.LoadJWTConfigFromEnv()

	u := handlers.NewUserHandler(userRepo)
	a := handlers.NewAuthHandler(userRepo, authRepo, jwtCfg)
	authMW := middleware.WithAuth(jwtCfg.Secret)

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
