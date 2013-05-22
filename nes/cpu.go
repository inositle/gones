package nes

import (
	"log"
)

const (
	InternalRAM   = 0x0000
	InputOutput   = 0x2000
	ExpansionMods = 0x5000
	CartridgeRAM  = 0x6000
	LowerBankCart = 0x8000
	UpperBankCart = 0xc000
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
    instrOpcodes []func()
}

func (cpu *Cpu) Init() error {
	log.Println("Initializing CPU...")
	return nil
}

func (cpu *Cpu) testAndSetZero(val uint8) {
}

func (cpu *Cpu) testAndSetNegative(val uint8) {
}

func (cpu *Cpu) Lda(location uint16) {
    cpu.regs.acc, _ = Ram.Read(location)

    cpu.testAndSetZero(cpu.regs.acc)
    cpu.testAndSetNegative(cpu.regs.acc)
}

func (cpu *Cpu) Step() (err error) {
    return
}
