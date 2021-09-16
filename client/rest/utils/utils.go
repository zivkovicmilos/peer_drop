package utils

import "strconv"

type PaginationLimits struct {
	Limit int
	Page  int
}

var (
	NoPagination = PaginationLimits{
		Limit: -1,
		Page:  -1,
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
