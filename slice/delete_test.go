package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete(t *testing.T) {
	testCases := []struct {
		name      string
		inSlice   []int
		index     int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "index 0",
			inSlice:   []int{0, 1, 2},
			index:     0,
			wantSlice: []int{1, 2},
			wantErr:   nil,
		},
		{
			name:      "index last",
			inSlice:   []int{0, 1, 2},
			index:     2,
			wantSlice: []int{0, 1},
			wantErr:   nil,
		},
		{
			name:      "index -1",
			inSlice:   []int{0, 1, 2},
			index:     -1,
			wantSlice: nil,
			wantErr:   ErrIndexOutOfRange,
		},
		{
			name:      "index 3",
			inSlice:   []int{0, 1, 2},
			index:     3,
			wantSlice: nil,
			wantErr:   ErrIndexOutOfRange,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Delete(tc.inSlice, tc.index)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantSlice, res)
		})
	}
}
