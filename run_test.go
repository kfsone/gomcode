package gomcode

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tearUp(t *testing.T) (Run, *strings.Builder) {
	// The Run class "execute"s by writing to an io.Writer for you,
	// so we're going to create a strings.Builder to capture that
	// output so we can test what gets written.
	writer := &strings.Builder{}
	assert.Equal(t, 0, writer.Len())

	r := NewRun(false, false, writer)
	assert.False(t, r.Checksum)
	assert.False(t, r.Comments)
	assert.Equal(t, 0, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, 0, writer.Len())

	return r, writer
}

func TestNewRun(t *testing.T) {
	_, _ = tearUp(t)
}

func TestRunQueue(t *testing.T) {
	r, writer := tearUp(t)

	// To be thorough, we'll check the length and values
	// of the queue as we build it up to ensure we get the
	// expected size AND order
	expecting := make([]string, 0, 10)
	checkGCodes := func() {
		assert.Equal(t, len(expecting), len(*r.cmdQueue))
		for idx, value := range expecting {
			assert.Equal(t, value, (*r.cmdQueue)[idx].GCode)
		}
	}

	// Queue one
	r.Queue(Code{GCode: "C1"})
	expecting = append(expecting, "C1")
	checkGCodes()

	// Queue another with the same code, repetition is valid.
	r.Queue(Code{GCode: "C1"})
	expecting = append(expecting, "C1")
	checkGCodes()

	// Queue a couple more individually
	r.Queue(Code{GCode: "C2"})
	expecting = append(expecting, "C2")
	r.Queue(Code{GCode: "C3"})
	expecting = append(expecting, "C3")
	checkGCodes()

	// Use the variadic version to queue 3 more
	r.Queue(Code{GCode: "C4"}, Code{GCode: "C3"}, Code{GCode: "C5"})
	expecting = append(expecting, "C4")
	expecting = append(expecting, "C3")
	expecting = append(expecting, "C5")
	checkGCodes()

	assert.Equal(t, 0, writer.Len())
}

func TestExecuteImmediateSingleBasic(t *testing.T) {
	r, writer := tearUp(t)

	err := r.ExecuteImmediate(Code{GCode: "E101", Comment: "-commentX-"})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(0), (*r.cmdHistory)[0].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "E101\n", writer.String())
	writer.Reset()

	r.Checksum = true

	err = r.ExecuteImmediate(Code{GCode: "E102", Comment: "-commentX-"})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(1), (*r.cmdHistory)[1].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "N1 E102*41\n", writer.String())
}

func TestExecuteImmediateSingleParams(t *testing.T) {
	r, writer := tearUp(t)

	err := r.ExecuteImmediate(NewCode("E103", "-commentX-", NewParamArray("b", 999, "c", 200)...))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(0), (*r.cmdHistory)[0].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "E103 B999 C200\n", writer.String())
	writer.Reset()

	r.Checksum = true
	err = r.ExecuteImmediate(NewCode("E104", "-commentX-", NewParamArray("b", 999, "c", 200)...))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(1), (*r.cmdHistory)[1].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "N1 E104 B999 C200*37\n", writer.String())
}

func TestExecuteImmediateSingleCmt(t *testing.T) {
	r, writer := tearUp(t)
	r.Comments = true

	err := r.ExecuteImmediate(NewCode("E105", "-comment1-"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(0), (*r.cmdHistory)[0].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "E105 ;-comment1-\n", writer.String())
	writer.Reset()

	r.Checksum = true

	err = r.ExecuteImmediate(NewCode("E106", "-comment2-"))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(1), (*r.cmdHistory)[1].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "N1 E106*45 ;-comment2-\n", writer.String())
}

func TestExecuteImmediateSingleCmtParams(t *testing.T) {
	r, writer := tearUp(t)
	r.Comments = true

	err := r.ExecuteImmediate(NewCode("E107", "-comment3-", NewParamArray("f", 11144, "G", "GVal")...))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(0), (*r.cmdHistory)[0].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "E107 F11144 GGVal ;-comment3-\n", writer.String())
	writer.Reset()

	r.Checksum = true

	err = r.ExecuteImmediate(NewCode("E108", "-comment4-", NewParamArray("H", 11145, "i", "iVal")...))
	assert.Nil(t, err)
	assert.Equal(t, 2, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(1), (*r.cmdHistory)[1].LineNo)

	// Make sure we actually emitted the op code
	assert.Equal(t, "N1 E108 H11145 IiVal*32 ;-comment4-\n", writer.String())
}

