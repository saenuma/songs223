package main

import (
	"image"
	"time"

	g143 "github.com/bankole7782/graphics143"
)

type SongFolder struct {
	Title         string
	Cover         string
	NumberOfSongs int
}

type SongDesc struct {
	SongName string
	SongPath string
	Length   string
}

const (
	FPS      = 24
	FontSize = 20
	PageSize = 8

	FoldersViewBtn    = 101
	NowPlayingViewBtn = 102
	OpenWDBtn         = 103
	InfoBtn           = 104

	Scale   = 0.8
	BoxSize = 40

	PlayPauseBtn = 501
	PrevBtn      = 502
	NextBtn      = 503

	Lyrics818Link = 601
	SaeNgLink     = 602
)

var (
	ObjCoords          map[int]g143.RectSpecs
	CurrentPage        int
	IsOutsidePlayer    bool
	scrollEventCount   = 0
	cursorEventsCount  int
	currentWindowFrame image.Image

	CurrentSongFolder SongFolder
	StartTime         time.Time

	CurrentPlayingSong SongDesc
	PausedSeconds      int
	CurrentPlaySeconds int
	TmpNowPlayingImg   image.Image
)
