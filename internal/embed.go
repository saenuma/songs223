package internal

import (
	_ "embed"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed no_cover.png
var NoCover []byte

//go:embed play.png
var PlayBytes []byte

//go:embed pause.png
var PauseBytes []byte

//go:embed previous.png
var PrevBytes []byte

//go:embed next.png
var NextBytes []byte
