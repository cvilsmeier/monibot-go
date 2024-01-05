package histogram

import (
	"sort"
	"strconv"
	"strings"
)

func StrValues(values []int64) string {
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
