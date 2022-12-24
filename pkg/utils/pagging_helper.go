package utils

import (
	"fmt"
)

func Paginate(limit, offset, total int, host, path string) (totalPages, currentPage int, next, prev Paging) {
	totalPages = total / limit
	pageLeft := ((total - offset) / limit) - 1
	currentPage = (total / limit) - pageLeft

	if currentPage > 1 {
		prevOffset := offset - limit
		prevPath := fmt.Sprintf("%s?offset=%d&limit=%d", path, prevOffset, limit)
		prevLink := fmt.Sprintf("%s://%s/%s", "http", host, prevPath)
		prev = Paging{
			URL:  prevLink,
			Path: prevPath,
		}
	}

	if totalPages != currentPage && totalPages != 0 {
		nextOffset := offset + limit
		nextPath := fmt.Sprintf("%s?offset=%d&limit=%d", path, nextOffset, limit)
		nextLink := fmt.Sprintf("%s://%s/%s", "http", host, nextPath)
		next = Paging{
			URL:  nextLink,
			Path: nextPath,
		}
	}

	return
}
