package main

import (
	"flag"
	"log"
    "github.com/inositle/gones/nes"
)


var (
    infile string

    rom nes.Rom
    cpu nes.Cpu
    apu nes.Apu
)

func init() {
	flag.StringVar(&infile, "r", "", "Rom to emulate")
}

func main() {
	flag.Parse()
	if infile == "" {
		log.Fatal("ROM file required!")
	}

    // Initialize hardware
    cpu.Init()

    // Load in ROM
    if err := rom.Init(infile); err != nil {
		log.Fatal("Error loading ROM. ", err)
        return
	}
    log.Println(rom)

    pc := uint16(0x8000)
    var b byte
    for {
        b = nes.Ram[pc]
        pc = nes.Disassemble(b, &cpu, pc)
    }
}
