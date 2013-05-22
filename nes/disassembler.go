package nes

import (
	"bytes"
	"encoding/binary"
	//"log"
    "fmt"
)

func absoluteAddress(pc uint16) (result uint16) {
	buf := bytes.NewBuffer(Ram[pc : pc+2])
	binary.Read(buf, binary.LittleEndian, &result)
	return
}

func immediateAddress(pc uint16) (result uint8) {
	buf := bytes.NewBuffer(Ram[pc : pc+1])
	binary.Read(buf, binary.LittleEndian, &result)
	return
}

func relativeAddress(pc uint16) (result int8) {
    return int8(immediateAddress(pc))
}

func Disassemble(cpu *Cpu, update bool) (new_pc uint16) {
    pc := cpu.r.pc
    opcode, _ := Ram.Read(pc)
	new_pc = pc + 1

    out := fmt.Sprintf("%04x: %02x  ", pc, opcode)
	switch opcode {
    case 0x00:
        out += fmt.Sprintf("BRK\n")
    case 0x09:
        out += fmt.Sprintf("ORA #%02x\n", immediateAddress(pc+1))
        new_pc += 1
    case 0x10:
        out += fmt.Sprintf("BPL $%02x\n", relativeAddress(pc+1))
        new_pc += 1
    case 0x20:
        out += fmt.Sprintf("JSR $%04x\n", absoluteAddress(pc+1))
        new_pc += 2
	case 0x4c:
		out += fmt.Sprintf("JMP $%04x\n", absoluteAddress(pc+1))
		new_pc += 2
	case 0x60:
		out += fmt.Sprintf("RTS\n")
	case 0x78:
		out += fmt.Sprintf("SEI\n")
    case 0x8d:
        out += fmt.Sprintf("STA [$%04x]\n", absoluteAddress(pc+1))
        new_pc += 2
    case 0x90:
        out += fmt.Sprintf("BCC $%02x\n", relativeAddress(pc+1))
        new_pc += 1
	case 0x9a:
		out += fmt.Sprintf("TXS\n")
	case 0xa0:
		out += fmt.Sprintf("LDY #%02x\n", immediateAddress(pc+1))
		new_pc += 1
	case 0xa2:
		out += fmt.Sprintf("LDX #%02x\n", immediateAddress(pc+1))
		new_pc += 1
    case 0xa9:
        out += fmt.Sprintf("LDA #%02x\n", immediateAddress(pc+1))
        new_pc += 1
	case 0xad:
		out += fmt.Sprintf("LDA [$%04x]\n", absoluteAddress(pc+1))
		new_pc += 2
    case 0xb0:
        out += fmt.Sprintf("BCS $%02x\n", relativeAddress(pc+1))
        new_pc += 1
    case 0xbd:
        out += fmt.Sprintf("LDA [$%04x + X]\n", absoluteAddress(pc+1))
        new_pc += 2
    case 0xca:
        out += fmt.Sprintf("DEX\n")
    case 0xc9:
        out += fmt.Sprintf("CMP #%02x\n", immediateAddress(pc+1))
        new_pc += 1
    case 0xd0:
        out += fmt.Sprintf("BNE $%02x\n", relativeAddress(pc+1))
        new_pc += 1
	case 0xd8:
		out += fmt.Sprintf("CLD\n")
    case 0xee:
        out += fmt.Sprintf("INC [$%04x]\n", absoluteAddress(pc+1))
        new_pc += 2
	default:
		out += fmt.Sprintf("Unknown: 0x%02x\n", opcode)
	}

    fmt.Printf(out)
    if update {
        cpu.r.pc = new_pc
    }

	return
}
