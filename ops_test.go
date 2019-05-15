package gomcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolIdx0(t *testing.T) {
	assert.Equal(t, Code{GCode: "T0", Comment: "select tool"}, ToolIdx(0))
}

func TestToolIdx1(t *testing.T) {
	assert.Equal(t, Code{GCode: "T1", Comment: "select tool"}, ToolIdx(1))
}

func TestToolIdx999(t *testing.T) {
	assert.Equal(t, Code{GCode: "T999", Comment: "select tool"}, ToolIdx(999))
}

func TestLineNo(t *testing.T) {
	assert.Equal(t, Code{GCode: "M110", Comment: "set line no", Parameters: NewParamArray("N", "12995"), LineNo: 12995}, LineNo(12995))
}

func TestHotendTemp(t *testing.T) {
	assert.Equal(t, Code{GCode: "M104", Comment: "set hotend temp", Parameters: NewParamArray("S", "204")}, HotendTemp(204))
}

func TestHotendTempMaxAuto(t *testing.T) {
	expected := Code{GCode: "M104", Comment: "set hotend temp and max auto", Parameters: NewParamArray("S", "222", "B", "180", "F", "")}
	assert.Equal(t, expected, HotendTempMaxAuto(222, 180))
}
