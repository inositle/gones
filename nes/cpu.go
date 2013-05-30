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

func (r *regs) getZeroFlag() bool {
	return (r.status & 0x2) != 0
}

func (r *regs) getOverflowFlag() bool {
    return (r.status & 0x40) != 0
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

func (r *regs) setOverflowFlag() {
    r.status |= 0x40
}

func (r *regs) setNegativeFlag() {
    r.status |= 0x80
}

func (r *regs) clearCarryFlag() {
	r.status &= (0xff - 0x1)
}

func (r *regs) clearZeroFlag() {
	r.status &= (0xff - 0x2)
}

func (r *regs) clearDecimalMode() {
    r.status &= (0xff - 0x8)
}

func (r *regs) clearOverflowFlag() {
    r.status &= (0xff - 0x40)
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
    c.r.sp = 0xff
    
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

func (c *Cpu) zeroPageAddress() (result uint16) {
    tmp, _ := Ram.Read(c.r.pc)
    result = uint16(tmp)
    c.r.pc += 1
    return
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
    _ = Ram.Write(addr, reg)
}

func (c *Cpu) jump(addr uint16) {
    c.r.pc = addr
}

func (c *Cpu) push(val uint8) {
    _ = Ram.Write(0x100 + uint16(c.r.sp), val)
    c.r.sp--
}

func (c *Cpu) pop() (val uint8) {
    c.r.sp++
    val, _ = Ram.Read(0x100 + uint16(c.r.sp))
    return
}

func (c *Cpu) bitTest(addr uint16) {
    mem, _ := Ram.Read(addr)
    res := c.r.acc & mem

    if res == 0 {
        c.r.setZeroFlag()
    } else {
        c.r.clearZeroFlag()
    }

    if (mem & 0x80) != 0 {
        c.r.setNegativeFlag()
    } else {
        c.r.clearNegativeFlag()
    }

    if (mem & 0x40) != 0 {
        c.r.setOverflowFlag()
    } else {
        c.r.clearOverflowFlag()
    }
}

func (c *Cpu) Step() (err error) {
    log.Printf("PC: 0x%04x\n", c.r.pc)
    opc, _ := Ram.Read(c.r.pc)
    c.r.pc += 1

    switch opc {
    case 0x10:  // BPL
        rel := c.relativeAddress()
        if !c.r.getNegativeFlag() {
            c.jump(c.r.pc + rel)
        }
    case 0x18:  // CLC
        c.r.clearCarryFlag()
    case 0x20:  // JSR
        dst := c.absoluteAddress()
        tmp := c.r.pc - 1
        c.push(uint8(tmp >> 8))
        c.push(uint8(tmp & 0xff))
        c.jump(dst)
    case 0x24:  // BIT [ZPG]
        c.bitTest(c.zeroPageAddress())
    case 0x38:  // SEC
        c.r.setCarryFlag()
    case 0x4c:  // JMP ABS
        c.jump(c.absoluteAddress())
    case 0x50:  // BVC
        rel := c.relativeAddress()
        if !c.r.getOverflowFlag() {
            c.jump(c.r.pc + rel)
        }
    case 0x60:  // RTS
        tmpl := c.pop()
        tmph := c.pop()
        dst := uint16(tmph) << 8 | uint16(tmpl)
        c.jump(dst + 1)
    case 0x70:  // BVS
        rel := c.relativeAddress()
        if c.r.getOverflowFlag() {
            c.jump(c.r.pc + rel)
        }
    case 0x78:  // SEI
        c.r.setInterruptDisable()
    case 0x85:  // STA [ZPG]
        c.stReg(c.zeroPageAddress(), c.r.acc)
    case 0x86:  // STX [ZPG]
        c.stReg(c.zeroPageAddress(), c.r.x)
    case 0x8d:  // STA [ABS]
        c.stReg(c.absoluteAddress(), c.r.acc)
    case 0x90:  // BCC
        rel := c.relativeAddress()
        if !c.r.getCarryFlag() {
            c.jump(c.r.pc + rel)
        }
    case 0x9a:  // TXS
        c.transfer(&c.r.sp, c.r.x)
    case 0xa2:  // LDX IMM
        c.ldReg(c.immediateAddress(), &c.r.x)
    case 0xa9:  // LDA IMM
        c.ldReg(c.immediateAddress(), &c.r.acc)
    case 0xad:  // LDA [ABS]
        c.ldReg(c.absoluteAddress(), &c.r.acc)
    case 0xb0:  // BCS
        rel := c.relativeAddress()
        if c.r.getCarryFlag() {
            c.jump(c.r.pc + rel)
        }
    case 0xd0:  // BNE
        rel := c.relativeAddress()
        if !c.r.getZeroFlag() {
            c.jump(c.r.pc + rel)
        }
    case 0xd8:  // CLD
        c.r.clearDecimalMode()
    case 0xea:  // NOP 
        break
    case 0xf0:  // BEQ
        rel := c.relativeAddress()
        if c.r.getZeroFlag() {
            c.jump(c.r.pc + rel)
        }
    default:
        log.Printf("Unimplemented opcode 0x%02x\n", opc)
        log.Fatalf("CPU State:\n%s\n", c)
    }
    //log.Printf("CPU State:\n%s\n", c)
    return
}
