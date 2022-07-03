package registry_test

import (
	"testing"

	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestGetImageRegistryByHostname(t *testing.T) {
	testCases := map[string]struct {
		input        string
		expect       string
		expectErrMsg string
	}{
		"valid": {
			input:  "asia.gcr.io",
			expect: reg.GoogleContainerRegistry,
		},
		"unknown hostname": {
			input:        "some.where.io",
			expectErrMsg: "unknown matcher registry handler by host some.where.io",
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			gotType, err := reg.GetImageRegistryByHostname(tc.input)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			} else {
				assert.Equal(t, tc.expect, gotType)
			}
		})
	}
}
