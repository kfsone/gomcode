package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParamArray(t *testing.T) {
	actual := NewParamArray("t", "1", "xyz", 222, "abc", "A")
	expected := []Param{Param{'T', "1"}, Param{'X', "222"}, Param{'A', "A"}}
	assert.Equal(t, expected, actual)
}

func TestNewCode(t *testing.T) {
	newCode := NewCode("T123", "---test---")
	myCode := Code{GCode: "T123", Comment: "---" + "test" + "---", HideChecksum: false, LineNo: 0}
	assert.Equal(t, myCode, newCode)
	assert.Equal(t, 0, len(newCode.Parameters))

	myCode.Parameters = []Param{Param{'T', "1"}, Param{'U', "2"}}
	newCode = NewCode("T123", "---test---", Param{'T', "1"}, Param{'U', "2"})
	assert.Equal(t, myCode, newCode)
}

func TestGCodeChecksumEmptyString(t *testing.T) {
	assert.Equal(t, 0, GCodeChecksum(""))
}

func TestGCodeChecksumA(t *testing.T) {
	assert.Equal(t, int('A'), GCodeChecksum("A"))
}

func TestGCodeChecksumAB(t *testing.T) {
	assert.Equal(t, int('A')^int('B'), GCodeChecksum("AB"))
}

func TestGCodeEqual(t *testing.T) {
	m101 := NewCode("M101", "")
	m102 := NewCode("M102", "")
	code := NewCode("M101", "")

	assert.Equal(t, "M101", m101.GCode)
	assert.Equal(t, "M102", m102.GCode)
	assert.Equal(t, "M101", code.GCode)

	assert.True(t, m101.Equal(m101))
	assert.True(t, m102.Equal(m102))
	assert.True(t, m101.Equal(code))
	assert.False(t, m101.Equal(m102))
	assert.False(t, code.Equal(m102))

	code.Comment = "---"
	assert.Equal(t, "", m101.Comment)
	assert.Equal(t, "---", code.Comment)
	assert.True(t, m101.Equal(code))

	code.HideChecksum = true
	assert.True(t, m101.Equal(code))

	code.LineNo += 1
	assert.True(t, m101.Equal(code))

	code.Parameters = []Param{Param{'A', "1"}}
	assert.False(t, m101.Equal(code))
	assert.False(t, m102.Equal(code))

	m101.Parameters = []Param{Param{'A', "1"}}
	assert.True(t, m101.Equal(code))
	assert.False(t, m102.Equal(code))
}

func TestGCodeOverride(t *testing.T) {
	code := NewCode("T101", "")
	_, ok := code.Parameter('A')
	assert.False(t, ok)
	assert.Equal(t, 0, len(code.Parameters))

	assert.Nil(t, code.Override('T', "99"))
	assert.Equal(t, 1, len(code.Parameters))
	value, ok := code.Parameter('T')
	assert.True(t, ok)
	assert.Equal(t, "99", value)
}

func TestGCodeEmitBasic(t *testing.T) {
	code := NewCode("M101", "")
	assert.Equal(t, "M101", code.Emit(0))
}

func TestGCodeBasicLineNo(t *testing.T) {
	code := NewCode("M101", "")
	// by default, emiting a line number also adds a checksum
	code.HideChecksum = true
	assert.Equal(t, "N335 M101", code.Emit(335))
	// now try with the checksum
	code.HideChecksum = false
	assert.Equal(t, "N335 M101*38", code.Emit(335))
}

func TestGCodeEmitCommentAndParams(t *testing.T) {
	code := NewCode("M123", "the comment", Param{'A', "111"}, Param{'B', "234"}, Param{'z', "935"})
	assert.Equal(t, "M123 A111 B234 Z935 ;the comment", code.Emit(0))
}

func TestGCodeEmitCommentAndParamsLineNo(t *testing.T) {
	code := NewCode("M123", "the comment", Param{'A', "111"}, Param{'B', "234"}, Param{'z', "935"})
	// by default, emiting a line number also adds a checksum
	code.HideChecksum = true
	assert.Equal(t, "N4294967297 M123 A111 B234 Z935 ;the comment", code.Emit(4294967297))
	// now try with the checksum
	code.HideChecksum = false
	assert.Equal(t, "N1844674407379551617 M123 A111 B234 Z935*98 ;the comment", code.Emit(1844674407379551617))
}
