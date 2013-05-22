package nes

import (
	"log"
    "fmt"
)

const (
	InternalRAM   = 0x0000
	InputOutput   = 0x2000
	ExpansionMods = 0x5000
	CartridgeRAM  = 0x6000
	LowerBankCart = 0x8000
	UpperBankCart = 0xc000
)

type regs struct {
	pc     uint16
	sp     uint8
	acc    uint8
	x      uint8
	y      uint8
	status uint8
}

func (r *regs) getCarryFlag() bool {
	return (r.status & 0x1) != 0
}

func (r *regs) getNegativeFlag() bool {
    return (r.status & 0x80) != 0
}

func (r *regs) setCarryFlag() {
	r.status |= 0x1
}

func (r *regs) setZeroFlag() {
	r.status |= 0x2
}

func (r *regs) setInterruptDisable() {
    r.status |= 0x4
}

func (r *regs) setNegativeFlag() {
    r.status |= 0x80
}

func (r *regs) clearZeroFlag() {
	r.status &= (0xff - 0x2)
}

func (r *regs) clearDecimalMode() {
    r.status &= (0xff - 0x8)
}

func (r *regs) clearNegativeFlag() {
    r.status &= (0xff - 0x80)
}

type Cpu struct {
	r regs
    instrOpcodes []func()
}

func (c *Cpu) Init() (e error) {
	log.Println("Initializing CPU...")
    c.r.pc = 0x8000
    
	return
}

func (c *Cpu) String() (out string) {
    out += fmt.Sprintf("PC = %02x, ACC = %02x, X = %02x, Y = %02x\n",
        c.r.pc, c.r.acc, c.r.x, c.r.y)
    out += fmt.Sprintf("SP = %02x, Status = %02x", 
        c.r.sp, c.r.status)
    return
}

func (c *Cpu) testAndSetZero(val uint8) {
    if val == 0 {
        c.r.setZeroFlag()
    } else {
        c.r.clearZeroFlag()
    }
}

func (c *Cpu) testAndSetSign(val uint8) {
    if (val & 0x80) != 0 {
        c.r.setNegativeFlag()
    } else {
        c.r.clearNegativeFlag()
    }
}

func (c *Cpu) immediateAddress() (result uint16) {
    result = c.r.pc
    c.r.pc += 1
    return
}

func (c *Cpu) relativeAddress() (result uint16) {
    tmp, _ := Ram.Read(c.r.pc)
    result = uint16(int8(tmp))
    c.r.pc += 1
    return
}

func (c *Cpu) absoluteAddress() (result uint16) {
    low, _ := Ram.Read(c.r.pc)
    high, _ := Ram.Read(c.r.pc + 1)
    result = (uint16(high) << 8) | uint16(low)
    c.r.pc += 2
    return
}

func (c *Cpu) transfer(dstReg *uint8, srcReg uint8) {
    *dstReg = srcReg
}

func (c *Cpu) ldReg(addr uint16, reg *uint8) {
    *reg, _ = Ram.Read(addr)

    c.testAndSetZero(*reg)
    c.testAndSetSign(*reg)
    return
}

func (c *Cpu) stReg(addr uint16, reg uint8) {
    _ = Ram.Write(addr, c.r.acc)
}

func (c *Cpu) jump(addr uint16) {
    c.r.pc = addr
}

func (c *Cpu) Step() (err error) {
    log.Printf("PC: 0x%04x\n", c.r.pc)
    opc, _ := Ram.Read(c.r.pc)
    c.r.pc += 1

    switch opc {
    case 0x10:  // BPL
        if c.r.getNegativeFlag() {
            c.jump(c.r.pc + c.relativeAddress())
        }
    case 0x78:  // SEI
        c.r.setInterruptDisable()
    case 0x8d:  // STA [ABS]
        c.stReg(c.absoluteAddress(), c.r.acc)
    case 0x9a:  // TXS
        c.transfer(&c.r.sp, c.r.x)
    case 0xa2:  // LDX IMM
        c.ldReg(c.immediateAddress(), &c.r.x)
    case 0xa9:  // LDA IMM
        c.ldReg(c.immediateAddress(), &c.r.acc)
    case 0xad:  // LDA [ABS]
        c.ldReg(c.absoluteAddress(), &c.r.acc)
    case 0xd8:  // CLD
        c.r.clearDecimalMode()
    default:
        log.Printf("Unimplemented opcode 0x%02x\n", opc)
        log.Fatalf("CPU State:\n%s\n", c)
    }
    //log.Printf("CPU State:\n%s\n", c)
    return
}
