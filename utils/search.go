package utils

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// SearchConfig defines which fields to search and their weight
type SearchConfig struct {
	Fields []string
	Fuzzy  bool // Use ILIKE (case-insensitive) vs exact match
}

// ApplySearch applies search conditions to a GORM query
func ApplySearch(db *gorm.DB, search string, config SearchConfig) *gorm.DB {
	if search == "" || len(config.Fields) == 0 {
		return db
	}

	// Trim and prepare search term
	search = strings.TrimSpace(search)

	// Build search conditions
	var conditions []string
	var args []interface{}

	for _, field := range config.Fields {
		if config.Fuzzy {
			conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", field))
			args = append(args, "%"+search+"%")
		} else {
			conditions = append(conditions, fmt.Sprintf("%s = ?", field))
			args = append(args, search)
		}
	}

	// Combine with OR
	query := strings.Join(conditions, " OR ")
	return db.Where(query, args...)
}

// ApplySort applies sorting to a GORM query
func ApplySort(db *gorm.DB, sortBy, sortDir string, allowedFields map[string]string) *gorm.DB {
	if sortBy == "" {
		return db
	}

	// Validate sort field
	dbField, ok := allowedFields[sortBy]
	if !ok {
		return db
	}

	// Validate sort direction
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc"
	}

	return db.Order(fmt.Sprintf("%s %s", dbField, strings.ToUpper(sortDir)))
}

// ApplyPagination applies pagination to a GORM query
func ApplyPagination(db *gorm.DB, page, pageSize int) *gorm.DB {
	if page < 1 || pageSize < 1 {
		return db
	}
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}
