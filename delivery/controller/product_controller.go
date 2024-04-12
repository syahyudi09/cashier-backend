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

type ProductController struct {
	router *gin.Engine
	pu     usecase.ProductUsecase
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	productName := c.PostForm("productName")
	priceStr := c.PostForm("price")
	categoryId := c.PostForm("categoryId")

	// Konversi priceStr ke float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		response := helper.APIResponse("Failed to parse price", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Membuat objek InputProduct dari nilai-nilai yang diambil
	product := model.InputProduct{
		ProductName: productName,
		Price:       price,
		CategoryId:  categoryId,
	}


	if !helper.ValidateStruct(c, product) {
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	product.Thumbnail = file.Filename

	path := "public/product/" + file.Filename

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		response := helper.APIResponse("Failed to upload avatar image", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	exists, err := pc.pu.ProductNameExits(product.ProductName)
	if err != nil {
		fmt.Printf("error an pc.pu.ProductNameExits(product.ProductName): %v", err)
		errorMessage := gin.H{"errors": "FAILED_TO_CHECK_NAME_EXISTENCE"}
		response := helper.APIResponse("FAILED_TO_CHECK_NAME_EXISTENCE", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	statusCode := http.StatusOK
	message := "CREATE_PRODUCT_SUCCESS"

	if exists {
		// Jika email sudah ada dalam basis data, respons akan menampilkan pesan bahwa email sudah ada.
		statusCode = http.StatusConflict
		message = "PRODUCT_NAME_ALREADY"
	} else {
		if err := pc.pu.Create(product); err != nil {
			fmt.Printf("error an cc.cu.CreateCategory(category): %v", err)
			errorMessage := gin.H{"errors": "FAILED_TO_CREATE_CATEGORY"}
			response := helper.APIResponse("Create Category Failed", http.StatusBadRequest, "error", errorMessage)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}

	response := helper.APIResponse(message, statusCode, "Success", product)
	c.JSON(http.StatusOK, response)
}

func (pc *ProductController) GetAll(c *gin.Context) {
	search, page, limit, err := utils.ParsePaginationParams(c)
	if err != nil {
		response := map[string]interface{}{
			"status":  "ERROR",
			"message": "Invalid pagination parameters",
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	product, totalDocs, err := pc.pu.GetAll(c, page, limit, search)
	if err != nil {
		fmt.Printf("error an ch.campaignUsecae.FindByID: %v", err)
		response := helper.APIResponse("GET_PRODUCT_FAILED", http.StatusInternalServerError, "ERROR", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	Formatter := formatter.FormatterProduct(product)
	pagination := utils.GeneratePaginationData(totalDocs, page, limit)
	response := helper.APIResponse("Successfully Get Category", http.StatusOK, "success", Formatter)

	data := gin.H{
		"product":    response,
		"pagination": pagination,
	}
	c.JSON(http.StatusOK, data)
}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	productId := c.Param("id")
	if productId == "" {
		fmt.Printf("err on categoryId := c.Param: %v", productId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_ID"})
		return
	}

	var updatedProduct model.UpdateProduct
	err := c.ShouldBindJSON(&updatedProduct)
	if err != nil {
		fmt.Printf("err %v", err)
		response := helper.APIResponse("UPDATE_FAILED", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if !helper.ValidateStruct(c, updatedProduct) {
		return
	}

	exists, err := pc.pu.CheckNameForUpdate(updatedProduct.ProductName, productId)
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

	if err := pc.pu.Update(productId, updatedProduct); err != nil {
		fmt.Printf("err on pc.pu.Update %v", err)
		response := helper.APIResponse("UPDATE_FAILED", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("UPDATE_SUCCESS", http.StatusOK, "Success", updatedProduct)
	c.JSON(http.StatusOK, response)
}

func (pc *ProductController) FindById(c *gin.Context) {
	ID := c.Param("id")

	// Check if ID is empty
	if len(ID) == 0 {
		response := helper.APIResponse("Category ID is empty", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := pc.pu.FindById(ID)
	if err != nil {
		fmt.Printf("error on cc.cu.FindById(param.ID): %v", err)
		response := helper.APIResponse("ID NOT FOUND", http.StatusNotFound, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.APIResponse("Get Category ID", http.StatusOK, "success", product)
	c.JSON(http.StatusOK, response)
}

func (pc *ProductController) Delete(c *gin.Context) {
	productId := c.Param("id")
	// Jika userID kosong, respons akan menunjukkan kesalahan karena ID pengguna tidak valid.
	if productId == "" {
		fmt.Printf("err %v", productId)
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "INVALID PRODUCT ID"})
		return
	}

	if _, err := uuid.Parse(productId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ERROR": "INVALID CATEGORY ID"})
		return
	}

	// Melakukan pencarian informasi pengguna berdasarkan userID.
	err := pc.pu.Delete(productId)
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

func NewProductController(r *gin.Engine, pu usecase.ProductUsecase) *ProductController {
	product := ProductController{
		router: r,
		pu:     pu,
	}

	authenticated := r.Group("/products")
	authenticated.Use(jwt.JwtAuthMiddleware())

	authenticated.POST("/create", product.CreateProduct) // create
	authenticated.GET("/index", product.GetAll)          // get all
	authenticated.PUT("/:id", product.UpdateProduct)     // update
	authenticated.GET("/:id", product.FindById)          // get All
	authenticated.DELETE("/:id", product.Delete)         // delete

	return &product
}
