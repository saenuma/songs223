package main

import (
	"bytes"
	"image"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func topBarPartOfMouseCallback(window *glfw.Window, widgetCode int) {
	switch widgetCode {
	case OpenWDBtn:
		rootPath, _ := GetRootPath()
		ExternalLaunch(rootPath)

	case FoldersViewBtn:
		IsOutsidePlayer = true
		ObjCoords = make(map[int]g143.RectSpecs)
		DrawFirstUI(window, CurrentPage)
		window.SetMouseButtonCallback(mouseBtnCallback)
		window.SetScrollCallback(FirstUIScrollCallback)

	case NowPlayingViewBtn:
		if CurrentPlayingSong.SongName != "" {
			ObjCoords = make(map[int]g143.RectSpecs)
			seconds := time.Since(StartTime).Seconds()

			DrawNowPlayingUI(window, CurrentPlayingSong, int(seconds))
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)
			window.SetScrollCallback(nil)
		}

	case InfoBtn:
		IsOutsidePlayer = true
		ObjCoords = make(map[int]g143.RectSpecs)
		DrawInfoUI(window)
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

	for code, RS := range ObjCoords {
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
		ObjCoords = make(map[int]g143.RectSpecs)
		folderIndex := widgetCode - 2000 - 1
		gottenFolder := GetFolders(CurrentPage)[folderIndex]
		DrawFolderUI(window, gottenFolder)
		window.SetMouseButtonCallback(folderUIMouseBtnCallback)
		window.SetScrollCallback(nil)
	}

	// for generated page buttons
	if widgetCode > 3000 && widgetCode < 4000 {
		ObjCoords = make(map[int]g143.RectSpecs)
		pageNum := widgetCode - 3000
		DrawFirstUI(window, pageNum)
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

	for code, RS := range ObjCoords {
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
		ObjCoords = make(map[int]g143.RectSpecs)

		songIndex := widgetCode - 4001
		songDesc := GetSongs(CurrentSongFolder)[songIndex]
		DrawNowPlayingUI(window, songDesc, 0)
		window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

		StartTime = time.Now()
		go playAudio(songDesc.SongPath, "00:00:00")
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

	for code, RS := range ObjCoords {
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

		ObjCoords = make(map[int]g143.RectSpecs)

		songs := GetSongs(CurrentSongFolder)
		var songIndex int
		for index, songDesc := range songs {
			if songDesc.SongName == CurrentPlayingSong.SongName {
				songIndex = index
				break
			}
		}
		if songIndex != 0 {
			songDesc := GetSongs(CurrentSongFolder)[songIndex-1]
			DrawNowPlayingUI(window, songDesc, 0)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			StartTime = time.Now()
			go playAudio(songDesc.SongPath, "00:00:00")
		} else {
			IsOutsidePlayer = true
			DrawFolderUI(window, CurrentSongFolder)
			window.SetMouseButtonCallback(folderUIMouseBtnCallback)
		}

	case PlayPauseBtn:
		if playerCancelFn != nil {
			playerCancelFn()
			playerCancelFn = nil
			seconds := time.Since(StartTime).Seconds()
			PausedSeconds = int(seconds)

			playImg, _, _ := image.Decode(bytes.NewReader(PlayBytes))
			playImg = imaging.Fit(playImg, BoxSize, BoxSize, imaging.Lanczos)

			ggCtx := gg.NewContextForImage(TmpNowPlayingImg)
			ggCtx.SetHexColor("#fff")
			ggCtx.DrawRectangle(float64(widgetRS.OriginX), float64(widgetRS.OriginY), BoxSize, BoxSize)
			ggCtx.Fill()
			ggCtx.DrawImage(playImg, widgetRS.OriginX, widgetRS.OriginY)

			// send the frame to glfw window
			windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
			g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
			window.SwapBuffers()
		} else {
			ObjCoords = make(map[int]g143.RectSpecs)
			DrawNowPlayingUI(window, CurrentPlayingSong, PausedSeconds)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			startTimeUnix := time.Now().Unix() - int64(PausedSeconds)
			StartTime = time.Unix(startTimeUnix, 0)

			go playAudio(CurrentPlayingSong.SongPath, "00:"+SecondsToMinutes(PausedSeconds))
		}

	case NextBtn:
		if playerCancelFn != nil {
			playerCancelFn()
		}

		ObjCoords = make(map[int]g143.RectSpecs)

		songs := GetSongs(CurrentSongFolder)
		var songIndex int
		for index, songDesc := range songs {
			if songDesc.SongName == CurrentPlayingSong.SongName {
				songIndex = index
				break
			}
		}
		if songIndex != len(songs)-1 {
			songDesc := GetSongs(CurrentSongFolder)[songIndex+1]
			DrawNowPlayingUI(window, songDesc, 0)
			window.SetMouseButtonCallback(nowPlayingMouseBtnCallback)

			StartTime = time.Now()
			go playAudio(songDesc.SongPath, "00:00:00")
		} else {
			IsOutsidePlayer = true
			DrawFolderUI(window, CurrentSongFolder)
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

	for code, RS := range ObjCoords {
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
	case Lyrics818Link:
		ExternalLaunch("https://sae.ng/lyrics818")

	case SaeNgLink:
		ExternalLaunch("https://sae.ng")
	}
}
