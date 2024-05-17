package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	testCases := []struct {
		name      string
		inSlice   []int
		wantRes   int
		wantMatch bool
	}{
		{
			name:      "match",
			inSlice:   []int{1, 2, 30},
			wantRes:   30,
			wantMatch: true,
		},
		{
			name:      "not match",
			inSlice:   []int{1, 2, 3},
			wantRes:   0,
			wantMatch: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, match := Find(tc.inSlice, func(val int) bool {
				if val > 10 {
					return true
				} else {
					return false
				}
			})
			assert.Equal(t, tc.wantMatch, match)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestFindAll(t *testing.T) {
	testCases := []struct {
		name    string
		inSlice []int
		wantRes []int
	}{
		{
			name:    "match",
			inSlice: []int{1, 2, 30, 40},
			wantRes: []int{30, 40},
		},
		{
			name:    "not match",
			inSlice: []int{1, 2, 3, 4},
			wantRes: []int{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := FindAll(tc.inSlice, func(val int) bool {
				if val > 10 {
					return true
				}
				return false
			})
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
