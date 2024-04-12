package controller

import (
	jwt "cashier/delivery/middleware"
	"cashier/formatter"
	"cashier/helper"
	"cashier/model"
	"cashier/usecase"
	"cashier/utils"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	router  *gin.Engine
	usecase usecase.UserUsecase
}

func (uc *UserController) GetAllUser(c *gin.Context) {
	search := c.Query("search")

	// Mengambil nilai halaman dan batas dari query string, defaultnya adalah halaman 1 dan batas 10.
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Mengonversi nilai halaman dan batas ke dalam tipe data integer.
	page, errPage := strconv.Atoi(pageStr)
	limit, errLimit := strconv.Atoi(limitStr)
	if errPage != nil || errLimit != nil {
		response := helper.APIResponse("INVALID_QUERY_PARAMS", http.StatusBadRequest, "ERROR", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var userFormatters []*formatter.UserFormatter
	var totalDocs int
	var err error

	if search != "" {
		userFormatters, totalDocs, err = uc.usecase.GetAllUser(c, page, limit, search)
	} else {
		userFormatters, totalDocs, err = uc.usecase.GetAllUser(c, page, limit, search)
	}

	if err != nil {
		response := helper.APIResponse("GET_USERS_FAILED", http.StatusInternalServerError, "ERROR", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	paginationUser := utils.GeneratePaginationData(totalDocs, page, limit)

	// Menyiapkan respons dengan pesan sukses dan data pengguna yang berhasil diambil.
	response := helper.APIResponse("SUCCESSFULLY_RETRIEVED_USERS", http.StatusOK, "SUCCESS", nil)

	// Menyusun data respons, data pengguna, dan data paginasi dalam sebuah map.
	data := gin.H{
		"user":       userFormatters,
		"response":   response,
		"pagination": paginationUser,
	}
	c.JSON(http.StatusOK, data)
}

func (uc *UserController) Update(c *gin.Context) {
	//Memeriksa ID pengguna yang diberikan dalam URL.
	userID := c.Param("id")
	if userID == "" {
		fmt.Printf("err on params: %v", userID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_ID"})
		return
	}

	//Mengikat data JSON yang diterima dengan struktur model UpdateUserInput.
	var update model.UpdateUserInput
	err := c.ShouldBindJSON(&update)
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("FAILED_TO_PROCESS_PRODUCT_UPDATE_REQUEST", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if !helper.ValidateStruct(c, update) { return }

	// Pengecekan apakah email yang baru ingin diubah tidak sama dengan email yang sudah ada di database
	exists, err := uc.usecase.CheckEmailForUpdate(update.Email, userID)
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("'FAILED_TO_RETRIEVE_DATA", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if exists {
		fmt.Println("exists:", exists)
		response := helper.APIResponse("EMAIL_ALREADY_EXISTS", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Update data pengguna
	err = uc.usecase.UpdateUser(userID, update)
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("FAILED_TO_UPDATE_DATA", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// Berhasil melakukan pembaruan
	response := helper.APIResponse("UPDATE_SUCCESSFUL", http.StatusOK, "success", update)
	c.JSON(http.StatusOK, response)
}

// FindByID mengambil informasi pengguna berdasarkan ID yang diberikan.
func (uc *UserController) FindByID(c *gin.Context) {

	// userID akan menyimpan nilai parameter ID dari request.
	userID := c.Param("id")

	// Jika userID kosong, respons akan menunjukkan kesalahan karena ID pengguna tidak valid.
	if userID == "" {
		fmt.Printf("err %v", userID)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "INVALID USER ID"})
		return
	}

	// Melakukan pencarian informasi pengguna berdasarkan userID.
	user, err := uc.usecase.FindById(userID)

	// Jika terjadi kesalahan dalam pencarian pengguna, respons akan menunjukkan kesalahan yang terjadi.
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("FAILED_TO_RETRIEVE_USER", http.StatusUnprocessableEntity, "ERROR", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Jika pencarian pengguna berhasil, respons akan berisi informasi pengguna yang ditemukan.
	response := helper.APIResponse("USER_FOUND_SUCCESSFULLY", http.StatusOK, "SUCCESS", user)
	c.JSON(http.StatusOK, response)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// Jika userID kosong, respons akan menunjukkan kesalahan karena ID pengguna tidak valid.
	if userID == "" {
		fmt.Printf("err %v", userID)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "INVALID USER ID"})
		return
	}

	if _, err := uuid.Parse(userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "INVALID USER ID"})
		return
	}

	// Melakukan pencarian informasi pengguna berdasarkan userID.
	err := uc.usecase.Delete(userID)
	// Jika terjadi kesalahan dalam pencarian pengguna, respons akan menunjukkan kesalahan yang terjadi.
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("FAILED_TO_DELETE_USER", http.StatusUnprocessableEntity, "ERROR", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Jika pencarian pengguna berhasil, respons akan berisi informasi pengguna yang ditemukan.
	response := helper.APIResponse("DELETE_SUCCESS", http.StatusOK, "SUCCESS", nil)
	c.JSON(http.StatusOK, response)
}

func NewUserController(r *gin.Engine, usecase usecase.UserUsecase) *UserController {
	controller := UserController{
		router:  r,
		usecase: usecase,
	}

	userRoute := r.Group("/users")
	userRoute.Use(jwt.JwtAuthMiddleware())

	userRoute.GET("/index", jwt.JwtAuthMiddleware(), controller.GetAllUser)
	userRoute.PUT("/update/:id", jwt.JwtAuthMiddleware(), controller.Update)
	userRoute.GET("/:id", jwt.JwtAuthMiddleware(), controller.FindByID)
	userRoute.DELETE("/delete/:id", jwt.JwtAuthMiddleware(), controller.DeleteUser)

	return &controller
}
