package cmd_test

import (
	"testing"

	"github.com/iomarmochtar/cir-rotator/app/cmd"
	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
)

func TestListAction(t *testing.T) {
	listTestCases := h.CombineMaps(commonTestCases, map[string]caseParam{
		"not providing any output": {
			cmdArgs:        []string{"-ho", "asia.gcr.io"},
			expectedErrMsg: "must specified one or more output",
		},
	})
	runCmdTestCases("list", cmd.ListAction(), listTestCases, t)
}
