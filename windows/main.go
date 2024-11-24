package main

import (
	"math"
	"runtime"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/saenuma/lyrics818/l8f"
	"github.com/saenuma/songs223a/internal"
)

func main() {
	runtime.LockOSThread()

	internal.GetRootPath()
	internal.ObjCoords = make(map[int]g143.Rect)

	window := g143.NewWindow(1200, 800, "Songs223: media player of songs with embedded lyrics", false)
	internal.DrawFirstUI(window, 1)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetScrollCallback(internal.FirstUIScrollCallback)
	window.SetCursorPosCallback(internal.CurPosCB)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		// update UI if song is playing
		if internal.CurrentPlayingSong.SongName != "" && !internal.IsOutsidePlayer && currentPlayer != nil && currentPlayer.IsPlaying() {
			seconds := time.Since(internal.StartTime).Seconds()
			secondsInt := int(math.Floor(seconds))
			if secondsInt != internal.CurrentPlaySeconds {
				internal.DrawNowPlayingUI(window, internal.CurrentPlayingSong, secondsInt)
			}
		}

		// play next song or stop
		if internal.CurrentPlayingSong.SongName != "" && !internal.IsOutsidePlayer {
			songLengthSeconds, _ := l8f.GetVideoLength(internal.CurrentPlayingSong.SongPath)
			if songLengthSeconds == int(time.Since(internal.StartTime).Seconds()) {
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

		time.Sleep(time.Second/time.Duration(internal.FPS) - time.Since(t))
	}
}
