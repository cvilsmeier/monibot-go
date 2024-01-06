package histogram

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestStringifyValues(t *testing.T) {
	ass := assert.New(t)
	ass.Eq("", StringifyValues([]int64{}))
	ass.Eq("0", StringifyValues([]int64{0}))
	ass.Eq("0:2", StringifyValues([]int64{0, 0}))
	ass.Eq("1,2", StringifyValues([]int64{2, 1}))
	ass.Eq("0:3,1,2,3:3,4", StringifyValues([]int64{0, 0, 0, 1, 2, 3, 3, 3, 4}))
	ass.Eq("0:3,1,2,3:3,4", StringifyValues([]int64{0, 4, 3, 0, 2, 3, 0, 3, 1}))
}

func TestParseValues(t *testing.T) {
	str := func(values []int64, err error) string {
		if err != nil {
			return err.Error()
		}
		var ss []string
		for _, v := range values {
			ss = append(ss, strconv.FormatInt(v, 10))
		}
		return strings.Join(ss, ",")
	}
	ass := assert.New(t)
	v, err := ParseValues("")
	ass.Eq("", str(v, err))
	v, err = ParseValues("1")
	ass.Eq("1", str(v, err))
	v, err = ParseValues("1:1")
	ass.Eq("1", str(v, err))
	v, err = ParseValues("1:2")
	ass.Eq("1,1", str(v, err))
	v, err = ParseValues("3:2,2:1,1")
	ass.Eq("1,2,3,3", str(v, err))
	v, err = ParseValues("-3:2")
	ass.Eq("cannot parse token #1 \"-3:2\": invalid value -3", str(v, err))
	v, err = ParseValues("-3:2")
	ass.Eq("cannot parse token #1 \"-3:2\": invalid value -3", str(v, err))
	v, err = ParseValues("3:0")
	ass.Eq("cannot parse token #1 \"3:0\": invalid count 0", str(v, err))
	v, err = ParseValues("foo")
	ass.Eq("cannot parse token #1 \"foo\": cannot parse value \"foo\": strconv.ParseInt: parsing \"foo\": invalid syntax", str(v, err))
	v, err = ParseValues("foo:1")
	ass.Eq("cannot parse token #1 \"foo:1\": cannot parse value \"foo\": strconv.ParseInt: parsing \"foo\": invalid syntax", str(v, err))
	v, err = ParseValues("1:foo")
	ass.Eq("cannot parse token #1 \"1:foo\": cannot parse count \"foo\": strconv.Atoi: parsing \"foo\": invalid syntax", str(v, err))
}
