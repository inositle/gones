package nes

type Memory [0x10000]byte

type MemoryError struct {
    ErrorText string
}

var Ram Memory

func (e MemoryError) Error() string {
    return e.ErrorText
}

func (m *Memory) Init() {
    return
}

func (m *Memory) Read(addr uint16) (b byte, err error) {
    if addr > 0x8000 && addr <= 0xffff {
        b = m[addr]
    } else {
        err = MemoryError{ErrorText:"Unimplemented memory access!"}
    }
    return
}
