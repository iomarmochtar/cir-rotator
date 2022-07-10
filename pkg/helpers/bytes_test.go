package helpers_test

import (
	"testing"

	"github.com/iomarmochtar/cir-rotator/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestByteCountIEC(t *testing.T) {
	assert.Equal(t, "984 B", helpers.ByteCountIEC(984), "under 1024")
	assert.Equal(t, "505.8 MiB", helpers.ByteCountIEC(530325786), "upper 1024")
}

func TestSizeUnitStrToFloat(t *testing.T) {
	testCase := map[string]struct {
		expectedErrMsg string
		result         float64
		input          string
	}{
		"unknown pattern": {
			input:          "1 KB",
			expectedErrMsg: "unknown pattern 1 KB",
		},
		"invalid number": {
			input:          "1.32.32 KiB",
			expectedErrMsg: "unknown pattern 1.32.32 KiB",
		},
		"successfully converted KiB": {
			input:  "1 KiB",
			result: 1024,
		},
		"successfully converted MiB": {
			input:  "1 MiB",
			result: 1024 * 1024,
		},
		"successfully converted GiB": {
			input:  "1 GiB",
			result: 1024 * 1024 * 1024,
		},
	}

	for title, tc := range testCase {
		t.Run(title, func(t *testing.T) {
			result, err := helpers.SizeUnitStrToFloat(tc.input)
			if tc.expectedErrMsg != "" {
				assert.EqualError(t, err, tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.result, result)
		})
	}
}
