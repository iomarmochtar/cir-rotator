package cmd_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/iomarmochtar/cir-rotator/app/cmd"
	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
)

func TestDeleteAction(t *testing.T) {
	deleteTestCases := h.CombineMaps(commonTestCases, map[string]caseParam{
		"not providing any params": {
			expectedErrMsg: `Required flag "ho" not set`,
		},
		"successfully deleting repositories": {
			cmdArgs: []string{"-u", "secret", "-p", "souce"},
			mockImageReg: func(w http.ResponseWriter, r *http.Request) error {
				data := readFixture("gcr/tag_list_no_child.json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(data)
				return err
			},
		},
		"using dry run option": {
			cmdArgs: []string{"-u", "secret", "-p", "souce", "--dry-run"},
			mockImageReg: func(w http.ResponseWriter, r *http.Request) error {
				// if it's there request for deleting then raise an error
				if r.Method == http.MethodDelete {
					return fmt.Errorf("will not deleting")
				}
				data := readFixture("gcr/tag_list_no_child.json")
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(data)
				return err
			},
		},
		"error while deleting image": {
			cmdArgs: []string{"-u", "secret", "-p", "souce"},
			mockImageReg: func(w http.ResponseWriter, r *http.Request) error {
				fixtureFile := "gcr/tag_list_no_child.json"
				// if it's there request for deleting then raise an error
				if r.Method == http.MethodDelete {
					fixtureFile = "gcr/error_delete_manifest.json"
				}
				data := readFixture(fixtureFile)
				w.WriteHeader(http.StatusOK)
				_, err := w.Write(data)
				return err
			},
			expectedErrMsg: "Failed to compute blob liveness for manifest: 'latest'",
		},
		"error if set worker count less than 1": {
			cmdArgs:        []string{"-ho", "https://asia.gcr.io/somepath", "-u", "secret", "-p", "souce", "--worker-count", "0"},
			expectedErrMsg: "invalid value for worker count: 0, make sure it's more than equal to 1",
		},
	})

	runCmdTestCases("delete", cmd.DeleteAction(), deleteTestCases, t)
}
