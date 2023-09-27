package data

import (
	"math"
	"strings"

	"github.com/rehmanm/greenlight/internal/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

type MetaData struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func CalculateMetaData(totalRecords, page, pageSize int) MetaData {
	if totalRecords == 0 {
		return MetaData{}
	}
	return MetaData{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func ValidateFilters(v *validator.Validator, filters Filters) {

	v.Check(filters.Page > 0, "page", "must be greater than zero")
	v.Check(filters.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(filters.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(filters.PageSize <= 100, "page_size", "must be maximum of 100")

	v.Check(validator.PermittedValue(filters.Sort, filters.SortSafeList...), "sort", "invalid sort value")

}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameters" + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}
func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
