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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryController struct {
	router *gin.Engine
	cu     usecase.CategoryUsecase
}

func (cc *CategoryController) Create(c *gin.Context) {
	var category model.CategoryInput
	err := c.ShouldBindJSON(&category)
	if err != nil {
		fmt.Printf("error an c.ShouldBindJSON(&category): %v", err)
		errorMessage := gin.H{"errors": "FAILED_TO_PROCESS_CATEGORY_REQUEST"}
		response := helper.APIResponse("Create Category Failed", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if !helper.ValidateStruct(c, category) {
		return
	}

	exists, err := cc.cu.CategoryNameExits(category.CategoryName)
	if err != nil {
		fmt.Printf("error an cc.cu.CategoryNameExits(category.CategoryName): %v", err)
		errorMessage := gin.H{"errors": "FAILED_TO_CHECK_NAME_EXISTENCE"}
		response := helper.APIResponse("FAILED_TO_CHECK_NAME_EXISTENCE", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	statusCode := http.StatusOK
	message := "CREATE_CATEGORY_SUCCESS"

	if exists {
		// Jika email sudah ada dalam basis data, respons akan menampilkan pesan bahwa email sudah ada.
		statusCode = http.StatusConflict
		message = "CATEGORY_NAME_ALREADY"
	} else {
		// Jika pendaftaran berhasil, respons akan menampilkan pesan sukses dan data pengguna yang terdafta
		if err := cc.cu.CreateCategory(category); err != nil {
			fmt.Printf("error an cc.cu.CreateCategory(category): %v", err)
			errorMessage := gin.H{"errors": "FAILED_TO_CREATE_CATEGORY"}
			response := helper.APIResponse("Create Category Failed", http.StatusBadRequest, "error", errorMessage)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	response := helper.APIResponse(message, statusCode, "Success", category)
	c.JSON(http.StatusOK, response)
}

func (cc *CategoryController) Update(c *gin.Context) {
	categoryId := c.Param("id")
	if categoryId == "" {
		fmt.Printf("err on categoryId := c.Param: %v", categoryId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_ID"})
		return
	}

	var updatedCategory model.CategoryUpdate
	err := c.ShouldBindJSON(&updatedCategory)
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("UPDATE_FAILED", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if !helper.ValidateStruct(c, updatedCategory) {
		return
	}

	exists, err := cc.cu.CheckNameForUpdate(updatedCategory.CategoryName, categoryId)
	if err != nil {
		fmt.Printf("err on cc.cu.CheckNameForUpdate(updatedCategory.CategoryName, categoryId) %v", err)
		response := helper.APIResponse("'FAILED_TO_RETRIEVE_DATA", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if exists {
		fmt.Println("err on exists:", exists)
		response := helper.APIResponse("NAME_ALREADY_EXISTS", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err := cc.cu.Update(categoryId, updatedCategory); err != nil {
		fmt.Printf("err on cc.cu.Update(categoryId, updatedCategory) %v", err)
		response := helper.APIResponse("UPDATE_FAILED", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("UPDATE_SUCCESS", http.StatusOK, "Success", updatedCategory)
	c.JSON(http.StatusOK, response)
}

func (cc *CategoryController) FindById(c *gin.Context) {
	ID := c.Param("id")

	// Check if ID is empty
	if len(ID) == 0 {
		response := helper.APIResponse("Category ID is empty", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category, err := cc.cu.FindById(ID)
	if err != nil {
		fmt.Printf("error on cc.cu.FindById(param.ID): %v", err)
		response := helper.APIResponse("ID NOT FOUND", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Get Category ID", http.StatusOK, "success", category)
	c.JSON(http.StatusOK, response)
}

func (cc *CategoryController) GetAll(c *gin.Context) {
	search, page, limit, err := utils.ParsePaginationParams(c)
	if err != nil {
		response := map[string]interface{}{
			"status":  "ERROR",
			"message": "Invalid pagination parameters",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category, totalDocs, err := cc.cu.GetAll(c, page, limit, search)
	if err != nil {
		response := helper.APIResponse("GET_CATEGORY_FAILED", http.StatusInternalServerError, "ERROR", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	Formatter := formatter.FormatterCategory(category)
	pagination := utils.GeneratePaginationData(totalDocs, page, limit)
	response := helper.APIResponse("Successfully Get Category", http.StatusOK, "success", Formatter)

	data := gin.H{
		"category":   response,
		"pagination": pagination,
	}
	c.JSON(http.StatusOK, data)
}

func (cc *CategoryController) Delete(c *gin.Context) {
	categoryId := c.Param("id")
	// Jika userID kosong, respons akan menunjukkan kesalahan karena ID pengguna tidak valid.
	if categoryId == "" {
		fmt.Printf("err %v", categoryId)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "INVALID CATEGORY ID"})
		return
	}

	if _, err := uuid.Parse(categoryId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "INVALID CATEGORY ID"})
		return
	}

	// Melakukan pencarian informasi pengguna berdasarkan userID.
	err := cc.cu.Delete(categoryId)
	// Jika terjadi kesalahan dalam pencarian pengguna, respons akan menunjukkan kesalahan yang terjadi.
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("FAILED_TO_DELETE", http.StatusUnprocessableEntity, "ERROR", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Jika pencarian pengguna berhasil, respons akan berisi informasi pengguna yang ditemukan.
	response := helper.APIResponse("DELETE_SUCCESS", http.StatusOK, "SUCCESS", nil)
	c.JSON(http.StatusOK, response)
}

func NewCategoryController(r *gin.Engine, cu usecase.CategoryUsecase) *CategoryController {
	category := CategoryController{
		router: r,
		cu:     cu,
	}

	route := r.Group("/category")
	route.Use(jwt.JwtAuthMiddleware())

	route.GET("/index", category.GetAll)
	route.POST("/create", category.Create)
	route.GET("/:id", category.FindById)
	route.PUT("/:id", category.Update)
	route.DELETE("/:id", category.Delete)

	return &category
}
