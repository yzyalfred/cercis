package utils

import "unsafe"

type GlobalSetting struct {
	littleEndian bool
}

func (this *GlobalSetting) Init() {
	this.littleEndian = this.systemLittleEndian()
}

func (this *GlobalSetting) systemLittleEndian() bool {
	n := 0x1234
	return *(*byte)(unsafe.Pointer(&n)) == 0x34
}

var gGlobalSetting *GlobalSetting

func LittleEndian() bool {
	return gGlobalSetting.littleEndian
}