package helpers

import (
	"fmt"
	"regexp"
	"strconv"
)

const (
	Kibibyte uint = 1024
	Mebibyte      = Kibibyte * 1024
	Gibibyte      = Mebibyte * 1024
	Tebibyte      = Gibibyte * 1024
	Pebibyte      = Tebibyte * 1024
	Exbibyte      = Pebibyte * 1024
)

var (
	imageSizeMatcher = regexp.MustCompile(`^(\d+(\.?\d+)?)\s?((Ki|Mi|Gi|Ti|Pi|Ei)?B)$`)
	unitShortsMapper = map[string]uint{
		"KiB": Kibibyte,
		"MiB": Mebibyte,
		"GiB": Gibibyte,
		"TiB": Tebibyte,
		"PiB": Pebibyte,
		"EiB": Exbibyte,
	}
)

// ByteCountIEC convert unsigned integer to human readable unit. take from https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
func ByteCountIEC(input uint) string {
	if input < Kibibyte {
		return fmt.Sprintf("%d B", input)
	}
	div, exp := Kibibyte, 0
	for n := input / Kibibyte; n >= Kibibyte; n /= Kibibyte {
		div *= Kibibyte
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(input)/float64(div), "KMGTPE"[exp])
}

func SizeUnitStrToFloat(input string) (float64, error) {
	matched := imageSizeMatcher.FindStringSubmatch(input)
	if len(matched) != 5 {
		return 0, fmt.Errorf("unknown pattern %s", input)
	}
	num, err := strconv.ParseFloat(matched[1], 64)
	if err != nil {
		return 0, err
	}
	var multiplier uint = 1
	if uMultiplier := unitShortsMapper[matched[3]]; uMultiplier != 0 {
		multiplier = uMultiplier
	}

	return num * float64(multiplier), nil
}
