package registry_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	hl "github.com/iomarmochtar/cir-rotator/pkg/helpers"
	mh "github.com/iomarmochtar/cir-rotator/pkg/http/mock_http"
	"github.com/iomarmochtar/cir-rotator/pkg/registry"
	reg "github.com/iomarmochtar/cir-rotator/pkg/registry"
	"github.com/stretchr/testify/assert"
)

var (
	gcrHost      = "asia.gcr.io"
	gcrHostHTTPS = fmt.Sprintf("https://%s", gcrHost)
)

func readFixture(loc string) []byte {
	data, err := ioutil.ReadFile(path.Join("..", "..", "testdata", loc))
	if err != nil {
		panic(errors.Wrap(err, "while reading fixtures data"))
	}
	return data
}

func readGCRResponseFixture[t any](loc string) t {
	var obj t
	data := readFixture(hl.SlashJoin("gcr", loc))
	if err := json.Unmarshal(data, &obj); err != nil {
		panic(errors.Wrapf(err, "while reading gcr fixture file %s", loc))
	}
	return obj
}

func TestNewGCR(t *testing.T) {
	testCases := map[string]struct {
		host        string
		expectedErr bool
	}{
		"host not set by base repository path": {
			host:        "asia.gcr.io",
			expectedErr: true,
		},
		"base repository is provided": {
			host:        "asia.gcr.io/parent_repo",
			expectedErr: false,
		},
		"sub base repository": {
			host:        "asia.gcr.io/parent_repo/sub1",
			expectedErr: false,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			imgReg, err := reg.NewGCR(tc.host, nil)
			if tc.expectedErr {
				assert.Nil(t, imgReg)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, imgReg)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGCR_Catalog(t *testing.T) {
	parentRepo := "parent"
	hostWithParentRepo := hl.SlashJoin(gcrHost, parentRepo)
	parentTagListURL := hl.SlashJoin(gcrHostHTTPS, "v2", parentRepo, "tags", "list")
	paretRepoResp := readGCRResponseFixture[reg.GCRTagsResponse]("tag_list_parent.json")
	sub1RepoResp := readGCRResponseFixture[reg.GCRTagsResponse]("tag_list_sub1.json")

	testCases := map[string]struct {
		mockHTTPClient     func(*mh.MockIHttpClient)
		expectErrMsg       string
		expectRepositories []reg.Repository
	}{
		"error while get tags from parent repository": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).Return(fmt.Errorf("an error while fetching catalog"))
			},
			expectErrMsg:       "an error while fetching catalog",
			expectRepositories: nil,
		},
		"error in response body": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					jbody.Errors = []reg.ErrorField{
						{
							Code:    "UNKNOWN",
							Message: "error message goes here",
						},
					}
					return nil
				})
			},
			expectErrMsg:       "[UNKNOWN] [error message goes here]",
			expectRepositories: nil,
		},
		"got an error while access child repo": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				childRepo := "child1"
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					jbody.Child = []string{childRepo}
					return nil
				})

				childURL := hl.SlashJoin(gcrHostHTTPS, "v2", parentRepo, childRepo, "tags", "list")
				m.EXPECT().GetMarshalReturnObj(childURL, gomock.Any()).Times(1).Return(fmt.Errorf("an error in accessing child repo"))
			},
			expectRepositories: nil,
			expectErrMsg:       "an error in accessing child repo",
		},
		"got invalid created timestamp": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					*jbody = readGCRResponseFixture[reg.GCRTagsResponse]("tag_list_parent_invalid_created_time.json")
					*jbody = readGCRResponseFixture[reg.GCRTagsResponse]("tag_list_parent_invalid_created_time.json")
					return nil
				})
			},
			expectRepositories: nil,
			expectErrMsg:       `while parse created time: strconv.ParseInt: parsing "": invalid syntax`,
		},
		"got invalid uploaded timestamp": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					*jbody = readGCRResponseFixture[reg.GCRTagsResponse]("tag_list_parent_invalid_uploaded_time.json")
					return nil
				})
			},
			expectRepositories: nil,
			expectErrMsg:       `while parse uploaded time: strconv.ParseInt: parsing "": invalid syntax`,
		},
		"invalid image size byte": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					*jbody = readGCRResponseFixture[reg.GCRTagsResponse]("tag_list_parent_invalid_byte_size.json")
					return nil
				})
			},
			expectRepositories: nil,
			expectErrMsg:       `while converting image size: strconv.ParseUint: parsing "abc": invalid syntax`,
		},
		"recursively access child url that has digest information if any": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				childRepo1 := "sub1"
				childRepo2 := "sub2"
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					*jbody = paretRepoResp
					jbody.Child = []string{childRepo1, childRepo2}
					return nil
				})

				childURL := hl.SlashJoin(gcrHostHTTPS, "v2", parentRepo, childRepo1, "tags", "list")
				m.EXPECT().GetMarshalReturnObj(childURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					*jbody = sub1RepoResp
					return nil
				})

				// no digest but has several child
				childURL2 := hl.SlashJoin(gcrHostHTTPS, "v2", parentRepo, childRepo2)
				m.EXPECT().GetMarshalReturnObj(hl.SlashJoin(childURL2, "tags", "list"), gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					*jbody = readGCRResponseFixture[reg.GCRTagsResponse]("empty_repo_with_child.json")
					return nil
				})

				m.EXPECT().GetMarshalReturnObj(hl.SlashJoin(childURL2, "cronjob-image", "tags", "list"), gomock.Any()).Times(1).Return(nil)
				m.EXPECT().GetMarshalReturnObj(hl.SlashJoin(childURL2, "job-script", "tags", "list"), gomock.Any()).Times(1).Return(nil)
			},
			expectRepositories: []reg.Repository{
				{
					Name: "asia.gcr.io/parent/sub1",
					Digests: []reg.Digest{
						{
							Name:           "sha256:C05ce64163cd2327d364933df75aa4850af425b6cbaec2f6af3b31e5246be0e2",
							ImageSizeBytes: 530325786,
							Tag:            []string{"latest", "abc"},
							Created:        time.Unix(1585800237411/1000, 0),
							Uploaded:       time.Unix(1585800278141/1000, 0),
						},
					},
				},
				{
					Name: "asia.gcr.io/parent",
					Digests: []reg.Digest{
						{
							Name:           "sha256:02123554c8d65d241a77dd8e238403ba8cc697afca7208bbe7564b634aa22fee",
							ImageSizeBytes: 365557176,
							Tag:            []string{"latest", "release-20210624-150000"},
							Created:        time.Unix(1624518709150/1000, 0),
							Uploaded:       time.Unix(1624518770462/1000, 0),
						},
					},
				},
			},
		},
		"empty repository will not be marked as result": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				m.EXPECT().GetMarshalReturnObj(parentTagListURL, gomock.Any()).Times(1).DoAndReturn(func(url string, jbody *reg.GCRTagsResponse) error {
					*jbody = readGCRResponseFixture[reg.GCRTagsResponse]("empty_repo_with_child.json")
					return nil
				})

				child1Url := hl.SlashJoin(gcrHostHTTPS, "v2", parentRepo, "cronjob-image", "tags", "list")
				m.EXPECT().GetMarshalReturnObj(child1Url, gomock.Any()).Times(1).Return(nil)

				child2Url := hl.SlashJoin(gcrHostHTTPS, "v2", parentRepo, "job-script", "tags", "list")
				m.EXPECT().GetMarshalReturnObj(child2Url, gomock.Any()).Times(1).Return(nil)
			},
			expectRepositories: nil,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mHc := mh.NewMockIHttpClient(ctrl)
			tc.mockHTTPClient(mHc)

			gcr, _ := registry.NewGCR(hostWithParentRepo, mHc)
			repositories, err := gcr.Catalog()

			assert.Equal(t, tc.expectRepositories, repositories)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGCR_Delete(t *testing.T) {
	sampleRepo := reg.Repository{
		Name: "asia.gcr.io/parent/sub1",
		Digests: []reg.Digest{
			{
				Name:           "sha256:C05ce64163cd2327d364933df75aa4850af425b6cbaec2f6af3b31e5246be0e2",
				ImageSizeBytes: 530325786,
				Tag:            []string{"latest", "abc"},
				Created:        time.Date(2020, time.April, 2, 11, 3, 57, 0, time.Local),
				Uploaded:       time.Date(2020, time.April, 2, 11, 4, 38, 0, time.Local),
			},
		},
	}

	testCases := map[string]struct {
		mockHTTPClient func(*mh.MockIHttpClient)
		isDryRun       bool
		repository     reg.Repository
		expectErrMsg   string
	}{
		"will not do api call for deletion if dry run is true": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				m.EXPECT().DeleteMarshalReturnObj(gomock.Any(), gomock.Any()).Times(0)
			},
			repository: sampleRepo,
			isDryRun:   true,
		},
		"will call api for deletion if dry run is false": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				manifestURL := hl.SlashJoin(gcrHostHTTPS, "v2", "parent", "sub1", "manifests")
				// expecting call in order since the tags will be deleted first before it's digest
				gomock.InOrder(
					m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "latest"), gomock.Any()).Times(1).Return(nil),
					m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "abc"), gomock.Any()).Times(1).Return(nil),
					m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "sha256:C05ce64163cd2327d364933df75aa4850af425b6cbaec2f6af3b31e5246be0e2"), gomock.Any()).Times(1).Return(nil),
				)
			},
			repository: sampleRepo,
			isDryRun:   false,
		},
		"an error while deleting tag": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				manifestURL := hl.SlashJoin(gcrHostHTTPS, "v2", "parent", "sub1", "manifests")
				m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "latest"), gomock.Any()).Times(1).Return(fmt.Errorf("an error while deleting tag"))
				m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "abc"), gomock.Any()).Times(0)
				m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "sha256:C05ce64163cd2327d364933df75aa4850af425b6cbaec2f6af3b31e5246be0e2"), gomock.Any()).Times(0)
			},
			repository:   sampleRepo,
			isDryRun:     false,
			expectErrMsg: "an error while deleting tag",
		},
		"an error while deleting digest": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				manifestURL := hl.SlashJoin(gcrHostHTTPS, "v2", "parent", "sub1", "manifests")
				gomock.InOrder(
					m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "latest"), gomock.Any()).Times(1).Return(nil),
					m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "abc"), gomock.Any()).Times(1).Return(nil),
					m.EXPECT().DeleteMarshalReturnObj(hl.SlashJoin(manifestURL, "sha256:C05ce64163cd2327d364933df75aa4850af425b6cbaec2f6af3b31e5246be0e2"), gomock.Any()).Times(1).Return(fmt.Errorf("an error in manifest deletion")),
				)
			},
			repository:   sampleRepo,
			isDryRun:     false,
			expectErrMsg: "an error in manifest deletion",
		},
		"error from delete response": {
			mockHTTPClient: func(m *mh.MockIHttpClient) {
				latestTagURL := hl.SlashJoin(gcrHostHTTPS, "v2", "parent", "sub1", "manifests", "latest")
				m.EXPECT().DeleteMarshalReturnObj(latestTagURL, gomock.Any()).Times(1).DoAndReturn(func(url string, errFields *reg.ErrorsField) error {
					*errFields = readGCRResponseFixture[reg.ErrorsField]("error_delete_manifest.json")
					return nil
				})
			},
			repository:   sampleRepo,
			isDryRun:     false,
			expectErrMsg: `Failed to compute blob liveness for manifest: 'latest'`,
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mHc := mh.NewMockIHttpClient(ctrl)
			tc.mockHTTPClient(mHc)

			gcr, err := registry.NewGCR(hl.SlashJoin(gcrHost, "parent"), mHc)
			assert.NoError(t, err)

			err = gcr.Delete(tc.repository, tc.isDryRun)
			if tc.expectErrMsg != "" {
				assert.EqualError(t, err, tc.expectErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