func getMultipleCodes(t *testing.T, comments bool, checksum bool) string {
	r, writer := tearUp(t)
	r.Comments = comments
	lineMultiplier := 0
	if checksum {
		lineMultiplier = 1
		r.Checksum = checksum
	}
	codes := []Code{NewCode("E201", ""), NewCode("E202", "-e202-"), NewCode("E204", "", NewParamArray('b', "133", 'c', "144", "d", "145")...), NewCode("E291", "-e291-", Param{'x', "987"})}

	err := r.ExecuteImmediate(codes...)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(*r.cmdHistory))
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, uint(1*lineMultiplier), (*r.cmdHistory)[0].LineNo)
	assert.Equal(t, uint(2*lineMultiplier), (*r.cmdHistory)[1].LineNo)
	assert.Equal(t, uint(3*lineMultiplier), (*r.cmdHistory)[2].LineNo)
	assert.Equal(t, uint(4*lineMultiplier), (*r.cmdHistory)[3].LineNo)

	return writer.String()
}

func TestExecuteImmediateMultiple(t *testing.T) {
	emitted := getMultipleCodes(t, false, false)
	assert.Equal(t, "E201\nE202\nE204 B133 C144 D145\nE291 X987\n", emitted)

	emitted = getMultipleCodes(t, true, false)
	assert.Equal(t, "E201\nE202 ;-e202-\nE204 B133 C144 D145\nE291 X987 ;-e291-\n", emitted)

	emitted = getMultipleCodes(t, false, true)
	assert.Equal(t, "N1 E201*41\nN2 E202*41\nN3 E204 B133 C144 D145*123\nN4 E291 X987*107\n", emitted)

	emitted = getMultipleCodes(t, true, true)
	assert.Equal(t, "N1 E201*41\nN2 E202*41 ;-e202-\nN3 E204 B133 C144 D145*123\nN4 E291 X987*107 ;-e291-\n", emitted)

}

func TestExecute(t *testing.T) {
	r, writer := tearUp(t)
	r.Comments = true
	r.Checksum = true

	r.Queue(NewCode("E300", ""), NewCode("E355", "-e355"))
	r.Queue(NewCode("E360", "", NewParamArray('S', 190, 'B', 180)...))
	r.Queue(NewCode("E321", "-e321", NewParamArray('x', 101, 'z', 15, 'y', 202, 'o', true)...))
	r.Queue(NewCode("G499", "", NewParamArray('o', false)...)) // expect: 'O' flag elided because its false
	assert.Equal(t, 5, len(*r.cmdQueue))                       // 2 on the first line
	assert.Equal(t, 0, len(*r.cmdHistory))
	assert.Equal(t, "E300", (*r.cmdQueue)[0].GCode)
	assert.Equal(t, "E355", (*r.cmdQueue)[1].GCode)
	assert.Equal(t, "E360", (*r.cmdQueue)[2].GCode)
	assert.Equal(t, "E321", (*r.cmdQueue)[3].GCode)
	assert.Equal(t, "G499", (*r.cmdQueue)[4].GCode)

	assert.Equal(t, 0, writer.Len())

	r.ExecuteImmediate(NewCode("M110", ""))
	assert.Equal(t, 1, len(*r.cmdHistory))
	assert.Equal(t, 5, len(*r.cmdQueue))
	assert.Equal(t, "M110", (*r.cmdHistory)[0].GCode)
	assert.Equal(t, "E300", (*r.cmdQueue)[0].GCode)
	assert.Equal(t, "G499", (*r.cmdQueue)[4].GCode)

	assert.Equal(t, "M110\n", writer.String())
	writer.Reset()

	r.Execute()
	assert.Equal(t, 0, len(*r.cmdQueue))
	assert.Equal(t, 6, len(*r.cmdHistory))
	assert.Equal(t, "M110", (*r.cmdHistory)[0].GCode)
	assert.Equal(t, "E300", (*r.cmdHistory)[1].GCode)
	assert.Equal(t, "E355", (*r.cmdHistory)[2].GCode)
	assert.Equal(t, "E360", (*r.cmdHistory)[3].GCode)
	assert.Equal(t, "E321", (*r.cmdHistory)[4].GCode)
	assert.Equal(t, "G499", (*r.cmdHistory)[5].GCode)

	assert.Equal(t, uint(0), (*r.cmdHistory)[0].LineNo)
	assert.Equal(t, uint(1), (*r.cmdHistory)[1].LineNo)
	assert.Equal(t, uint(2), (*r.cmdHistory)[2].LineNo)
	assert.Equal(t, uint(3), (*r.cmdHistory)[3].LineNo)
	assert.Equal(t, uint(4), (*r.cmdHistory)[4].LineNo)
	assert.Equal(t, uint(5), (*r.cmdHistory)[5].LineNo)

	lines := []string{
		"N1 E300*41",
		"N2 E355*42 ;-e355",
		"N3 E360 S190 B180*61",
		"N4 E321 X101 Z15 Y202 O*63 ;-e321",
		"N5 G499*40",
		"",
	}

	for idx, line := range strings.Split(writer.String(), "\n") {
		assert.Equal(t, lines[idx], line)
	}
}
