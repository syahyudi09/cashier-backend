package utils

import (
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
)
func ParsePaginationParams(c *gin.Context) (search string, page, limit int, err error) {
	search = c.Query("search")

	// Mengambil nilai halaman dan batas dari query string, defaultnya adalah halaman 1 dan batas 10.
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Mengonversi nilai halaman dan batas ke dalam tipe data integer.
	page, errPage := strconv.Atoi(pageStr)
	limit, errLimit := strconv.Atoi(limitStr)
	if errPage != nil || errLimit != nil {
		err = fmt.Errorf("invalid pagination parameters")
	}
	return search, page, limit, err
}