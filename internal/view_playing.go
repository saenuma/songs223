package internal

import (
	"bytes"
	"image"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/l8f"
)

func DrawNowPlayingUI(window *glfw.Window, songDesc SongDesc, seconds int) {
	wWidth, wHeight := window.GetSize()

	IsOutsidePlayer = false
	CurrentPlayingSong = songDesc
	CurrentPlaySeconds = seconds

	ggCtx := DrawTopBar(window)

	// Scale down the image and write frame
	currFrame, _ := l8f.ReadLaptopFrame(songDesc.SongPath, seconds)
	displayFrameW := int(Scale * float64((*currFrame).Bounds().Dx()))
	displayFrameH := int(Scale * float64((*currFrame).Bounds().Dy()))
	tmp := imaging.Fit(*currFrame, displayFrameW, displayFrameH, imaging.Lanczos)
	ggCtx.DrawImage(tmp, (wWidth-displayFrameW)/2, 80)

	aStr := CurrentSongFolder.Title + " / " + songDesc.SongName
	aStrW, _ := ggCtx.MeasureString(aStr)
	ggCtx.SetHexColor("#444")

	aStrY := float64(displayFrameH) + 90 + FontSize
	ggCtx.DrawString(aStr, (float64(wWidth)-aStrW)/2, aStrY)

	window.SetTitle(aStr + "  | Songs223")

	// write time elapsed
	elapsedTimeStr := SecondsToMinutes(seconds)
	ggCtx.DrawString(elapsedTimeStr, 50, aStrY)

	// write stop time
	totalSeconds, _ := l8f.GetVideoLength(songDesc.SongPath)
	stopTimeStr := SecondsToMinutes(totalSeconds)
	stopTimeStrW, _ := ggCtx.MeasureString(stopTimeStr)
	ggCtx.DrawString(stopTimeStr, float64(wWidth)-50-stopTimeStrW, aStrY)

	// draw controls
	prevImg, _, _ := image.Decode(bytes.NewReader(PrevBytes))
	prevImg = imaging.Fit(prevImg, BoxSize, BoxSize, imaging.Lanczos)
	pauseImg, _, _ := image.Decode(bytes.NewReader(PauseBytes))
	pauseImg = imaging.Fit(pauseImg, BoxSize, BoxSize, imaging.Lanczos)
	nextImg, _, _ := image.Decode(bytes.NewReader(NextBytes))
	nextImg = imaging.Fit(nextImg, BoxSize, BoxSize, imaging.Lanczos)

	controlsY := displayFrameH + 90 + FontSize + 20
	ggCtx.DrawImage(prevImg, 500, controlsY)
	prevRS := g143.NRectSpecs(500, controlsY, BoxSize, BoxSize)
	ObjCoords[PrevBtn] = prevRS

	ggCtx.DrawImage(pauseImg, 600, controlsY)
	pauseRS := g143.NRectSpecs(600, controlsY, BoxSize, BoxSize)
	ObjCoords[PlayPauseBtn] = pauseRS

	ggCtx.DrawImage(nextImg, 700, controlsY)
	nextRS := g143.NRectSpecs(700, controlsY, BoxSize, BoxSize)
	ObjCoords[NextBtn] = nextRS

	// save the frame
	TmpNowPlayingImg = ggCtx.Image()

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
}
