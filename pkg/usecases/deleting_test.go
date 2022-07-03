package usecases_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	mr "github.com/iomarmochtar/cir-rotator/pkg/registry/mock_registry"
	"github.com/iomarmochtar/cir-rotator/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRepositories(t *testing.T) {
	testCases := map[string]struct {
		mockImageRegistry func(*mr.MockImageRegistry)
		repositories      []reg.Repository
		skipList          []string
		dryRun            bool
		expectErrMsg      string
	}{
		"got an error while deleting image": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Delete(sampleRepos[0], false).Times(1).Return(fmt.Errorf("failure"))
			},
			dryRun:       false,
			repositories: sampleRepos,
			expectErrMsg: "failure",
		},
		"found in skip list": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Delete(sampleRepos[1], false).Times(1).Return(nil)
			},
			skipList:     []string{"image-1:latest"},
			dryRun:       false,
			repositories: sampleRepos,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			imgReg := mr.NewMockImageRegistry(ctrl)
			tc.mockImageRegistry(imgReg)

			err := usecases.DeleteRepositories(imgReg, tc.repositories, tc.skipList, tc.dryRun)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
