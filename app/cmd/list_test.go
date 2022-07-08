package cmd_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/iomarmochtar/cir-rotator/app/cmd"
	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
)

func TestListAction(t *testing.T) {
	testCases := map[string]struct {
		cmdArgs        []string
		expectedErrMsg string
		mockImageReg   func(w http.ResponseWriter) error
	}{
		"not providing any output": {
			cmdArgs:        []string{"-ho", "asia.gcr.io"},
			expectedErrMsg: "must specified one or more output",
		},
		"basic auth authentication params not provided": {
			cmdArgs:        []string{"-ho", "asia.gcr.io", "--output-table"},
			expectedErrMsg: "you must set oauth token or basic auth params (username & password)",
		},
		"successfully listing repository": {
			cmdArgs: []string{"--type", "gcr", "--output-table", "-u", "secret", "-p", "souce"},
			mockImageReg: func(w http.ResponseWriter) error {
				data, err := ioutil.ReadFile(path.Join("..", "..", "testdata", "gcr", "tag_list_no_child.json"))
				if err != nil {
					return err
				}
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(data)
				return err
			},
		},
		"error while listing repository": {
			cmdArgs:        []string{"--type", "gcr", "--output-table", "-u", "secret", "-p", "souce"},
			expectedErrMsg: "invalid character 'o' looking for beginning of value",
			mockImageReg: func(w http.ResponseWriter) error {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("ok"))
				return err
			},
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			if tc.mockImageReg != nil {
				ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					assert.NoError(t, tc.mockImageReg(w))
				}))
				// at moment for supporting gcr pattern
				normHost := strings.Replace(ts.URL, "https://", "", 1)
				os.Setenv("REGISTRY_HOST", fmt.Sprintf("%s/repo", normHost))
				os.Setenv("ALLOW_INSECURE_SSL", "true")

				defer os.Unsetenv("REGISTRY_HOST")
				defer os.Unsetenv("ALLOW_INSECURE_SSL")
				defer ts.Close()
			}
			command := cmd.ListAction()
			set := flag.NewFlagSet("test", 0)
			app := &cli.App{Writer: ioutil.Discard}
			assert.NoError(t, set.Parse(append([]string{"list"}, tc.cmdArgs...)))

			ctx := cli.NewContext(app, set, nil)
			err := command.Run(ctx)

			if tc.expectedErrMsg != "" {
				assert.EqualError(t, err, tc.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
