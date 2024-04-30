package main

import (
	"bytes"
	"image"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/songs223a/internal"
)

func topBarPartOfMouseCallback(window *glfw.Window, widgetCode int) {
	switch widgetCode {
	case internal.OpenWDBtn:
		rootPath, _ := internal.GetRootPath()
		internal.ExternalLaunch(rootPath)

	case internal.FoldersViewBtn:
		internal.IsOutsidePlayer = true
		internal.ObjCoords = make(map[int]g143.RectSpecs)
		internal.DrawFirstUI(window, internal.CurrentPage)
		window.SetMouseButtonCallback(mouseBtnCallback)
		window.SetScrollCallback(internal.FirstUIScrollCallback)

	case internal.NowPlayingViewBtn:
		if internal.CurrentPlayingSong.SongName != "" {
			internal.ObjCoords = make(map[int]g143.RectSpecs)
			seconds := time.Since(internal.StartTime).Seconds()

			internal.DrawNowPlayingUI(window, internal.CurrentPlayingSong, int(seconds))
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)
			window.SetScrollCallback(nil)
		}

	case internal.InfoBtn:
		internal.IsOutsidePlayer = true
		internal.ObjCoords = make(map[int]g143.RectSpecs)
		internal.DrawInfoUI(window)
		window.SetMouseButtonCallback(infoUIMouseBtnCallback)
		window.SetScrollCallback(nil)
	}

}

func mouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range internal.ObjCoords {
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

	// for generated folder buttons
	if widgetCode > 2000 && widgetCode < 3000 {
		internal.ObjCoords = make(map[int]g143.RectSpecs)
		folderIndex := widgetCode - 2000 - 1
		gottenFolder := internal.GetFolders(internal.CurrentPage)[folderIndex]
		internal.DrawFolderUI(window, gottenFolder)
		window.SetMouseButtonCallback(folderUIMouseBtnCallback)
		window.SetScrollCallback(nil)
	}

	// for generated page buttons
	if widgetCode > 3000 && widgetCode < 4000 {
		internal.ObjCoords = make(map[int]g143.RectSpecs)
		pageNum := widgetCode - 3000
		internal.DrawFirstUI(window, pageNum)
	}

}

func folderUIMouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range internal.ObjCoords {
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
		internal.ObjCoords = make(map[int]g143.RectSpecs)

		songIndex := widgetCode - 4001
		songDesc := internal.GetSongs(internal.CurrentSongFolder)[songIndex]
		internal.DrawNowPlayingUI(window, songDesc, 0)
		window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

		internal.StartTime = time.Now()
		go playAudio(songDesc.SongPath)
	}

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

	for code, RS := range internal.ObjCoords {
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
	case internal.PrevBtn:
		if currentPlayer != nil {
			currentPlayer.Pause()
		}

		internal.ObjCoords = make(map[int]g143.RectSpecs)

		songs := internal.GetSongs(internal.CurrentSongFolder)
		var songIndex int
		for index, songDesc := range songs {
			if songDesc.SongName == internal.CurrentPlayingSong.SongName {
				songIndex = index
				break
			}
		}
		if songIndex != 0 {
			songDesc := internal.GetSongs(internal.CurrentSongFolder)[songIndex-1]
			internal.DrawNowPlayingUI(window, songDesc, 0)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			internal.StartTime = time.Now()
			go playAudio(songDesc.SongPath)
		} else {
			internal.IsOutsidePlayer = true
			internal.DrawFolderUI(window, internal.CurrentSongFolder)
			window.SetMouseButtonCallback(folderUIMouseBtnCallback)
		}

	case internal.PlayPauseBtn:
		if currentPlayer != nil && currentPlayer.IsPlaying() {
			currentPlayer.Pause()
			seconds := time.Since(internal.StartTime).Seconds()
			internal.PausedSeconds = int(seconds)

			playImg, _, _ := image.Decode(bytes.NewReader(internal.PlayBytes))
			playImg = imaging.Fit(playImg, internal.BoxSize, internal.BoxSize, imaging.Lanczos)

			ggCtx := gg.NewContextForImage(internal.TmpNowPlayingImg)
			ggCtx.SetHexColor("#fff")
			ggCtx.DrawRectangle(float64(widgetRS.OriginX), float64(widgetRS.OriginY), internal.BoxSize, internal.BoxSize)
			ggCtx.Fill()
			ggCtx.DrawImage(playImg, widgetRS.OriginX, widgetRS.OriginY)

			// send the frame to glfw window
			windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
			window.SwapBuffers()

		} else if currentPlayer != nil && !currentPlayer.IsPlaying() {
			internal.ObjCoords = make(map[int]g143.RectSpecs)
			internal.DrawNowPlayingUI(window, internal.CurrentPlayingSong, internal.PausedSeconds)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			startTimeUnix := time.Now().Unix() - int64(internal.PausedSeconds)
			internal.StartTime = time.Unix(startTimeUnix, 0)

			go continueAudio()
		}

	case internal.NextBtn:
		if currentPlayer != nil {
			currentPlayer.Pause()
		}

		internal.ObjCoords = make(map[int]g143.RectSpecs)

		songs := internal.GetSongs(internal.CurrentSongFolder)
		var songIndex int
		for index, songDesc := range songs {
			if songDesc.SongName == internal.CurrentPlayingSong.SongName {
				songIndex = index
				break
			}
		}
		if songIndex != len(songs)-1 {
			songDesc := internal.GetSongs(internal.CurrentSongFolder)[songIndex+1]
			internal.DrawNowPlayingUI(window, songDesc, 0)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			internal.StartTime = time.Now()
			go playAudio(songDesc.SongPath)
		} else {
			internal.IsOutsidePlayer = true
			internal.DrawFolderUI(window, internal.CurrentSongFolder)
			window.SetMouseButtonCallback(folderUIMouseBtnCallback)
		}
	}
}

func infoUIMouseBtnCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Release {
		return
	}

	xPos, yPos := window.GetCursorPos()
	xPosInt := int(xPos)
	yPosInt := int(yPos)

	// wWidth, wHeight := window.GetSize()

	// var widgetRS g143.RectSpecs
	var widgetCode int

	for code, RS := range internal.ObjCoords {
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

	switch widgetCode {
	case internal.Lyrics818Link:
		internal.ExternalLaunch("https://sae.ng/lyrics818")

	case internal.SaeNgLink:
		internal.ExternalLaunch("https://sae.ng")
	}
}
