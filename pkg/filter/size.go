package filter

import (
	"github.com/antonmedv/expr"
	h "github.com/iomarmochtar/cir-rotator/pkg/helpers"
)

var (
	cachedSizeStr map[string]float64
)

func init() {
	cachedSizeStr = make(map[string]float64)
}

type imageSize struct{}

func (imageSize) SizeStr(input string) float64 {
	if cstr := cachedSizeStr[input]; cstr != 0 {
		return cachedSizeStr[input]
	}

	result, err := h.SizeUnitStrToFloat(input)
	if err != nil {
		panic(err)
	}
	cachedSizeStr[input] = result
	return result
}

func (imageSize) EqualSize(a uint, b float64) bool        { return float64(a) == b }
func (imageSize) GreaterEqualSize(a uint, b float64) bool { return float64(a) >= b }
func (imageSize) GreaterSize(a uint, b float64) bool      { return float64(a) > b }
func (imageSize) LessEqualSize(a uint, b float64) bool    { return float64(a) <= b }
func (imageSize) LessSize(a uint, b float64) bool         { return float64(a) < b }

func sizeOperations() []expr.Option {
	return []expr.Option{
		expr.Operator("==", "EqualSize"),
		expr.Operator(">=", "GreaterEqualSize"),
		expr.Operator(">", "GreaterSize"),
		expr.Operator("<=", "LessEqualSize"),
		expr.Operator("<", "LessSize"),
	}
}
