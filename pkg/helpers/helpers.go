package helpers

import (
	"bufio"
	"os"
	"strings"
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
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
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
