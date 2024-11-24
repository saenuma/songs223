package internal

import (
	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func DrawInfoUI(window *glfw.Window) {
	wWidth, wHeight := window.GetSize()

	ggCtx := DrawTopBar(window)

	fontPath := GetDefaultFontPath()
	ggCtx.LoadFontFace(fontPath, 40)
	ggCtx.SetHexColor("#444")
	infoY := 80 + 20
	ggCtx.DrawString("Info Page", 40, float64(infoY)+FontSize)

	ggCtx.LoadFontFace(fontPath, 20)

	msg1 := "To make a song that songs223 supports, please use "

	ggCtx.DrawString(msg1, 60, float64(infoY)+60+FontSize)
	ggCtx.SetHexColor("#444")
	lStrW, lStrH := ggCtx.MeasureString("lyrics818")
	lStrY := infoY + 60 + FontSize + 10
	ggCtx.DrawRoundedRectangle(60, float64(lStrY), lStrW+20, lStrH+20, 4)
	ObjCoords[Lyrics818Link] = g143.NewRect(60, lStrY, int(lStrW)+20, int(lStrH)+20)
	ggCtx.Fill()

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString("lyrics818", 70, float64(lStrY)+FontSize+5)

	msg2 := "Brought to you with love by "
	msg2Y := lStrY + 100
	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(msg2, 60, float64(msg2Y))

	sStr := "https://sae.ng"
	sStrW, sStrH := ggCtx.MeasureString(sStr)
	sStrY := msg2Y + 10
	ggCtx.SetHexColor("#444")
	ggCtx.DrawRoundedRectangle(60, float64(sStrY), sStrW+20, sStrH+20, 4)
	ggCtx.Fill()
	ObjCoords[SaeNgLink] = g143.NewRect(60, sStrY, int(sStrW)+20, int(sStrH)+20)
	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(sStr, 70, float64(sStrY)+FontSize+5)

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
}
