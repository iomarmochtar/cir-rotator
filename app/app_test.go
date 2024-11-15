package app_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/iomarmochtar/cir-rotator/app"
	mc "github.com/iomarmochtar/cir-rotator/app/config/mock_config"
	fl "github.com/iomarmochtar/cir-rotator/pkg/filter"
	mf "github.com/iomarmochtar/cir-rotator/pkg/filter/mock_filter"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	mr "github.com/iomarmochtar/cir-rotator/pkg/registry/mock_registry"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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

func TestApp_ListRepositories(t *testing.T) {
	testCases := map[string]struct {
		mockConfig         func(*gomock.Controller) *mc.MockIConfig
		mockImageRegistry  func(*mr.MockImageRegistry)
		mockIncludeFilter  func(*gomock.Controller) fl.IFilterEngine
		mockExcludeFilter  func(*gomock.Controller) fl.IFilterEngine
		expectErrMsg       string
		expectRepositories []reg.Repository
	}{
		"error while get repository catalog": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Catalog().Times(1).Return(nil, fmt.Errorf("an error while fetching catalog"))

				mockConfig := mc.NewMockIConfig(ctrl)

				mockConfig.EXPECT().ImageRegistry().Return(mockReg)
				mockConfig.EXPECT().RepositoryList().Times(1).Return([]reg.Repository{})
				mockConfig.EXPECT().IncludeEngine().Times(0)
				mockConfig.EXPECT().ExcludeEngine().Times(0)
				return mockConfig
			},
			expectErrMsg:       "an error while fetching catalog",
			expectRepositories: nil,
		},
		"error in include filter": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)

				mif := mf.NewMockIFilterEngine(ctrl)
				mif.EXPECT().Process(gomock.Any()).Times(1).Return(false, fmt.Errorf("error in include filter"))

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().RepositoryList().Times(1).Return([]reg.Repository{})
				mockConfig.EXPECT().ImageRegistry().Times(1).Return(mockReg)
				mockConfig.EXPECT().IncludeEngine().Times(1).Return(mif)
				mockConfig.EXPECT().ExcludeEngine().Times(1).Return(nil)
				return mockConfig
			},
			expectErrMsg:       "error in include filter",
			expectRepositories: nil,
		},
		"error in exclude filter": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)

				mif := mf.NewMockIFilterEngine(ctrl)
				mif.EXPECT().Process(gomock.Any()).Times(1).Return(true, nil)

				mef := mf.NewMockIFilterEngine(ctrl)
				mef.EXPECT().Process(gomock.Any()).Times(1).Return(false, fmt.Errorf("error in exclude filter"))

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().RepositoryList().Times(1).Return([]reg.Repository{})
				mockConfig.EXPECT().ImageRegistry().Times(1).Return(mockReg)
				mockConfig.EXPECT().IncludeEngine().Times(1).Return(mif)
				mockConfig.EXPECT().ExcludeEngine().Times(1).Return(mef)
				return mockConfig
			},
			expectErrMsg:       "error in exclude filter",
			expectRepositories: nil,
		},
		"if filters are not provided then will returning all results": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().RepositoryList().Times(1).Return([]reg.Repository{})
				mockConfig.EXPECT().ImageRegistry().Times(1).Return(mockReg)
				mockConfig.EXPECT().IncludeEngine().Times(1).Return(nil)
				mockConfig.EXPECT().ExcludeEngine().Times(1).Return(nil)
				return mockConfig
			},
			expectRepositories: sampleRepos,
		},
		"if match with exclude filter then it will not passed as result": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)

				mif := mf.NewMockIFilterEngine(ctrl)
				mif.EXPECT().Process(gomock.Any()).Times(2).Return(true, nil)

				mef := mf.NewMockIFilterEngine(ctrl)
				mef.EXPECT().Process(gomock.Any()).AnyTimes().DoAndReturn(func(arg fl.Fields) (bool, error) {
					if arg.Digest == sampleRepos[1].Digests[0].Name {
						return true, nil
					}
					return false, nil
				})

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().RepositoryList().Times(1).Return([]reg.Repository{})
				mockConfig.EXPECT().ImageRegistry().Times(1).Return(mockReg)
				mockConfig.EXPECT().IncludeEngine().Times(1).Return(mif)
				mockConfig.EXPECT().ExcludeEngine().Times(1).Return(mef)
				return mockConfig
			},
			expectRepositories: []reg.Repository{sampleRepos[0]},
		},
		"include filter is provided but none match with it": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)

				mif := mf.NewMockIFilterEngine(ctrl)
				mif.EXPECT().Process(gomock.Any()).AnyTimes().Return(false, nil)

				mef := mf.NewMockIFilterEngine(ctrl)
				mef.EXPECT().Process(gomock.Any()).Times(0)

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().RepositoryList().Times(1).Return([]reg.Repository{})
				mockConfig.EXPECT().ImageRegistry().Times(1).Return(mockReg)
				mockConfig.EXPECT().IncludeEngine().Times(1).Return(mif)
				mockConfig.EXPECT().ExcludeEngine().Times(1).Return(mef)
				return mockConfig
			},
			expectRepositories: nil,
		},
		"exclude filter is provided but no one match with it": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Catalog().Times(1).Return(sampleRepos, nil)

				mef := mf.NewMockIFilterEngine(ctrl)
				mef.EXPECT().Process(gomock.Any()).AnyTimes().Return(false, nil)

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().RepositoryList().Times(1).Return([]reg.Repository{})
				mockConfig.EXPECT().ImageRegistry().Times(1).Return(mockReg)
				mockConfig.EXPECT().IncludeEngine().Times(1).Return(nil)
				mockConfig.EXPECT().ExcludeEngine().Times(1).Return(mef)
				return mockConfig
			},
			expectRepositories: sampleRepos,
		},
		"repository list provided in config initialization": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().RepositoryList().Times(1).Return(sampleRepos)
				return mockConfig
			},
			expectRepositories: sampleRepos,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConfig := tc.mockConfig(ctrl)

			repositories, err := app.New(mockConfig).ListRepositories()
			assert.Equal(t, tc.expectRepositories, repositories)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			}
		})
	}
}

