package filter_test

import (
	"testing"
	"time"

	fl "github.com/iomarmochtar/cir-rotator/pkg/filter"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// the valid one will retruning IFilterEngine
	engine, err := fl.New([]string{"Now() >= CreatedAt", `Repository matches '.*secret-souce.*'`})
	assert.NoError(t, err)
	assert.NotNil(t, engine)

	// valid to compile rules
	engine, err = fl.New([]string{`Repository matches '*secret-souce$'`})
	assert.Error(t, err, "error in rules")
	assert.Nil(t, engine)
}

func TestNew_Process(t *testing.T) {
	testCases := map[string]struct {
		filters        []string
		fields         fl.Fields
		expectedResult bool
		expectedErrMsg string
	}{
		"will match by one of the conditions": {
			filters:        []string{"1 != 1", "5 < 4", `'strongman' matches '.*man$'`},
			fields:         fl.Fields{},
			expectedResult: true,
		},
		"passing variable for can be used in filter": {
			filters:        []string{`"latest" in Tags`},
			fields:         fl.Fields{Tags: []string{"latest", "release-abc"}},
			expectedResult: true,
		},
		"error while processing filters": {
			filters:        []string{`50 in ["halo"]`},
			expectedErrMsg: "reflect.Value.MapIndex: value of type int is not assignable to type string (1:5)\n | (50 in [\"halo\"])\n | ....^",
		},
		"custom helper Date()": {
			filters:        []string{"1 + 1 == 2 && CreatedAt >= Date('1991-06-13')"},
			fields:         fl.Fields{CreatedAt: time.Date(1991, time.June, 15, 0, 0, 0, 0, time.Local)},
			expectedResult: true,
		},
		"custom helper Now() and Duration()": {
			filters: []string{
				"Duration('1h') > Duration('3h') || Now() + Duration('1h') > CreatedAt || true",
			},
			fields:         fl.Fields{CreatedAt: time.Date(1991, time.June, 15, 0, 0, 0, 0, time.Local)},
			expectedResult: true,
		},
		"wrong in Duration()": {
			filters:        []string{"Duration('wrong')"},
			expectedErrMsg: "time: invalid duration \"wrong\" (1:2)\n | (Duration('wrong'))\n | .^",
			expectedResult: false,
		},
		"custom Duration d": {
			filters:        []string{"Duration('1d') == Duration('24h')"},
			expectedResult: true,
		},
		"custom Duration M": {
			filters:        []string{"Duration('3M') == Duration('2160h')"},
			expectedResult: true,
		},
		"custom Duration Y": {
			filters:        []string{"Duration('1Y') == Duration('8640h')"},
			expectedResult: true,
		},
		"custom Duration combination": {
			filters:        []string{"Duration('2d3M1Y1h') == Duration('10849h')"},
			expectedResult: true,
		},
		"custom Duration unknown unit": {
			filters:        []string{"Duration('1x') == Duration('1h')"},
			expectedErrMsg: "time: unknown unit \"x\" in duration \"1x\" (1:2)\n | (Duration('1x') == Duration('1h'))\n | .^",
			expectedResult: false,
		},
		"wrong date pattern in Date()": {
			filters:        []string{"Date('13-06-1991')"},
			expectedErrMsg: "parsing time \"13-06-1991\" as \"2006-01-02\": cannot parse \"6-1991\" as \"2006\" (1:2)\n | (Date('13-06-1991'))\n | .^",
			expectedResult: false,
		},
		"Operators": {
			filters: []string{
				"Duration('1h') < Duration('1s') || Duration('3h') >= Duration('5h')", // false
				"UploadedAt < CreatedAt && UploadedAt - CreatedAt <= Duration('3s')",  // false
				"Duration('1m') <= Duration('1s')",                                    // false
				"UploadedAt <= CreatedAt",                                             // false
				"UploadedAt - CreatedAt == Duration('48h')",                           // true
			},
			expectedResult: true,
			fields: fl.Fields{
				CreatedAt:  time.Date(1991, time.June, 13, 0, 0, 0, 0, time.Local),
				UploadedAt: time.Date(1991, time.June, 15, 0, 0, 0, 0, time.Local),
			},
		},
		"equal time operator": {
			filters: []string{
				`Date('2020-02-02') == Date('2020-02-02')`,
			},
			expectedResult: true,
		},
		"size equal": {
			filters: []string{
				"ImageSize == SizeStr('1 KiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024,
			},
			expectedResult: true,
		},
		"size gte": {
			filters: []string{
				"ImageSize >= SizeStr('1 MiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024 * 1024,
			},
			expectedResult: true,
		},
		"size gte not match": {
			filters: []string{
				"ImageSize >= SizeStr('1 MiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024*1024 - 1,
			},
			expectedResult: false,
		},
		"size gt": {
			filters: []string{
				"ImageSize > SizeStr('2 MiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024*1024*2 + 1,
			},
			expectedResult: true,
		},
		"size lte": {
			filters: []string{
				"ImageSize <= SizeStr('10 MiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024 * 1024 * 10,
			},
			expectedResult: true,
		},
		"size lte not match": {
			filters: []string{
				"ImageSize <= SizeStr('10 MiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024*1024*10 + 1,
			},
			expectedResult: false,
		},
		"size lt": {
			filters: []string{
				"ImageSize < SizeStr('10 MiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024*1024*10 - 1,
			},
			expectedResult: true,
		},
		"size lt not match": {
			filters: []string{
				"ImageSize < SizeStr('10 MiB')",
			},
			fields: fl.Fields{
				ImageSize: 1024 * 1024 * 10,
			},
			expectedResult: false,
		},
		"SizeStr wrong pattern": {
			filters: []string{
				"ImageSize < SizeStr('not valid')",
			},
			fields: fl.Fields{
				ImageSize: 1024 * 1024 * 10,
			},
			expectedErrMsg: "unknown pattern not valid (1:14)\n | (ImageSize < SizeStr('not valid'))\n | .............^",
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			filterEngine, err := fl.New(tc.filters)
			assert.NoError(t, err)

			result, err := filterEngine.Process(tc.fields)
			if tc.expectedErrMsg != "" {
				assert.EqualError(t, err, tc.expectedErrMsg)
				assert.False(t, result)
			} else {
				assert.Equal(t, result, tc.expectedResult)
				assert.NoError(t, err)
			}
		})
	}
}
