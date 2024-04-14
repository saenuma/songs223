package main

import (
	_ "embed"
)

//go:embed Roboto-Light.ttf
var DefaultFont []byte

//go:embed no_cover.png
var NoCover []byte
