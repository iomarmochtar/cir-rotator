package config_test

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	c "github.com/iomarmochtar/cir-rotator/app/config"
	"github.com/stretchr/testify/assert"
)

func dummyWriter(fileNamePtrn string, content []byte, perms fs.FileMode) (path string, err error) {
	file, err := os.CreateTemp("", fileNamePtrn)
	if err != nil {
		return "", err
	}
	path = file.Name()
	if err = os.WriteFile(path, content, perms); err != nil {
		return "", nil
	}
	if err = os.Chmod(path, perms); err != nil {
		return "", err
	}
	return path, nil
}

func TestConfig_Init(t *testing.T) {
	type tcArg struct {
		config         *c.Config
		beforeExec     func(*tcArg) error
		afterExec      func(*testing.T, *c.Config)
		expectedErrMsg string
	}

	testCases := map[string]*tcArg{
		"registry host is not set": {
			config: &c.Config{
				RegUsername: "user",
				RegPassword: "secret",
			},
			expectedErrMsg: "registry host is required",
		},
		"service account path is set but not exists": {
			config: &c.Config{
				RegistryHost:       "asia.gcr.io",
				ServiceAccountPath: "/tmp/not_found.json",
			},
			expectedErrMsg: "open /tmp/not_found.json: no such file or directory",
		},
		"an error while reading service account file": {
			config: &c.Config{
				RegistryHost: "asia.gcr.io",
			},
			beforeExec: func(tc *tcArg) error {
				path, err := dummyWriter("should-be-err", []byte(`{"hello": "world"}`), 000)
				if err != nil {
					return err
				}
				tc.config.ServiceAccountPath = path
				tc.expectedErrMsg = fmt.Sprintf("open %s: permission denied", path)
				return nil
			},
		},
		"unknown registry type": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "reg.some.where",
				RegistryType: "custreg",
			},
			expectedErrMsg: "unknown image registry type custreg",
		},
		"not specified the parent repo in gcr type": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io",
			},
			expectedErrMsg: "you must specified parent repository after registry host, eg; asia.gcr.io/parent_repo",
		},
		"minimum configs for gcr type are set": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io/parent",
			},
		},
		"an error while reading skip list file": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io/parent",
			},
			beforeExec: func(tc *tcArg) error {
				path, err := dummyWriter("should-be-err", []byte(`"image-1:latest\nimage-2:latest"`), 000)
				if err != nil {
					return err
				}
				tc.config.SkipListPath = path
				tc.expectedErrMsg = fmt.Sprintf("open %s: permission denied", path)
				return nil
			},
		},
		"error while init include filter": {
			config: &c.Config{
				RegUsername:    "user",
				RegPassword:    "secret",
				RegistryHost:   "asia.gcr.io/parent",
				IncludeFilters: []string{"whoami"},
			},
			expectedErrMsg: "unknown name whoami (1:2)\n | (whoami)\n | .^",
		},
		"error while init exclude filter": {
			config: &c.Config{
				RegUsername:    "user",
				RegPassword:    "secret",
				RegistryHost:   "asia.gcr.io/parent",
				ExcludeFilters: []string{"whoami"},
			},
			expectedErrMsg: "unknown name whoami (1:2)\n | (whoami)\n | .^",
		},
		"valid with filters": {
			config: &c.Config{
				RegUsername:    "user",
				RegPassword:    "secret",
				RegistryHost:   "asia.gcr.io/parent",
				ExcludeFilters: []string{"Repository matches 'asia.gcr.io/parent/sub1.*'"},
				IncludeFilters: []string{"Now() + Duration('7d') > CreatedAt"},
			},
		},
		"error reading repo list file": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io/parent",
			},
			beforeExec: func(tc *tcArg) error {
				path, err := dummyWriter("should-be-err", []byte(`"dummy"`), 000)
				if err != nil {
					return err
				}
				tc.config.RepoListPath = path
				tc.expectedErrMsg = fmt.Sprintf("error while reading repository list file: open %s: permission denied", path)
				return nil
			},
		},
		"error unmarshall repository json file": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io/parent",
			},
			beforeExec: func(tc *tcArg) error {
				path, err := dummyWriter("should-be-err", []byte(`"dummy"`), os.ModePerm)
				if err != nil {
					return err
				}
				tc.config.RepoListPath = path
				tc.expectedErrMsg = "unmarshaling repository list file: json: cannot unmarshal string into Go value of type []registry.Repository"
				return nil
			},
		},
		"success read repository json file": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io/parent",
			},
			beforeExec: func(tc *tcArg) error {
				path, err := dummyWriter("valid", []byte(`[{"repository": "asia.gcr.io/parent/ok"}]`), os.ModePerm)
				if err != nil {
					return err
				}
				tc.config.RepoListPath = path
				return nil
			},
			afterExec: func(t *testing.T, c *c.Config) {
				assert.Equal(t, "asia.gcr.io/parent/ok", c.RepositoryList()[0].Name)
			},
		},
		"read skip list": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io/parent",
				SkipListPath: "../../testdata/skip_list.txt",
			},
			afterExec: func(t *testing.T, c *c.Config) {
				assert.Equal(t, []string{"asia.gcr.io/parent1/repo1:latest", "asia.gcr.io/parent1/repo2:release-abc"}, c.SkipList())
			},
		},
		"http worker count": {
			config: &c.Config{
				RegUsername:  "user",
				RegPassword:  "secret",
				RegistryHost: "asia.gcr.io/parent",
				WorkerCount:  10,
			},
			afterExec: func(t *testing.T, c *c.Config) {
				assert.Equal(t, 10, c.HTTPWorkerCount())
			},
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(r *testing.T) {
			if tc.beforeExec != nil {
				assert.NoError(t, tc.beforeExec(tc))
			}
			err := tc.config.Init()
			if tc.expectedErrMsg != "" {
				assert.EqualError(t, err, tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}
			if tc.afterExec != nil {
				tc.afterExec(t, tc.config)
			}
		})
	}
}
