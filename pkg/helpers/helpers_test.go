package helpers_test

import (
	"os"
	"testing"
	"time"

	"github.com/iomarmochtar/cir-rotator/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestIsInList(t *testing.T) {
	assert.True(t, helpers.IsInList(10, []int{3, 4, 10}))
	assert.False(t, helpers.IsInList("adi", []string{"sorry", "not", "found"}))
}

func TestReadLines(t *testing.T) {
	expectedLines := []string{`{"errors":[{"code":"NAME_UNKNOWN","message":"Failed to compute blob liveness for manifest: 'latest'"}]}`}
	result, err := helpers.ReadLines("../../testdata/gcr/error_delete_manifest.json")
	assert.NoError(t, err)
	assert.Equal(t, expectedLines, result, "read the contents of file")

	result, err = helpers.ReadLines("/this/path/is/not/found.txt")
	assert.True(t, os.IsNotExist(err), "expected error for file not found")
	assert.Empty(t, result)
}

func TestSlachJoin(t *testing.T) {
	assert.Equal(t, "hello/world/gogo", helpers.SlashJoin("hello", "world", "gogo"))
}

func TestCombineMaps(t *testing.T) {
	map1 := map[string]string{"name1": "val1", "name2": "val2"}
	map2 := map[string]string{"name2": "val3", "name4": "val4"}
	expected := map[string]string{"name1": "val1", "name2": "val3", "name4": "val4"}
	assert.Equal(t, expected, helpers.CombineMaps(map1, map2))
}

func TestFileExist(t *testing.T) {
	assert.True(t, helpers.FileExist("../../testdata/gcr/error_delete_manifest.json"))
	assert.False(t, helpers.FileExist("/path/not/exists.txt"))
}

func TestHumanizeDuration(t *testing.T) {
	dur, _ := time.ParseDuration("2h5m")
	assert.Equal(t, "2 hours 5 minutes", helpers.HumanizeDuration(dur))
	dur, _ = time.ParseDuration("4320m1s")
	assert.Equal(t, "3 days 1 second", helpers.HumanizeDuration(dur))
}
