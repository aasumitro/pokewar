package utils_test

import (
	"github.com/aasumitro/pokewar/pkg/utils"
	"testing"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		limit       int
		offset      int
		total       int
		expectedTP  int
		expectedCP  int
		expectedNxt utils.Paging
		expectedPrv utils.Paging
	}{
		{10, 0, 50, 5, 1, utils.Paging{URL: "http://example.com/test?offset=10&limit=10", Path: "test?offset=10&limit=10"}, utils.Paging{}},
		{10, 10, 50, 5, 2, utils.Paging{URL: "http://example.com/test?offset=20&limit=10", Path: "test?offset=20&limit=10"}, utils.Paging{URL: "http://example.com/test?offset=0&limit=10", Path: "test?offset=0&limit=10"}},
		{10, 20, 50, 5, 3, utils.Paging{URL: "http://example.com/test?offset=30&limit=10", Path: "test?offset=30&limit=10"}, utils.Paging{URL: "http://example.com/test?offset=10&limit=10", Path: "test?offset=10&limit=10"}},
		{10, 30, 50, 5, 4, utils.Paging{URL: "http://example.com/test?offset=40&limit=10", Path: "test?offset=40&limit=10"}, utils.Paging{URL: "http://example.com/test?offset=20&limit=10", Path: "test?offset=20&limit=10"}},
		{10, 40, 50, 5, 5, utils.Paging{}, utils.Paging{URL: "http://example.com/test?offset=30&limit=10", Path: "test?offset=30&limit=10"}},
	}
	for _, test := range tests {
		totalPages, currentPage, next, prev := utils.Paginate(test.limit, test.offset, test.total, "example.com", "test")
		if totalPages != test.expectedTP {
			t.Errorf("Paginate(%d, %d, %d): expected totalPages %d, got %d", test.limit, test.offset, test.total, test.expectedTP, totalPages)
		}
		if currentPage != test.expectedCP {
			t.Errorf("Paginate(%d, %d, %d): expected currentPage %d, got %d", test.limit, test.offset, test.total, test.expectedCP, currentPage)
		}
		if next != test.expectedNxt {
			t.Errorf("Paginate(%d, %d, %d): expected next %v, got %v", test.limit, test.offset, test.total, test.expectedNxt, next)
		}
		if prev != test.expectedPrv {
			t.Errorf("Paginate(%d, %d, %d): expected prev %v, got %v", test.limit, test.offset, test.total, test.expectedPrv, prev)
		}
	}
}
