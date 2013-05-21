package nes

const (
	RW1          = 0x4000
	RW2          = 0x4004
	TriWave      = 0x4008
	Noise        = 0x400c
	DmcPlayMode  = 0x4010
	DmcDeltaCnt  = 0x4011
	DmcPlayCode  = 0x4012
	DmcCodeLen   = 0x4013
	ChanEnable   = 0x4015
	FrameCntCtrl = 0x4017
)

type Apu struct {
}
