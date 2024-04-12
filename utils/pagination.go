package utils

// PaginationData represents pagination information.
type PaginationData struct {
	TotalDocs     int
	Limit         int
	TotalPages    int
	Page          int
	PagingCounter int
	HasPrevPage   bool
	HasNextPage   bool
	PrevPage      int
	NextPage      int
}

// GeneratePaginationData generates pagination data based on input parameters.
func GeneratePaginationData(totalDocs, limit, page int) PaginationData {
	totalPages := (totalDocs + limit - 1) / limit
	hasPrevPage := page > 1
	hasNextPage := page < totalPages
	prevPage := page - 1
	nextPage := page + 1

	// Use the correct uppercase spelling for PagingCounter
	PagingCounter := page

	return PaginationData{
		TotalDocs:     totalDocs,
		Limit:         limit,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: PagingCounter,
		HasPrevPage:   hasPrevPage,
		HasNextPage:   hasNextPage,
		PrevPage:      prevPage,
		NextPage:      nextPage,
	}
}
