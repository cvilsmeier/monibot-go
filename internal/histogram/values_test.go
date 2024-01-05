package histogram

import (
	"testing"

	"github.com/cvilsmeier/monibot-go/internal/assert"
)

func TestStrValues(t *testing.T) {
	ass := assert.New(t)
	ass.Eq("", StrValues([]int64{}))
	ass.Eq("0", StrValues([]int64{0}))
	ass.Eq("0:2", StrValues([]int64{0, 0}))
	ass.Eq("1,2", StrValues([]int64{2, 1}))
	ass.Eq("0:3,1,2,3:3,4", StrValues([]int64{0, 0, 0, 1, 2, 3, 3, 3, 4}))
	ass.Eq("0:3,1,2,3:3,4", StrValues([]int64{0, 4, 3, 0, 2, 3, 0, 3, 1}))
}
