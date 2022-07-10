package helpers_test

import (
	"os"
	"testing"
	"time"

	"github.com/iomarmochtar/cir-rotator/pkg/helpers"
	"github.com/stretchr/testify/assert"
)

func TestConvertTimeStrToUnix(t *testing.T) {
	parsed, err := helpers.ConvertTimeStrToUnix("1585800237411")
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2020, time.April, 2, 11, 3, 57, 0, time.Local), parsed)

	parsed, err = helpers.ConvertTimeStrToUnix("unknown")
	assert.Error(t, err)
	assert.Equal(t, time.Time{}, parsed)
}

func TestConvertTimeStrToReadAble(t *testing.T) {
	os.Setenv("TZ", "Asia/Jakarta")
	defer os.Unsetenv("TZ")

	parsed, err := helpers.ConvertTimeStrToReadAble("1585800237411")
	assert.NoError(t, err)
	assert.Equal(t, "2020-04-02T11:03:57+07:00", parsed, "parsed successfully")

	parsed, err = helpers.ConvertTimeStrToReadAble("2022-06-13")
	assert.Error(t, err, "wrong pattern")
	assert.Equal(t, "", parsed)
}
