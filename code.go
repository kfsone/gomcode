package gomcode

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Param struct {
	Key   rune
	Value string
}

func NewParamArray(params ...string) (parameters []Param) {
	parameters = make([]Param, 0, len(params)/2)
	for i := 0; i < len(params); i += 2 {
		key := unicode.ToUpper(rune(params[i][0]))
		value := params[i+1]
		parameters = append(parameters, Param{key, value})
	}
	return parameters
}

type Code struct {
	GCode        string
	Comment      string
	Parameters   []Param
	HideChecksum bool
	LineNo       uint
}

func NewCode(gcode string, comment string, params ...Param) Code {
	code := Code{GCode: gcode, Comment: comment}
	for _, param := range params {
		if !unicode.IsLetter(param.Key) {
			panic("Invalid parameter: " + string(param.Key))
		}
		code.Override(param.Key, param.Value)
	}
	return code
}

func GCodeChecksum(code string) int {
	sum := 0
	for _, c := range code {
		sum ^= int(c)
	}
	return sum & 255
}

func (lhs Code) Equal(rhs Code) bool {
	return lhs.GCode == rhs.GCode && reflect.DeepEqual(lhs.Parameters, rhs.Parameters)
}

func (c Code) Parameter(key rune) (value string, ok bool) {
	for _, param := range c.Parameters {
		if param.Key == key {
			return param.Value, true
		}
	}
	return "", false
}

func (c *Code) Override(key rune, value string) error {
	key = unicode.ToUpper(key)
	if !unicode.IsLetter(key) {
		return errors.New(fmt.Sprintf("Invalid key: %c", key))
	}
	for _, param := range c.Parameters {
		if param.Key == key {
			param.Value = value
			return nil
		}
	}
	c.Parameters = append(c.Parameters, Param{key, value})
	return nil
}

func (c *Code) Emit(lineNo uint) string {
	// Max atoms will be:
	//  Nxxx    line number
	//  Mxxx    code
	//  kxxx    parameter
	atoms := make([]string, 0, 1+1+len(c.Parameters))
	if lineNo > 0 {
		atoms = append(atoms, fmt.Sprintf("N%d", lineNo))
	}
	atoms = append(atoms, c.GCode)
	for _, param := range c.Parameters {
		atoms = append(atoms, fmt.Sprintf("%c%s", param.Key, param.Value))
	}
	code := strings.Join(atoms, " ")

	if lineNo > 0 && !c.HideChecksum {
		code += "*" + strconv.Itoa(GCodeChecksum(code))
	}

	if c.Comment != "" {
		code += " ;" + c.Comment
	}

	return code
}
