package controller

import (
	jwt "cashier/delivery/middleware"
	"cashier/helper"
	"cashier/model"
	"cashier/usecase"
	"cashier/utils/token"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// ErrInvalidRefreshToken is returned when the refresh token is invalid.
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthController struct {
	router  *gin.Engine
	usecase usecase.UserUsecase
}

func (ac *AuthController) RegisterUser(c *gin.Context) {
	// Mengikat data JSON dari request ke strac t model RegisterUserInput.
	var register model.RegisterUserInput
	err := c.ShouldBindJSON(&register)
	// Jika terjadi kesalahan dalam pemrosesan JSON, respons akan menunjukkan kesalahan tersebut.
	if err != nil {
		errorMessage := gin.H{"errors": "FAILED_TO_PROCESS_REGISTER_REQUEST"}
		response := helper.APIResponse("FAILED_TO_REGISTER_USER", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Validasi struktur data register.
	if !helper.ValidateStruct(c, register) {
		return
	}

	// Memeriksa keberadaan email dalam basis data.
	exists, err := ac.usecase.EmailExits(register.Email)
	if err != nil {
		errorMessage := gin.H{"errors": "FAILED_TO_CHECK_EMAIL_EXISTENCE"}
		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	statusCode := http.StatusOK
	message := "CREATE_SUCCESSFUL"

	if exists {
		// Jika email sudah ada dalam basis data, respons akan menampilkan pesan bahwa email sudah ada.
		statusCode = http.StatusConflict
		message = "EMAIL_ALREADY_EXISTS"
	} else {
		// Jika pendaftaran berhasil, respons akan menampilkan pesan sukses dan data pengguna yang terdafta
		if err := ac.usecase.RegisterUser(register); err != nil {
			response := helper.APIResponse("FAILED_TO_REGISTER_USER", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	// Berhasil melakukan register
	response := helper.APIResponse(message, statusCode, "success", register)
	c.JSON(http.StatusOK, response)
}

func (ac *AuthController) Login(c *gin.Context) {
	// Mengecek apakah data JSON yang diterima sesuai dengan struktur model LoginUserInput.
	var input model.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errorMessage := gin.H{"ERRORS": "INVALID_JSON_FORMAT"}
		response := helper.APIResponse("LOGIN_FAILED", http.StatusBadRequest, "ERROR", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Validasi struktur data input.
	if !helper.ValidateStruct(c, input) {
		return
	}
	// Melakukan proses login pengguna.
	result, err := ac.usecase.LoginUser(input)
	if err != nil {

		// Jika terjadi kesalahan saat proses login, fungsi akan memberikan respons dengan kesalahan tersebut.
		fmt.Println("err on ac .usecase.LoginUser(input)", err)
		errorMessage := gin.H{"ERRORS": "LOGIN_FAILED"}
		response := helper.APIResponse("LOGIN_FAILED", http.StatusBadRequest, "ERROR", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Jika proses login berhasil, fungsi akan memberikan respons dengan pesan sukses dan data pengguna yang berhasil login.
	response := helper.APIResponse("SUCCESSFULLY_LOGIN", http.StatusOK, "success", result)
	c.JSON(http.StatusOK, response)
}


func (ac *AuthController) RefreshToken(c *gin.Context) {
	var request struct {
	RefreshToken string `json:"refreshToken"`
	}

	// Mendapatkan refreshToken dari form
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST_BODY"})
		return
	}

	// Memverifikasi refreshToken
	claims, err := token.VerifyRefreshToken(request.RefreshToken)
	if err != nil {
		// Respon jika refreshToken tidak valid
		c.JSON(http.StatusUnauthorized, gin.H{"error": "INVALID_REFRESH_TOKEN"})
		return
	}

	// Mendapatkan userID dan userRole dari klaim refreshToken
	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INVALID_CLAIMS"})
		return
	}
	userRole, ok := claims["user_role"].(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "INVALID_CLAIMS"})
		return
	}

	// Membuat payload baru untuk accessToken
	accessToken, err := token.GenerateToken(userID, userRole)
	if err != nil {
		// Respon jika gagal membuat accessToken baru
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	// Membuat refreshToken baru dengan masa berlaku yang sama
	newRefreshToken, err := token.GenerateRefreshToken(userID, userRole)
	if err != nil {
		// Respon jika gagal membuat refreshToken baru
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Respon dengan refreshToken baru dan accessToken
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": newRefreshToken,
	})
}

func NewAuthController(r *gin.Engine, usecase usecase.UserUsecase) *AuthController {
	controller := AuthController{
		router:  r,
		usecase: usecase,
	}

	auth := r.Group("/auth")

	auth.POST("auth/register", jwt.JwtAuthMiddleware(), controller.RegisterUser)
	r.POST("/auth/login", controller.Login)
	auth.POST("/refresh-token", controller.RefreshToken)

	return &controller
}