func TestApp_DeleteRepositories(t *testing.T) {
	// construct sample data from sample repo with some additions to match with the test case
	var repoWithMoreDigest []reg.Repository
	repoWithMoreDigest = append(repoWithMoreDigest, sampleRepos...)
	deleteRepoDigest := reg.Digest{
		Name:           "sha256:01551c49819f8bda0a8bdc6216e5793404b0adb4937d407e99a590c0c5cb8078",
		ImageSizeBytes: 488119934,
		Tag:            []string{"release-prod-abc"},
		Created:        time.Date(2022, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
		Uploaded:       time.Date(2022, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
	}
	repoWithMoreDigest[0].Digests = append(repoWithMoreDigest[0].Digests, deleteRepoDigest)
	repoWithMoreDigest = append(repoWithMoreDigest, reg.Repository{
		Name: "image-3",
		Digests: []reg.Digest{
			{
				Name:           "sha256:B0ac9df37ff356753cd20f4475d4b8d3a543b4d45db2390c0275be2ee7a0xyz",
				ImageSizeBytes: 488118834,
				Tag:            []string{"abc", "def"},
				Created:        time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
				Uploaded:       time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			},
		},
	})

	testCases := map[string]struct {
		mockConfig   func(*gomock.Controller) *mc.MockIConfig
		repositories []reg.Repository
		expectErrMsg string
	}{
		"got an error while deleting image": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Delete(sampleRepos[0]).Times(1).Return(fmt.Errorf("failure"))

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().ImageRegistry().Times(1).Return(mockReg)
				mockConfig.EXPECT().SkipList().Times(1).Return([]string{})
				mockConfig.EXPECT().IsDryRun().Times(2).Return(false)
				mockConfig.EXPECT().HTTPWorkerCount().Times(1).Return(1)
				mockConfig.EXPECT().SkipDeletionErr().Times(1).Return(false)
				return mockConfig
			},
			repositories: sampleRepos,
			expectErrMsg: "error while deleting repository image-1: failure",
		},
		"will not returning any error if skip-error provided": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				mockReg.EXPECT().Delete(gomock.Any()).Times(2).Return(fmt.Errorf("failure"))

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().ImageRegistry().Times(2).Return(mockReg)
				mockConfig.EXPECT().SkipList().Return([]string{})
				mockConfig.EXPECT().IsDryRun().Times(2).Return(false)
				mockConfig.EXPECT().HTTPWorkerCount().Times(1).Return(1)
				mockConfig.EXPECT().SkipDeletionErr().Times(2).Return(true)
				return mockConfig
			},
			repositories: sampleRepos,
		},
		"found in skip list": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockReg := mr.NewMockImageRegistry(ctrl)
				expetedDeleteRepoImage1 := reg.Repository{Name: "image-1", Digests: []reg.Digest{deleteRepoDigest}}
				mockReg.EXPECT().Delete(repoWithMoreDigest[1]).Times(1).Return(nil)
				mockReg.EXPECT().Delete(expetedDeleteRepoImage1).Times(1).Return(nil)

				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().ImageRegistry().Times(2).Return(mockReg)
				mockConfig.EXPECT().SkipList().Times(1).Return([]string{"image-1:latest", "image-3:abc", "image-3:def"})
				mockConfig.EXPECT().IsDryRun().Times(2).Return(false)
				mockConfig.EXPECT().HTTPWorkerCount().Times(1).Return(1)
				return mockConfig
			},
			repositories: repoWithMoreDigest,
		},
		"will not calling registry delete api when dry run is set": {
			mockConfig: func(ctrl *gomock.Controller) *mc.MockIConfig {
				mockConfig := mc.NewMockIConfig(ctrl)
				mockConfig.EXPECT().ImageRegistry().Times(0)
				mockConfig.EXPECT().SkipList().Times(1).Return([]string{})
				mockConfig.EXPECT().IsDryRun().Times(3).Return(true)
				mockConfig.EXPECT().HTTPWorkerCount().Times(1).Return(1)
				return mockConfig
			},
			repositories: repoWithMoreDigest,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConfig := tc.mockConfig(ctrl)
			err := app.New(mockConfig).DeleteRepositories(tc.repositories)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
