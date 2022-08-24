package cmd_test

import (
	"bytes"
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
	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
)

type caseParam struct {
	cmdArgs        []string
	expectedErrMsg string
	mockImageReg   func(w http.ResponseWriter, r *http.Request) error
	beforeRunExec  func() error
	afterRunExec   func() error
}

func readFixture(fpath string) []byte {
	data, err := ioutil.ReadFile(path.Join("..", "..", "testdata", fpath))
	if err != nil {
		panic(err)
	}
	return data
}

var (
	commonTestCases = map[string]caseParam{
		"basic auth authentication params not provided": {
			cmdArgs:        []string{"-ho", "asia.gcr.io", "--output-table"},
			expectedErrMsg: "you must set oauth token or basic auth params (username & password)",
		},
		"successfully listing repository": {
			cmdArgs: []string{"--output-table", "-u", "secret", "-p", "souce"},
			mockImageReg: func(w http.ResponseWriter, r *http.Request) error {
				data := readFixture("gcr/tag_list_no_child.json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(data)
				return err
			},
		},
		"error while listing repository": {
			cmdArgs:        []string{"--output-table", "-u", "secret", "-p", "souce"},
			expectedErrMsg: "invalid character 'o' looking for beginning of value",
			mockImageReg: func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("ok"))
				return err
			},
		},
		"output json": {
			cmdArgs: []string{"--output-json", "/tmp/dump_output_path.json", "-u", "user", "-p", "secret"},
			beforeRunExec: func() error {
				// make sure the output file is not exist
				outputPath := "/tmp/dump_output_path.json"
				if h.FileExist(outputPath) {
					return os.Remove(outputPath)
				}
				return nil
			},
			afterRunExec: func() error {
				// make sure the output file is not exist
				outputPath := "/tmp/dump_output_path.json"
				if h.FileExist(outputPath) {
					return nil
				}
				return fmt.Errorf("%s is not found", outputPath)
			},
			mockImageReg: func(w http.ResponseWriter, r *http.Request) error {
				data := readFixture("gcr/tag_list_no_child.json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(data)
				return err
			},
		},
	}
)

func runCmdTestCases(name string, command *cli.Command, testCases map[string]caseParam, t *testing.T) {
	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			if tc.beforeRunExec != nil {
				assert.NoError(t, tc.beforeRunExec())
			}

			if tc.mockImageReg != nil {
				ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NoError(t, tc.mockImageReg(w, r))
				}))
				// at moment for supporting gcr pattern
				normHost := fmt.Sprintf("%s/repo", strings.Replace(ts.URL, "https://", "", 1))
				tc.cmdArgs = append(tc.cmdArgs, "-ho", normHost, "--allow-insecure", "--type", "gcr")

				t.Cleanup(func() {
					ts.Close()
				})
			}
			set := flag.NewFlagSet("test", 0)
			app := &cli.App{Writer: ioutil.Discard}
			assert.NoError(t, set.Parse(append([]string{name}, tc.cmdArgs...)))

			ctx := cli.NewContext(app, set, nil)
			err := command.Run(ctx)

			if tc.expectedErrMsg != "" {
				assert.Error(t, err)
				fmt.Println(err.Error())
				assert.True(t, strings.Contains(err.Error(), tc.expectedErrMsg))
			} else {
				assert.NoError(t, err)
			}

			if tc.afterRunExec != nil {
				assert.NoError(t, tc.afterRunExec())
			}
		})
	}
}

func TestNew(t *testing.T) {
	// get version
	buf := new(bytes.Buffer)
	app := cmd.New()
	app.Writer = buf
	err := app.Run([]string{"cir-rotator", "--version"})
	output := buf.String()

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("cir-rotator version %s\n", app.Version), output)
}
