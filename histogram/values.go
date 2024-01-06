package histogram

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// StringifyValues formats histogram values into a string.
// The result is a comma-separated list of "value:count" pairs.
// Each value is a non-negative int64 value, each count is a
// greater or equal to 1 int.
// If count is 1, the ":count" part is left out, so "42:1" and
// "42" are sematically equal.
func StringifyValues(values []int64) string {
	n := len(values)
	if n == 0 {
		return ""
	}
	if n == 1 {
		return strconv.FormatInt(values[0], 10)
	}
	tmp := append([]int64(nil), values...)
	// cannot use slices.Sort() because we support go 1.20
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i] < tmp[j]
	})
	add := func(toks []string, v int64, c int) []string {
		if c == 1 {
			toks = append(toks, strconv.FormatInt(v, 10))
		} else {
			toks = append(toks, strconv.FormatInt(v, 10)+":"+strconv.Itoa(c))
		}
		return toks
	}
	var toks []string
	value := tmp[0]
	count := 1
	for i := 1; i < n; i++ {
		next := tmp[i]
		if next == value {
			count++
		} else {
			toks = add(toks, value, count)
			value = next
			count = 1
		}
	}
	if count >= 1 {
		toks = add(toks, value, count)
	}
	return strings.Join(toks, ",")
}

// ParseValues parses a "value:count" comma-separated
// list of histogram values into a slice of int64 values.
// The resulting slice contains each value ":count" times.
// The resulting slice is sorted in ascending order.
func ParseValues(s string) ([]int64, error) {
	if s == "" {
		return nil, nil
	}
	var values []int64
	toks := strings.Split(s, ",")
	for itok, tok := range toks {
		v, c, err := parseValueAndCount(tok)
		if err != nil {
			return nil, fmt.Errorf("cannot parse token #%d %q: %w", itok+1, tok, err)
		}
		for i := 0; i < c; i++ {
			values = append(values, v)
		}
	}
	// cannot use slices.Sort() because we support go 1.20
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values, nil
}

func parseValueAndCount(s string) (int64, int, error) {
	parts := strings.Split(s, ":")
	np := len(parts)
	if np < 1 || 2 < np {
		return 0, 0, fmt.Errorf("expected 1 or 2 parts but was %d", np)
	}
	value, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot parse value %q: %w", parts[0], err)
	}
	if value < 0 {
		return 0, 0, fmt.Errorf("invalid value %d", value)
	}
	count := 1
	if np == 2 {
		count, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, fmt.Errorf("cannot parse count %q: %w", parts[1], err)
		}
		if count < 1 {
			return 0, 0, fmt.Errorf("invalid count %d", count)
		}
	}
	return value, count, nil
}
