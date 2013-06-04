package nes

import (
    "fmt"
    "log"
    "io/ioutil"
)

type Rom struct {
    romBankCnt uint8
    RomBanks [2][]byte
    vromBankCnt uint8
    vromBanks [][]byte
    flags1 byte
    flags2 byte
    ramBankCnt uint8
    isPal uint8
    reserved [6]byte
}

func (r Rom) String() string {
    out := fmt.Sprintf("%d ROM banks\n%d VROM banks\n%d RAM banks\n",
        r.romBankCnt, r.vromBankCnt, r.ramBankCnt)
    out += fmt.Sprintf("%d Flags1\n%d Flags2\n%d isPal\n",
        r.flags1, r.flags2, r.isPal)
    return out
}

func (r *Rom) Init(infile string) error {
	log.Println("Initializing ROM...")

    file, err := ioutil.ReadFile(infile)
	if err != nil {
		fmt.Println(err)
        return err
	}
    log.Printf("ROM is 0x%x bytes\n", len(file))

    if string(file[:3]) != "NES" {
        log.Fatal("Malformed NES ROM file!")
    }

    r.romBankCnt = uint8(file[4])
    r.vromBankCnt = uint8(file[5])
    r.flags1 = uint8(file[6])
    r.flags2 = uint8(file[7])
    r.ramBankCnt = uint8(file[8])
    r.isPal = uint8(file[9])

    // Allocate and fill ROM banks
    //r.RomBanks = make([][]byte, 2)
    for i := 0; i < int(r.romBankCnt); i++ {
        offset := 0x10 + i * 0x4000
        bank := make([]byte, 0x4000)
        copy(bank, file[offset:offset + 0x4000])
        r.RomBanks[i] = bank

        // Initialize RAM
        copy(Ram.Raw[0x8000+i*0x4000:], bank)
    }
    if r.romBankCnt == 1 {
        r.RomBanks[1] = make([]byte, 0x4000)
        copy(r.RomBanks[1], r.RomBanks[0])
        copy(Ram.Raw[0xc000:], r.RomBanks[1])
    }

    // Allocate and fille VROM banks
    for i := range r.vromBanks {
        offset := 0x10 + (i + int(r.romBankCnt)) * 0x2000
        bank := make([]byte, 0x2000)
        copy(bank, file[offset:offset + 0x2000])
        r.vromBanks[i] = bank
    }

    return nil
}
