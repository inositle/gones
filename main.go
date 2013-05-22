package main

import (
	"flag"
	"github.com/inositle/gones/nes"
	"log"
)

var (
	infile  string
	bDisasm bool

	rom nes.Rom
	cpu nes.Cpu
	apu nes.Apu
)

func init() {
	flag.StringVar(&infile, "r", "", "Rom to emulate")
	flag.BoolVar(&bDisasm, "d", false, "Disassemble instead of execute")
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

	if bDisasm {
		for {
            nes.Disassemble(&cpu, true)
		}
	} else {
        for {
            nes.Disassemble(&cpu, false)
            cpu.Step()
        }
    }
}
