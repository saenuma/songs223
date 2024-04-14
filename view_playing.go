package main

import (
	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/l8f"
)

var currentPlayingSong SongDesc

const scale = 0.8

func drawNowPlayingUI(window *glfw.Window, songDesc SongDesc, seconds int) {
	outsidePlayer = false
	wWidth, wHeight := window.GetSize()

	currentPlayingSong = songDesc

	ggCtx := drawTopBar(window)

	// ggCtx.SetHexColor("#444")
	// ggCtx.DrawString(currentSongFolder.Title+" / "+songDesc.SongName, 200, 80+30)

	// scale down the image
	currFrame, _ := l8f.ReadLaptopFrame(songDesc.SongPath, seconds)
	displayFrameW := int(scale * float64((*currFrame).Bounds().Dx()))
	displayFrameH := int(scale * float64((*currFrame).Bounds().Dy()))
	tmp := imaging.Fit(*currFrame, displayFrameW, displayFrameH, imaging.Lanczos)
	ggCtx.DrawImage(tmp, (wWidth-displayFrameW)/2, 100)

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}

func nowPlayingMouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			// widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	topBarPartOfMouseCallback(window, widgetCode)

}
