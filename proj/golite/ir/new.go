package ir

import (
	"bytes"
	"fmt"
)

type New struct {
	target   int
	dataType string
	size     int
}

func NewNew(target int, dataType string, size int) *New {
	return &New{target, dataType, size}
}

func (instr *New) GetTargets() []int {
	target := []int{}
	target = append(target, instr.target)
	return target
}

func (instr *New) GetSources() []int { return []int{} }

func (instr *New) GetImmediate() *int { return nil }

func (instr *New) GetSourceString() string {
	return instr.dataType
}

func (instr *New) GetLabel() string { return "" }

func (instr *New) SetLabel(newLabel string) {}

func (instr *New) String() string {
	var out bytes.Buffer
	targetReg := fmt.Sprintf("r%v",instr.target)
	out.WriteString(fmt.Sprintf("    new %s,%s",targetReg,instr.dataType))
	return out.String()
}

func (instr *New) TranslateToAssembly(funcVarDict map[int]int, paramRegIds map[int]int) []string {
	instruction := []string{}

	// prepare for malloc, push x0... to stack
	offset := 16
	for i := 0; i < len(paramRegIds); i++ {
		instruction = append(instruction, fmt.Sprintf("\tstr x%v,[x29,%v]",i,offset))
		offset += 8
	}

	space := instr.size * 8
	instruction = append(instruction, fmt.Sprintf("\tmov x0,#%v", space))
	instruction = append(instruction, "\tbl malloc")
	targetOffset := funcVarDict[instr.target]
	instruction = append(instruction, fmt.Sprintf("\tstr x0,[x29,#%v]",targetOffset))

	// restore registers after malloc
	offset = 16
	for i := 0; i < len(paramRegIds); i++ {
		instruction = append(instruction, fmt.Sprintf("\tldr x%v,[x29,%v]",i,offset))
		offset += 8
	}

	return instruction
}