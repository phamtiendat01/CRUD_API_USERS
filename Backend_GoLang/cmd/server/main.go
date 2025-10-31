package main

import (
	docs "crud_api_us/docs"
	"crud_api_us/internal/router"
)

// @title           User API (Gin + Swagger)
// @version         1.0
// @description     CRUD người dùng mẫu, sạch và tối giản.
// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Nhập dạng: Bearer <JWT>
func main() {
	// Thiết lập runtime cho Swagger (không bắt buộc nhưng rõ ràng)
	docs.SwaggerInfo.Title = "User API (Gin + Swagger)"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"

	r := router.New()
	_ = r.Run(":8080")
}
