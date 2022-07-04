package filter

import (
	"fmt"
	"strings"
	"time"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/rs/zerolog/log"
)

type Fields struct {
	datetime
	imageSize
	Repository string
	Digest     string
	ImageSize  uint
	Tags       []string
	CreatedAt  time.Time
	UploadedAt time.Time
}

//go:generate mockgen -destination mock_filter/mock_filter.go -source filter.go IFilterEngine
type IFilterEngine interface {
	Process(fields Fields) (result bool, err error)
}

const filtersJoiner = " || "

type Engine struct {
	program *vm.Program
}

func New(filters []string) (IFilterEngine, error) {
	finalFilter := make([]string, len(filters))
	// add bracket for group each filter then use OR logic between them
	for idx, f := range filters {
		finalFilter[idx] = fmt.Sprintf("(%s)", f)
	}

	options := append([]expr.Option{expr.Env(&Fields{})}, datetimeOperations()...)
	options = append(options, sizeOperations()...)

	strFilter := strings.Join(finalFilter, filtersJoiner)
	log.Debug().Str("filter", strFilter).Msg("compiling filter")

	program, err := expr.Compile(strFilter, options...)
	if err != nil {
		return nil, err
	}

	return &Engine{program: program}, nil
}

func (f Engine) Process(fields Fields) (result bool, err error) {
	output, err := expr.Run(f.program, fields)
	if err != nil {
		return false, err
	}

	return output.(bool), nil
}
