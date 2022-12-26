package utils_test

import (
	"github.com/aasumitro/pokewar/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInArray(t *testing.T) {
	assert.Equalf(t, true,
		utils.InArray("a", []string{"a", "b", "c"}),
		"InArray(%v, %v)", "a", []string{"a", "b", "c"})

	assert.Equalf(t, false,
		utils.InArray("d", []string{"a", "b", "c"}),
		"InArray(%v, %v)", "a", []string{"a", "b", "c"})

	assert.Equalf(t, false,
		utils.InArray(4, []int{1, 2, 3}),
		"InArray(%v, %v)", "a", []string{"a", "b", "c"})

	assert.Equalf(t, true,
		utils.InArray(1, []int{1, 2, 3}),
		"InArray(%v, %v)", "a", []string{"a", "b", "c"})
}
