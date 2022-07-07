package filter

// taken and modified from: https://github.com/antonmedv/expr/blob/2c1881a9909453f9f1047886d9be0b94e0ba5c48/docs/examples/dates_test.go

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/antonmedv/expr"
)

const (
	hourDayNum   = 24
	hourMonthNum = hourDayNum * 30
	hourYearNum  = hourMonthNum * 12
)

var (
	customDurationRe     = regexp.MustCompile(`(\d+)(d|M|Y)`)
	customDurationMapper = map[string]int{
		"d": hourDayNum,
		"M": hourMonthNum,
		"Y": hourYearNum,
	}
)

// customDurationFinder replace the known custom duration to hour unit
func customDurationFinder(ptrn string) (string, error) {
	matches := customDurationRe.FindAllStringSubmatch(ptrn, -1)
	for im := range matches {
		match := matches[im]
		numMatch, err := strconv.Atoi(match[1])
		if err != nil {
			return "", err
		}
		convertedNum := customDurationMapper[match[2]] * numMatch
		ptrn = strings.ReplaceAll(ptrn, match[0], fmt.Sprintf("%dh", convertedNum))
	}
	return ptrn, nil
}

type datetime struct{}

func (datetime) Date(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}
func (datetime) Duration(s string) time.Duration {
	s, err := customDurationFinder(s)
	if err != nil {
		panic(err)
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}
func (datetime) Now() time.Time                                     { return time.Now() }
func (datetime) Equal(a, b time.Time) bool                          { return a.Equal(b) }
func (datetime) Before(a, b time.Time) bool                         { return a.Before(b) }
func (datetime) BeforeOrEqual(a, b time.Time) bool                  { return a.Before(b) || a.Equal(b) }
func (datetime) After(a, b time.Time) bool                          { return a.After(b) }
func (datetime) AfterOrEqual(a, b time.Time) bool                   { return a.After(b) || a.Equal(b) }
func (datetime) AddTimeWDur(a time.Time, b time.Duration) time.Time { return a.Add(b) }
func (datetime) SubTimeWTime(a, b time.Time) time.Duration          { return a.Sub(b) }
func (datetime) EqualDuration(a, b time.Duration) bool              { return a == b }
func (datetime) BeforeDuration(a, b time.Duration) bool             { return a < b }
func (datetime) BeforeOrEqualDuration(a, b time.Duration) bool      { return a <= b }
func (datetime) AfterDuration(a, b time.Duration) bool              { return a > b }
func (datetime) AfterOrEqualDuration(a, b time.Duration) bool       { return a >= b }

func datetimeOperations() []expr.Option {
	return []expr.Option{
		// Operators override for date comprising.
		expr.Operator("==", "Equal"),
		expr.Operator("<", "Before"),
		expr.Operator("<=", "BeforeOrEqual"),
		expr.Operator(">", "After"),
		expr.Operator(">=", "AfterOrEqual"),

		// Time and duration manipulation.
		expr.Operator("+", "AddTimeWDur"),
		expr.Operator("-", "SubTimeWTime"),

		// Operators override for duration comprising.
		expr.Operator("==", "EqualDuration"),
		expr.Operator("<", "BeforeDuration"),
		expr.Operator("<=", "BeforeOrEqualDuration"),
		expr.Operator(">", "AfterDuration"),
		expr.Operator(">=", "AfterOrEqualDuration"),
	}
}
