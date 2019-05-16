package main

import "strconv"

type ToolId uint

func UintStr(value uint) string {
	return strconv.Itoa(int(value))
}

func ToolIdx(toolidx ToolId) Code {
	return NewCode("T"+UintStr(uint(toolidx)), "select tool")
}

func LineNo(lineNo uint) Code {
	if lineNo < 1 {
		panic("Cannot set line number less than 0")
	}
	code := NewCode("M110", "set line no", Param{'N', UintStr(lineNo)})
	code.LineNo = lineNo
	return code
}

func HotendTemp(celcius uint) Code {
	return NewCode("M104", "set hotend temp", Param{'S', UintStr(celcius)})
}

func HotendTempMaxAuto(celcius uint, maxAuto uint) Code {
	if maxAuto == 0 {
		return HotendTemp(celcius)
	}
	return NewCode("M104", "set hotend temp and max auto", Param{'S', UintStr(celcius)}, Param{'B', UintStr(maxAuto)}, Param{'F', ""})
}
