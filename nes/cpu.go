package nes

import (
	"log"
)

type regs6502 struct {
	pc     uint16
	sp     uint8
	acc    uint8
	x      uint8
	y      uint8
	status uint8
}

func (r *regs6502) getCarry() bool {
	return (r.status & 0x1) == 1
}

func (r *regs6502) setCarry(val uint8) {
	r.status |= ((val & 1) << 0)
}

type Cpu struct {
	regs regs6502
}

func (cpu *Cpu) Init() error {
	log.Println("Initializing CPU...")
	return nil
}
