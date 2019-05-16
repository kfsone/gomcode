package main

import (
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

func getRune(param interface{}) rune {
	switch typedval := param.(type) {
	case rune:
		{
			return unicode.ToUpper(typedval)
		}
	case byte:
		{
			return unicode.ToUpper(rune(typedval))
		}
	case []byte:
		{
			return unicode.ToUpper(rune(typedval[0]))
		}
	case string:
		{
			return unicode.ToUpper(rune(typedval[0]))
		}
	default:
		{
			panic("Parameter key must be a character")
		}
	}
}

func NewParamArray(params ...interface{}) (parameters []Param) {
	parameters = make([]Param, 0, len(params)/2)
	if len(params)&1 != 0 {
		panic("NewParamArray requires even number of arguments")
	}
	for i := 0; i < len(params); i++ {
		key := unicode.ToUpper(getRune(params[i]))
		i++
		switch value := params[i].(type) {
		case bool:
			{
				if value {
					parameters = append(parameters, Param{key, ""})
				}
				continue
			}
		default:
			{
				parameters = append(parameters, Param{key, fmt.Sprint(value)})
			}
		}
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
		if err := code.Override(param.Key, param.Value); err != nil {
			panic(err)
		}
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
		return fmt.Errorf("Invalid key: %c", key)
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
	if lineNo > 0 && c.GCode != "M110" {
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
