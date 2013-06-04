package nes

import (
	"log"
    "time"
    "fmt"
)

type Ppu struct {
    // Resolution: 256x240
    // 341 pixels / scanline
    // 262 scanlines / frame
    lastVB time.Time
    cycleCount uint32
    vint bool
}

const VBLANK_PERIOD = 16670000

func (p *Ppu) Init () (e error) {
	log.Println("Initializing PPU...")
    Ram.RegisterCallback(0x2002, p.StatusRegCb)

    p.lastVB = time.Now()
    return
}

func (p Ppu) String() (out string) {
    out += fmt.Sprintf("CC = %04x, VINT: %v\n", 
        p.cycleCount, p.vint)
    return
}

func (p *Ppu) StatusRegCb (addr uint16, isRead bool) (e error) {
    log.Printf("Accessing %04x\n", addr)

    /*
    if now := time.Now(); now.Sub(p.lastVB) >= VBLANK_PERIOD {
        p.lastVB = now
    }
    */

    // bit 7 is reset to 0 on reads, as are $2005 and $2006
    if isRead {
        val := Ram.Raw[addr]
        Ram.Raw[addr] = (val & 0x7f)
        Ram.Raw[0x2005] = 0x0
        Ram.Raw[0x2006] = 0x0
    }
    return
}

func (p *Ppu) Step () (e error) {
    if p.cycleCount == 341 * 262 {
        p.cycleCount = 0
    }
    scanlines := p.cycleCount / 341
    //pixelOfs := p.cycleCount % 341
    if scanlines >= 0 && scanlines < 20 {
        // VINT period
    } else if scanlines == 20 {
        // first scanline is dummy
        p.vint = false
        Ram.Raw[0x2002] &= 0x7f
    } else if scanlines > 20 && scanlines <= 260 {
    } else if scanlines > 260 {
        // do nothing, set VINT
        p.vint = true
        Ram.Raw[0x2002] |= 0x80
    } else {
        log.Panic("Scanlines did not reset!")
    }
    p.cycleCount++
    return
}
