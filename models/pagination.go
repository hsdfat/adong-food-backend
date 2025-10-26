package models

// PaginationParams contains pagination parameters from query string
type PaginationParams struct {
	Page     int    `form:"page" binding:"omitempty,min=0"`
	PageSize int    `form:"per_page" binding:"omitempty,min=0,max=100"` // Changed to per_page to match common conventions
	Search   string `form:"search"`
	SortBy   string `form:"sort_by"`
	SortDir  string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

// PaginationMeta contains pagination metadata matching frontend ResourceCollection interface
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
}

// ResourceCollection is the response wrapper matching frontend interface
type ResourceCollection struct {
	Data interface{}     `json:"data"`
	Meta *PaginationMeta `json:"meta"`
}

// GetPaginationParams extracts and validates pagination parameters with defaults
func GetPaginationParams(page, pageSize int, search, sortBy, sortDir string) PaginationParams {
	// Set defaults
	// if page < 1 {
	// 	page = 1
	// }
	// if pageSize < 1 {
	// 	pageSize = 10
	// }
	// if pageSize > 100 {
	// 	pageSize = 100
	// }
	if sortDir == "" {
		sortDir = "asc"
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		SortBy:   sortBy,
		SortDir:  sortDir,
	}
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(page, perPage int, total int64) *PaginationMeta {
	if perPage < 1 || page < 1 {
		return &PaginationMeta{
			CurrentPage: page,
			LastPage:    1,
			From:        0,
			To:          int(total),
			PerPage:     perPage,
			Total:       int(total),
		}
	}
	totalInt := int(total)
	lastPage := (totalInt + perPage - 1) / perPage
	if lastPage < 1 {
		lastPage = 1
	}

	// Calculate from and to
	from := 0
	to := 0
	if totalInt > 0 {
		from = (page-1)*perPage + 1
		to = from + perPage - 1
		if to > totalInt {
			to = totalInt
		}
	}

	return &PaginationMeta{
		CurrentPage: page,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		PerPage:     perPage,
		Total:       totalInt,
	}
}
