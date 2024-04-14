package main

import (
	"os"
	"path/filepath"
	"strings"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/l8f"
)

var currentSongFolder SongFolder

func getSongs(songFolder SongFolder) []SongDesc {
	rootPath, _ := GetRootPath()
	currentFolderPath := filepath.Join(rootPath, songFolder.Title)
	dirEs, _ := os.ReadDir(currentFolderPath)

	songs := make([]SongDesc, 0)
	for _, dirE := range dirEs {
		if !strings.HasSuffix(dirE.Name(), ".l8f") {
			continue
		}

		songPath := filepath.Join(currentFolderPath, dirE.Name())
		songName := strings.ReplaceAll(dirE.Name(), ".l8f", "")
		songLengthSeconds, _ := l8f.GetVideoLength(songPath)
		songLengthStr := SecondsToMinutes(songLengthSeconds)

		songs = append(songs, SongDesc{SongName: songName, SongPath: songPath, Length: songLengthStr})
	}

	return songs
}
func drawFolderUI(window *glfw.Window, songFolder SongFolder) {
	wWidth, wHeight := window.GetSize()

	currentSongFolder = songFolder

	ggCtx := drawTopBar(window)

	coverW := 300
	songCoverImg, _ := imaging.Open(songFolder.Cover)
	songCoverImg = imaging.Fit(songCoverImg, coverW, coverW, imaging.Lanczos)
	ggCtx.DrawImage(songCoverImg, 40, 80)

	songsX := coverW + 40 + 20

	fontPath := getDefaultFontPath()
	ggCtx.LoadFontFace(fontPath, 40)
	ggCtx.SetHexColor("#444")
	ggCtx.DrawString(songFolder.Title, float64(songsX), 80+fontSize+20)

	ggCtx.LoadFontFace(fontPath, 20)

	// songs UI
	songs := getSongs(songFolder)
	currentY := 80 + 40 + 30
	for i, songDesc := range songs {
		ggCtx.SetHexColor("#444")
		ggCtx.DrawString(songDesc.SongName, float64(songsX), float64(currentY)+fontSize)

		ggCtx.SetHexColor("#888")
		sLW, _ := ggCtx.MeasureString(songDesc.Length)
		ggCtx.DrawString(songDesc.Length, float64(wWidth)-sLW-40, float64(currentY)+fontSize)

		aSongRS := g143.NRectSpecs(songsX, currentY, wWidth-40, fontSize)
		objCoords[4000+i+1] = aSongRS
		currentY += 60
	}

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}

func folderUiMouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
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

	// for generated page buttons
	if widgetCode > 4000 && widgetCode < 5000 {
		objCoords = make(map[int]g143.RectSpecs)

		songIndex := widgetCode - 4001
		songDesc := getSongs(currentSongFolder)[songIndex]
		drawNowPlayingUI(window, songDesc)
		window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

		go playAudio(songDesc.SongPath, "00:00:00")
	}

}
