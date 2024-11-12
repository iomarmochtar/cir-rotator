package registry_test

import (
	"testing"

	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestGetImageRegistryTypeByHostname(t *testing.T) {
	testCases := map[string]struct {
		input        string
		expect       string
		expectErrMsg string
	}{
		"gcr: valid gcr.io": {
			input:  "asia.gcr.io",
			expect: reg.GoogleContainerRegistry,
		},
		"gcr: valid artifact registry domain": {
			input:  "asia-southeast2-docker.pkg.dev",
			expect: reg.GoogleContainerRegistry,
		},
		"type generic": {
			input:        "some.where.io",
			expectErrMsg: "unknown matcher registry handler by host some.where.io",
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			gotType, err := reg.GetImageRegistryTypeByHostname(tc.input)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			} else {
				assert.Equal(t, tc.expect, gotType)
			}
		})
	}
}
