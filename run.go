package gomcode

import "io"

type Run struct {
	Checksum bool
	Comments bool
	writer   io.Writer

	LineNo     uint
	cmdHistory *[]Code
	cmdQueue   *[]Code
}

func NewRun(checksum bool, comments bool, writer io.Writer) Run {
	history, queue := make([]Code, 0, 1024), make([]Code, 0, 1024)
	return Run{Checksum: checksum, Comments: comments, writer: writer, cmdHistory: &history, cmdQueue: &queue}
}

func (r *Run) Reset() {
	*r.cmdQueue = (*r.cmdQueue)[:0]
	*r.cmdHistory = (*r.cmdHistory)[:0]
	r.LineNo = 0
}

func (r *Run) Queue(codes ...Code) {
	for _, code := range codes {
		*r.cmdQueue = append(*r.cmdQueue, code)
	}
}

func (r *Run) executeCode(cmds ...Code) error {
	lineNo := uint(0)
	if r.Checksum {
		lineNo = r.LineNo + 1
	}
	for _, code := range cmds {
		if r.Comments == false && code.Comment != "" {
			code.Comment = ""
		}
		if !r.Checksum || code.GCode == "M110" {
			code.HideChecksum = true
		}
		if _, err := r.writer.Write([]byte(code.Emit(lineNo) + "\n")); err != nil {
			return err
		}
		if lineNo > 0 {
			if code.GCode != "M110" {
				code.LineNo = lineNo
			} else {
				code.LineNo = 0
			}
			r.LineNo = code.LineNo
			lineNo += 1
		}
		*r.cmdHistory = append(*r.cmdHistory, code)
	}
	return nil
}

func (r *Run) ExecuteImmediate(cmds ...Code) error {
	return r.executeCode(cmds...)
}

func (r *Run) Execute() error {
	err := r.executeCode(*r.cmdQueue...)
	*r.cmdQueue = (*r.cmdQueue)[:0]
	return err
}
