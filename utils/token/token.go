package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	AccessTokenLifetime  = 5 * time.Hour
	RefreshTokenLifetime = 5 * time.Hour
)

// GenerateToken menghasilkan token akses JWT.
func GenerateToken(userID, userRole string) (string, error) {
	expirationTime := time.Now().Add(AccessTokenLifetime).Unix()

	claims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
		"user_role":  userRole,
		"exp":        expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func GenerateRefreshToken(userID, userRole string) (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	claims := jwt.MapClaims{
		"sub":       userID,
		"user_role": userRole,
		"exp":       time.Now().Add(RefreshTokenLifetime).Unix(),
	}

	refreshToken.Claims = claims

	return refreshToken.SignedString([]byte(os.Getenv("API_SECRET")))
}

// TokenValid memeriksa validitas token akses JWT.
func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			return fmt.Errorf("TOKEN_EXPIRED")
		}
		return nil
	} else {
		return fmt.Errorf("INVALID_TOKEN")
	}
}

// ExtractToken mengambil token dari permintaan HTTP.
func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

// ExtractTokenID mengambil ID pengguna dari token akses JWT.
func ExtractTokenID(c *gin.Context) (uint, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}

// ExtractBearerToken mengambil token bearer dari header.
func ExtractBearerToken(header string) string {
	parts := strings.Split(header, " ")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

func VerifyRefreshToken(refreshToken string) (jwt.MapClaims, error) {
	// Ambil secret key dari variabel lingkungan
	secret := os.Getenv("API_SECRET")
	if secret == "" {
		// Jika secret key tidak diatur, kembalikan error
		return nil, fmt.Errorf("JWT_REFRESH_TOKEN_SECRET is not set")
	}

	// Parse refreshToken dengan secret key yang diambil
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode tipe token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Gunakan secret key sebagai interface untuk verifikasi token
		return []byte(secret), nil
	})
	if err != nil {
		// Jika terjadi error dalam parsing token, kembalikan error
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	expirationTime, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("expiration time not found in claims")
	}

	expiration := time.Unix(int64(expirationTime), 0)
	if time.Now().After(expiration) {
		return nil, fmt.Errorf("token expired")
	}

	// Kembalikan klaim JWT jika token valid
	return claims, nil
}
