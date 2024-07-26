package main

import (
	"math"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/l8f"
)

func main() {
	runtime.LockOSThread()

	GetRootPath()
	ObjCoords = make(map[int]g143.RectSpecs)

	window := g143.NewWindow(1200, 800, "Songs223: media player of songs with embedded lyrics", false)
	DrawFirstUI(window, 1)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	// respond to mouse movements
	window.SetCursorPosCallback(curPosCB)

	window.SetScrollCallback(FirstUIScrollCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		// update UI if song is playing
		if CurrentPlayingSong.SongName != "" && !IsOutsidePlayer && currentPlayer != nil && currentPlayer.IsPlaying() {
			seconds := time.Since(StartTime).Seconds()
			secondsInt := int(math.Floor(seconds))
			if secondsInt != CurrentPlaySeconds {
				DrawNowPlayingUI(window, CurrentPlayingSong, secondsInt)
			}
		}

		// play next song or stop
		if CurrentPlayingSong.SongName != "" && !IsOutsidePlayer {
			songLengthSeconds, _ := l8f.GetVideoLength(CurrentPlayingSong.SongPath)
			if songLengthSeconds == int(time.Since(StartTime).Seconds()) {
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
					go playAudio(songDesc.SongPath)
				} else {
					IsOutsidePlayer = true
					DrawFolderUI(window, CurrentSongFolder)
					window.SetMouseButtonCallback(folderUIMouseBtnCallback)
				}
			}
		}

		time.Sleep(time.Second/time.Duration(FPS) - time.Since(t))
	}
}
