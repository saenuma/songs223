package main

import (
	"bytes"
	"image"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/l8f"
	"github.com/saenuma/songs223a/internal"
)

var currentPlayingSong SongDesc

const (
	scale   = 0.8
	boxSize = 40

	PlayPauseBtn = 501
	PrevBtn      = 502
	NextBtn      = 503
)

var pausedSeconds int
var tmpNowPlayingImg image.Image

func drawNowPlayingUI(window *glfw.Window, songDesc SongDesc, seconds int) {
	outsidePlayer = false
	wWidth, wHeight := window.GetSize()

	currentPlayingSong = songDesc

	ggCtx := drawTopBar(window)

	// scale down the image and write frame
	currFrame, _ := l8f.ReadLaptopFrame(songDesc.SongPath, seconds)
	displayFrameW := int(scale * float64((*currFrame).Bounds().Dx()))
	displayFrameH := int(scale * float64((*currFrame).Bounds().Dy()))
	tmp := imaging.Fit(*currFrame, displayFrameW, displayFrameH, imaging.Lanczos)
	ggCtx.DrawImage(tmp, (wWidth-displayFrameW)/2, 80)

	aStr := currentSongFolder.Title + " / " + songDesc.SongName
	aStrW, _ := ggCtx.MeasureString(aStr)
	ggCtx.SetHexColor("#444")

	aStrY := float64(displayFrameH) + 90 + fontSize
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
	prevImg, _, _ := image.Decode(bytes.NewReader(internal.PrevBytes))
	prevImg = imaging.Fit(prevImg, boxSize, boxSize, imaging.Lanczos)
	pauseImg, _, _ := image.Decode(bytes.NewReader(internal.PauseBytes))
	pauseImg = imaging.Fit(pauseImg, boxSize, boxSize, imaging.Lanczos)
	nextImg, _, _ := image.Decode(bytes.NewReader(internal.NextBytes))
	nextImg = imaging.Fit(nextImg, boxSize, boxSize, imaging.Lanczos)

	controlsY := displayFrameH + 90 + fontSize + 20
	ggCtx.DrawImage(prevImg, 500, controlsY)
	prevRS := g143.NRectSpecs(500, controlsY, boxSize, boxSize)
	objCoords[PrevBtn] = prevRS

	ggCtx.DrawImage(pauseImg, 600, controlsY)
	pauseRS := g143.NRectSpecs(600, controlsY, boxSize, boxSize)
	objCoords[PlayPauseBtn] = pauseRS

	ggCtx.DrawImage(nextImg, 700, controlsY)
	nextRS := g143.NRectSpecs(700, controlsY, boxSize, boxSize)
	objCoords[NextBtn] = nextRS

	// save the frame
	tmpNowPlayingImg = ggCtx.Image()

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

	wWidth, wHeight := window.GetSize()

	var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range objCoords {
		if g143.InRectSpecs(RS, xPosInt, yPosInt) {
			widgetRS = RS
			widgetCode = code
			break
		}
	}

	if widgetCode == 0 {
		return
	}

	topBarPartOfMouseCallback(window, widgetCode)

	switch widgetCode {
	case PrevBtn:
		if playerCancelFn != nil {
			playerCancelFn()
		}

		objCoords = make(map[int]g143.RectSpecs)

		songs := getSongs(currentSongFolder)
		var songIndex int
		for index, songDesc := range songs {
			if songDesc.SongName == currentPlayingSong.SongName {
				songIndex = index
				break
			}
		}
		if songIndex != 0 {
			songDesc := getSongs(currentSongFolder)[songIndex-1]
			drawNowPlayingUI(window, songDesc, 0)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			startTime = time.Now()
			go playAudio(songDesc.SongPath, "00:00:00")
		} else {
			outsidePlayer = true
			drawFolderUI(window, currentSongFolder)
			window.SetMouseButtonCallback(folderUiMouseBtnCallback)
		}

	case PlayPauseBtn:
		if playerCancelFn != nil {
			playerCancelFn()
			playerCancelFn = nil
			seconds := time.Since(startTime).Seconds()
			pausedSeconds = int(seconds)

			playImg, _, _ := image.Decode(bytes.NewReader(internal.PlayBytes))
			playImg = imaging.Fit(playImg, boxSize, boxSize, imaging.Lanczos)

			ggCtx := gg.NewContextForImage(tmpNowPlayingImg)
			ggCtx.SetHexColor("#fff")
			ggCtx.DrawRectangle(float64(widgetRS.OriginX), float64(widgetRS.OriginY), boxSize, boxSize)
			ggCtx.Fill()
			ggCtx.DrawImage(playImg, widgetRS.OriginX, widgetRS.OriginY)

			// send the frame to glfw window
			windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
			window.SwapBuffers()
		} else {
			objCoords = make(map[int]g143.RectSpecs)
			drawNowPlayingUI(window, currentPlayingSong, pausedSeconds)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			startTimeUnix := time.Now().Unix() - int64(pausedSeconds)
			startTime = time.Unix(startTimeUnix, 0)

			go playAudio(currentPlayingSong.SongPath, "00:"+SecondsToMinutes(pausedSeconds))
		}

	case NextBtn:
		if playerCancelFn != nil {
			playerCancelFn()
		}

		objCoords = make(map[int]g143.RectSpecs)

		songs := getSongs(currentSongFolder)
		var songIndex int
		for index, songDesc := range songs {
			if songDesc.SongName == currentPlayingSong.SongName {
				songIndex = index
				break
			}
		}
		if songIndex != len(songs)-1 {
			songDesc := getSongs(currentSongFolder)[songIndex+1]
			drawNowPlayingUI(window, songDesc, 0)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			startTime = time.Now()
			go playAudio(songDesc.SongPath, "00:00:00")
		} else {
			outsidePlayer = true
			drawFolderUI(window, currentSongFolder)
			window.SetMouseButtonCallback(folderUiMouseBtnCallback)
		}
	}
}
