package utils

import "strconv"

type PaginationLimits struct {
	Limit int
	Page  int
}

type SortParams struct {
	SortParam     string
	SortDirection string
}

var (
	DateFormat = "02.01.2006."
)

// Sort constants
var (
	// Sort params

	SORT_NAME           = "name"
	SORT_NUM_WORKSPACES = "numWorkspaces"
	SORT_DATE_CREATED   = "dateCreated"
	SORT_PUBLIC_KEY_ID  = "publicKeyID"

	// Sort direction

	SORT_DIR_ASC  = "asc"
	SORT_DIR_DESC = "desc"
)

var (
	NoPagination = PaginationLimits{
		Limit: -1,
		Page:  -1,
	}

	DefaultSort = SortParams{
		SortParam:     SORT_NAME,
		SortDirection: SORT_DIR_ASC,
	}
)

func ParsePagination(limit string, page string) PaginationLimits {
	if limit != "" && page != "" {
		intLimit, limitErr := strconv.Atoi(limit)
		intPage, pageErr := strconv.Atoi(page)

		if limitErr != nil || pageErr != nil || intLimit < 0 || intPage < 0 {
			return NoPagination
		}

		return PaginationLimits{
			Limit: intLimit,
			Page:  intPage,
		}
	} else {
		// No pagination requested
		return NoPagination
	}
}

// isValidSort checks if the sort param is a known sort param
func isValidSort(sortParam string) bool {
	return sortParam == SORT_NAME ||
		sortParam == SORT_DATE_CREATED ||
		sortParam == SORT_NUM_WORKSPACES ||
		sortParam == SORT_PUBLIC_KEY_ID
}

// isValidSortDirection checks if the sort direction is known
func isValidSortDirection(sortDirection string) bool {
	return sortDirection == SORT_DIR_ASC ||
		sortDirection == SORT_DIR_DESC
}

func ParseSortParams(sortParam string, sortDirection string) SortParams {
	if sortParam != "" &&
		sortDirection != "" &&
		isValidSort(sortParam) &&
		isValidSortDirection(sortDirection) {

		return SortParams{
			SortParam:     sortParam,
			SortDirection: sortDirection,
		}
	} else {
		// No sort specified, use default one
		return DefaultSort
	}
}
