package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed fonts/SourceHanSansCN-Medium.otf
var chineseFontData []byte

var resourceSourceHanSansSCVFTtf = &fyne.StaticResource{
	StaticName:    "SourceHanSansCN-Medium.otf",
	StaticContent: chineseFontData,
}
