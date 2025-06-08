package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ParseChapterFlags(chapter int, chaptersStr string, from, to int) ([]int, error) {
	var result []int

	switch {
	case chapter > 0:
		result = append(result, chapter)

	case chaptersStr != "":
		split := strings.Split(chaptersStr, ",")
		for _, s := range split {
			n, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil {
				return nil, fmt.Errorf("invalid chapter number: %s", s)
			}
			result = append(result, n)
		}

	case from > 0 && to > 0 && from <= to:
		for i := from; i <= to; i++ {
			result = append(result, i)
		}

	default:
		return nil, errors.New("please provide at least one valid chapter selection method")
	}

	return result, nil
}
