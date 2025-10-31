package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// WithAuth xác thực access token trong header Authorization: Bearer <token>
// - Dùng MapClaims để tương thích mọi kiểu claims (tránh nhầm với claims của refresh token).
// - Trích xuất uid (uid|user_id|sub) và role.
func WithAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// chỉ chấp nhận HMAC (HS256/384/512)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// Lấy uid từ các key phổ biến: uid | user_id | sub
		uid := extractInt(claims["uid"])
		if uid == 0 {
			uid = extractInt(claims["user_id"])
		}
		if uid == 0 {
			uid = extractInt(claims["sub"])
		}

		role := ""
		if v, ok := claims["role"].(string); ok {
			role = v
		} else if v, ok := claims["rol"].(string); ok {
			role = v
		}

		if uid == 0 || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		c.Set("uid", uid)
		c.Set("role", role)
		c.Next()
	}
}

// extractInt: hỗ trợ số (float64 từ JSON), int, hoặc string số
func extractInt(v any) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case int:
		return x
	case int64:
		return int(x)
	case string:
		if i, err := strconv.Atoi(x); err == nil {
			return i
		}
	}
	return 0
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	allow := map[string]struct{}{}
	for _, r := range roles {
		allow[r] = struct{}{}
	}
	return func(c *gin.Context) {
		roleAny, ok := c.Get("role")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if _, ok := allow[roleAny.(string)]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
