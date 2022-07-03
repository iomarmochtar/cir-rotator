package usecases_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	fl "github.com/iomarmochtar/cir-rotator/pkg/filter"
	mf "github.com/iomarmochtar/cir-rotator/pkg/filter/mock_filter"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	mr "github.com/iomarmochtar/cir-rotator/pkg/registry/mock_registry"
	"github.com/iomarmochtar/cir-rotator/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

var (
	sampleRepos = []reg.Repository{
		{
			Name: "image-1",
			Digests: []reg.Digest{
				{
					Name:           "sha256:B0ac9df37ff356753cd20f4475d4b8d3a543b4d45db2390c0275be2ee7a09b2e",
					ImageSizeBytes: 488118834,
					Tag:            []string{"latest", "release-abc-def"},
					Created:        time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
					Uploaded:       time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
				},
			},
		},
		{
			Name: "image-2",
			Digests: []reg.Digest{
				{
					Name:           "sha256:C05ce64163cd2327d364933df75aa4850af425b6cbaec2f6af3b31e5246be0e2",
					ImageSizeBytes: 530325786,
					Tag:            []string{"abc", "xyz"},
					Created:        time.Date(2022, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
					Uploaded:       time.Date(2022, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
				},
			},
		},
	}
)

func TestListRepositories(t *testing.T) {
	testCases := map[string]struct {
		mockImageRegistry  func(*mr.MockImageRegistry)
		mockIncludeFilter  func(*gomock.Controller) fl.IFilterEngine
		mockExcludeFilter  func(*gomock.Controller) fl.IFilterEngine
		expectErrMsg       string
		expectRepositories []reg.Repository
	}{
		"error while get repository catalog": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Catalog().Times(1).Return(nil, fmt.Errorf("an error while fetching catalog"))
			},
			mockIncludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				return nil
			},
			mockExcludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				return nil
			},
			expectErrMsg:       "an error while fetching catalog",
			expectRepositories: nil,
		},
		"error in include filter": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)
			},
			mockIncludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).Times(1).Return(false, fmt.Errorf("error in include filter"))
				return m
			},
			mockExcludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).Times(0)
				return m
			},
			expectErrMsg:       "error in include filter",
			expectRepositories: nil,
		},
		"error in exclude filter": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)
			},
			mockIncludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).Times(1).Return(true, nil)
				return m
			},
			mockExcludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).Times(1).Return(false, fmt.Errorf("error in exclude filter"))
				return m
			},
			expectErrMsg:       "error in exclude filter",
			expectRepositories: nil,
		},
		"if filters are not provided then will returning all results": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)
			},
			mockIncludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				return nil
			},
			mockExcludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				return nil
			},
			expectRepositories: sampleRepos,
		},
		"if match with exclude filter then it will not passed as result": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)
			},
			mockIncludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).Times(2).Return(true, nil)
				return m
			},
			mockExcludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).AnyTimes().DoAndReturn(func(arg fl.Fields) (bool, error) {
					if arg.Digest == sampleRepos[1].Digests[0].Name {
						return true, nil
					}
					return false, nil
				})
				return m
			},
			expectRepositories: []reg.Repository{sampleRepos[0]},
		},
		"include filter is provided but none match with it": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)
			},
			mockIncludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).AnyTimes().Return(false, nil)
				return m
			},
			mockExcludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).Times(0)
				return m
			},
			expectRepositories: nil,
		},
		"exclude filter is provided but no one match with it": {
			mockImageRegistry: func(m *mr.MockImageRegistry) {
				m.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)
			},
			mockIncludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				return nil
			},
			mockExcludeFilter: func(ctrl *gomock.Controller) fl.IFilterEngine {
				m := mf.NewMockIFilterEngine(ctrl)
				m.EXPECT().Process(gomock.Any()).AnyTimes().Return(false, nil)
				return m
			},
			expectRepositories: sampleRepos,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			imgReg := mr.NewMockImageRegistry(ctrl)
			tc.mockImageRegistry(imgReg)

			iFilter := tc.mockIncludeFilter(ctrl)
			eFilter := tc.mockExcludeFilter(ctrl)

			repositories, err := usecases.ListRepositories(imgReg, iFilter, eFilter)
			assert.Equal(t, tc.expectRepositories, repositories)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			}
		})
	}
}
