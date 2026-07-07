package azure

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var symbolArray []string = []string{"Y", "M", "D", "T", "H", "M", "S"}

var timemap map[int]time.Duration = map[int]time.Duration{
	0: 0,
	1: 0,
	2: time.Hour * 24,
	4: time.Hour,
	5: time.Minute,
	6: time.Second,
}

func parseISO8601Duration(s string) (time.Duration, error) {
	position, symbol := 0, 0
	var duration time.Duration = 0
	quantifier := ""
	foundT := false

	if len(s) < 3 {
		return 0, errInvalid(s)
	}
	if s[position] != 'P' {
		return 0, errInvalid(s)
	}
	for position = 1; position < len(s); position++ {
		currentChar := s[position]
		if currentChar >= '0' && currentChar <= '9' {
			quantifier += string(currentChar)
			continue
		}
		for {
			if symbol >= len(symbolArray) {
				return 0, errInvalid(s)
			}
			if symbol > 3 && !foundT {
				return 0, errInvalid(s)
			}
			if string(currentChar) == symbolArray[symbol] {
				if symbol == 0 || symbol == 1 { // Year or Month
					if quantifier != "" {
						return 0, fmt.Errorf("year and month are not supported, invalid string: %q", s)
					}
				}
				if symbol == 3 { // 'T'
					foundT = true
					break
				}
				quant, err := strconv.Atoi(quantifier)
				if err != nil {
					return 0, errInvalid(s)
				}
				duration += time.Duration(quant) * timemap[symbol]
				quantifier = ""
				symbol += 1
				break
			}
			symbol += 1
		}
	}
	return duration, nil
}

func errInvalid(s string) error {
	return fmt.Errorf("invalid ISO8601 duration %q", s)
}

// formatISO8601Duration returns d in ISO 8601 format, rounded down to the nearest
// minute. Zero and negative durations return an error.
func formatISO8601Duration(d time.Duration) (string, error) {
	if d < time.Minute {
		return "", errors.New("given duration must be 1 minute or more")
	}
	builder := strings.Builder{}
	// We only care about time, start with PT always
	builder.WriteString("PT")

	h := d / time.Hour
	d %= time.Hour
	m := d / time.Minute

	if h > 0 {
		builder.WriteString(fmt.Sprintf("%dH", h))
	}
	if m > 0 {
		builder.WriteString(fmt.Sprintf("%dM", m))
	}
	return builder.String(), nil
}
