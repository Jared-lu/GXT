package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		inSlice   []int
		val       int
		index     int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "index 0",
			inSlice:   []int{0, 1, 2},
			val:       10,
			index:     0,
			wantSlice: []int{10, 0, 1, 2},
			wantErr:   nil,
		},
		{
			name:      "index last",
			inSlice:   []int{0, 1, 2},
			val:       10,
			index:     3,
			wantSlice: []int{0, 1, 2, 10},
			wantErr:   nil,
		},
		{
			name:      "index -1",
			inSlice:   []int{0, 1, 2},
			val:       10,
			index:     -1,
			wantSlice: nil,
			wantErr:   ErrIndexOutOfRange,
		},
		{
			name:      "index 4",
			inSlice:   []int{0, 1, 2},
			val:       10,
			index:     4,
			wantSlice: nil,
			wantErr:   ErrIndexOutOfRange,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Add(tc.inSlice, tc.val, tc.index)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantSlice, res)
		})
	}
}
