package main

import (
	"os"

	docs "crud_api_us/docs"
	"crud_api_us/internal/router"
)

// @title           User API (Gin + Swagger)
// @version         1.0
// @description     CRUD người dùng mẫu, sạch và tối giản.
// @BasePath        /api/v1
// @schemes         http https
// (NÊN bỏ @host để Swagger tự dùng host hiện tại, hoặc để rỗng qua runtime)

func main() {
	// Lấy PORT do Railway/Render/Heroku cấp, fallback 8080 khi chạy local
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Thiết lập thông tin Swagger (nên để Host động, tránh hardcode localhost)
	docs.SwaggerInfo.Title = "User API (Gin + Swagger)"
	docs.SwaggerInfo.Version = "1.0"
	// Nếu bạn thật sự muốn set host thủ công cho production, dùng env:
	if h := os.Getenv("SWAGGER_HOST"); h != "" {
		docs.SwaggerInfo.Host = h // ví dụ "your-service.up.railway.app"
	} else {
		docs.SwaggerInfo.Host = "" // để trống => Swagger UI dùng current host
	}
	docs.SwaggerInfo.BasePath = "/api/v1"

	r := router.New()

	// Bind đúng cổng
	if err := r.Run(":" + port); err != nil {
		panic(err)
	}
}
