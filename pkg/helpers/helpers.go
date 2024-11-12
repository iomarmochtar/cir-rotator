package helpers

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func IsInList[T comparable](needle T, haystack []T) bool {
	for _, val := range haystack {
		if val == needle {
			return true
		}
	}
	return false
}

func ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	//nolint:errcheck
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func SlashJoin(parts ...string) string {
	return strings.Join(parts, "/")
}

func CombineMaps[k string, v any](s1 map[k]v, s2 map[k]v) map[k]v {
	result := make(map[k]v)
	merge := func(kvS map[k]v, kvD map[k]v) {
		for a, b := range kvS {
			kvD[a] = b
		}
	}
	merge(s1, result)
	merge(s2, result)
	return result
}

func FileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// HumanizeDuration convert duration to readable string.
// taken from: https://gist.github.com/harshavardhana/327e0577c4fed9211f65?permalink_comment_id=2366908#gistcomment-2366908
func HumanizeDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	chunks := []struct {
		singularName string
		amount       int64
	}{
		{"day", days},
		{"hour", hours},
		{"minute", minutes},
		{"second", seconds},
	}

	parts := []string{}

	for _, chunk := range chunks {
		switch chunk.amount {
		case 0:
			continue
		case 1:
			parts = append(parts, fmt.Sprintf("%d %s", chunk.amount, chunk.singularName))
		default:
			parts = append(parts, fmt.Sprintf("%d %ss", chunk.amount, chunk.singularName))
		}
	}

	return strings.Join(parts, " ")
}
