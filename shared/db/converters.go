package db

import (
	"fmt"
	"strconv"
	"strings"
)

func IntSliceToString(nums []int) string {
	if len(nums) == 0 {
		return ""
	}

	strs := make([]string, len(nums))
	for i, n := range nums {
		strs[i] = strconv.Itoa(n)
	}

	return strings.Join(strs, ",")
}

func StringToIntSlice(s string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}

	parts := strings.Split(s, ",")
	result := make([]int, len(parts))

	for i, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, fmt.Errorf("invalid number at index %d: %w", i, err)
		}
		result[i] = n
	}

	return result, nil
}
