package middleware

import (
	"net/http"
	"cashier/utils/token"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			var errorMessage string
			switch err.Error() {
			case "TOKEN_EXPIRED":
				errorMessage = "Access token expired"
			case "INVALID_TOKEN":
				errorMessage = "Invalid access token"
			default:
				errorMessage = "Access token invalid"
			}

			// Set metadata error_message dengan pesan kesalahan
			c.Set("error_message", errorMessage)
			
			// Menyiapkan data JSON menggunakan gin.H
			responseData := gin.H{
				"status":  "error",
				"code":    http.StatusUnauthorized,
				"message": "ACCESS_TOKEN_INVALID",
				"error":   errorMessage,
			}
			
			// Mengirim respons JSON
			c.JSON(http.StatusUnauthorized, responseData)
			
			// Menghentikan eksekusi selanjutnya
			c.Abort()
			return
		}

		// Lanjutkan jika token valid
		c.Next()
	}
}

