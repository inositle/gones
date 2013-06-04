package nes

import (
    "log"
)

type callback func(addr uint16, isRead bool) error

type Memory struct {
    Raw [0x10000]byte
    callbacks map[uint16] callback
    inCallback bool
}

type MemoryError struct {
    ErrorText string
}

var Ram Memory

func (e MemoryError) Error() string {
    return e.ErrorText
}

func init() {
    log.Printf("Initializing RAM...\n")
    Ram.callbacks = make(map[uint16] callback)
}

func (m *Memory) RegisterCallback (addr uint16, f callback) {
    log.Printf("Registering callback @ 0x%04x: %v\n", addr, f)
    m.callbacks[addr] = f
}

func (m *Memory) Read(addr uint16) (b byte, err error) {
    // TODO: Raise "exception" for invalid memory regions
    if m.inCallback {
        log.Fatal("Do not call Read/Write from within a callback!")
    }

    if addr >= 0x0 && addr <= 0xffff {
        b = m.Raw[addr]
    } else {
        err = MemoryError{ErrorText:"Unimplemented memory access!"}
    }
    // Post Callbacks
    if cb, ok := m.callbacks[addr]; ok {
        m.inCallback = true
        cb(addr, true)
        m.inCallback = false
    }
    return
}

func (m *Memory) Write(addr uint16, b byte) (err error) {
    if m.inCallback {
        log.Fatal("Do not call Read/Write from within a callback!")
    }

    if addr >= 0x0 && addr <= 0xffff {
        m.Raw[addr] = b
    } else {
        err = MemoryError{ErrorText:"Unimplemented memory access!"}
    }
    // Post Callbacks
    if cb, ok := m.callbacks[addr]; ok {
        m.inCallback = true
        cb(addr, false)
        m.inCallback = false
    }
    return
}
