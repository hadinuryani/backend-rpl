package middleware

import (
	"strings"
	"time"

	"ic-plus-backend/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a configured CORS middleware.
// It allows origins from config, Vercel domains, and ngrok URLs.
func CORSMiddleware() gin.HandlerFunc {
	cfg := config.GetConfig()

	allowedSet := make(map[string]bool, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		allowedSet[strings.TrimSpace(o)] = true
	}

	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {

			// DEBUG LOG
			println("CORS ORIGIN:", origin)

			// Allow configured origins
			if allowedSet[origin] {
				return true
			}

			// Allow all Vercel domains
			if strings.Contains(origin, "vercel.app") {
				return true
			}

			// Allow any ngrok tunnel
			if strings.HasSuffix(origin, ".ngrok-free.app") ||
				strings.HasSuffix(origin, ".ngrok.io") ||
				strings.HasSuffix(origin, ".ngrok-free.dev") ||
				strings.Contains(origin, "ngrok") {
				return true
			}

			return false
		},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
			"HEAD",
		},

		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"Origin",
			"Accept",
			"X-Requested-With",
			"ngrok-skip-browser-warning",
		},

		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	})
}