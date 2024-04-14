package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	fps      = 10
	fontSize = 20
	pageSize = 8

	FoldersViewBtn    = 101
	NowPlayingViewBtn = 102
	OpenWDBtn         = 103
	InfoBtn           = 104
)

var objCoords map[int]g143.RectSpecs
var currentPage int

func main() {
	runtime.LockOSThread()

	GetRootPath()
	objCoords = make(map[int]g143.RectSpecs)

	window := g143.NewWindow(1200, 800, "Songs223: media player of songs with embedded lyrics", false)
	drawFirstUI(window, 1)

	// respond to the mouse
	window.SetMouseButtonCallback(mouseBtnCallback)

	for !window.ShouldClose() {
		t := time.Now()
		glfw.PollEvents()

		time.Sleep(time.Second/time.Duration(fps) - time.Since(t))
	}
}

func totalPages() int {
	rootPath, _ := GetRootPath()
	dirFIs, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	return int(math.Ceil(float64(len(dirFIs)) / float64(pageSize)))
}

func getFolders(page int) []SongFolder {
	rootPath, _ := GetRootPath()
	ret := make([]SongFolder, 0)

	dirFIs, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err.Error())
		return ret
	}

	noCoverPath := filepath.Join(os.TempDir(), "no_cover.png")
	os.WriteFile(noCoverPath, NoCover, 0777)

	beginIndex := (page - 1) * pageSize
	endIndex := beginIndex + pageSize

	var toCheckDirFIs []fs.DirEntry
	if endIndex > len(dirFIs) {
		toCheckDirFIs = dirFIs[beginIndex:]
	} else {
		toCheckDirFIs = dirFIs[beginIndex:endIndex]
	}

	for _, dirFI := range toCheckDirFIs {
		if !dirFI.IsDir() {
			continue
		}

		coverPath := noCoverPath
		if DoesPathExists(filepath.Join(rootPath, dirFI.Name(), "cover.jpg")) {
			coverPath = filepath.Join(rootPath, dirFI.Name(), "cover.jpg")
		} else if DoesPathExists(filepath.Join(rootPath, dirFI.Name(), "Cover.jpg")) {
			coverPath = filepath.Join(rootPath, dirFI.Name(), "Cover.jpg")
		}

		innerDirFIs, err := os.ReadDir(filepath.Join(rootPath, dirFI.Name()))
		if err != nil {
			fmt.Println(err)
			continue
		}

		l8fCount := 0

		for _, innerDirFI := range innerDirFIs {
			if strings.HasSuffix(innerDirFI.Name(), ".l8f") {
				l8fCount += 1
				continue
			}
		}

		ret = append(ret, SongFolder{dirFI.Name(), coverPath, l8fCount})
	}

	return ret
}

func drawTopBar(window *glfw.Window) *gg.Context {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := getDefaultFontPath()
	err := ggCtx.LoadFontFace(fontPath, 20)
	if err != nil {
		panic(err)
	}

	// folders button
	foldersStr := "Folders"
	foldersStrW, foldersStrH := ggCtx.MeasureString(foldersStr)
	foldersBtnW := foldersStrW + 80
	foldersBtnH := foldersStrH + 30
	ggCtx.SetHexColor("#B75F5F")
	foldersBtnX := 200
	ggCtx.DrawRoundedRectangle(float64(foldersBtnX), 10, foldersBtnW, foldersBtnH, foldersBtnH/2)
	ggCtx.Fill()

	foldersBtnRS := g143.NRectSpecs(foldersBtnX, 10, int(foldersBtnW), int(foldersBtnH))
	objCoords[FoldersViewBtn] = foldersBtnRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(foldersStr, float64(20+foldersBtnX), 10+foldersStrH+15)

	ggCtx.SetHexColor("#633232")
	ggCtx.DrawCircle(10+float64(foldersBtnX)+foldersBtnW-40, 10+foldersBtnH/2, 10)
	ggCtx.Fill()

	// Now Playing Button
	npStr := "Now Playing"
	npStrW, npStrH := ggCtx.MeasureString(npStr)
	npBtnW := npStrW + 80
	npBtnH := npStrH + 30
	npBtnX := foldersBtnW + float64(foldersBtnRS.OriginX) + 20
	ggCtx.SetHexColor("#81577F")
	ggCtx.DrawRoundedRectangle(npBtnX, 10, npBtnW, npBtnH, npBtnH/2)
	ggCtx.Fill()

	npRS := g143.NRectSpecs(int(npBtnX), 10, int(npBtnW), int(npBtnH))
	objCoords[NowPlayingViewBtn] = npRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(npStr, 30+npBtnX, 10+npStrH+15)

	ggCtx.SetHexColor("#633260")
	ggCtx.DrawCircle(float64(npRS.OriginX)+npBtnW-30, 10+npBtnH/2, 10)
	ggCtx.Fill()

	// Open Working Directory button
	owdStr := "Open Working Directory"
	owdStrWidth, owdStrHeight := ggCtx.MeasureString(owdStr)
	openWDBtnWidth := owdStrWidth + 60
	openWDBtnHeight := owdStrHeight + 30
	ggCtx.SetHexColor("#56845A")
	openWDBtnOriginX := float64(npRS.OriginX+npRS.Width) + 20
	ggCtx.DrawRoundedRectangle(openWDBtnOriginX, 10, openWDBtnWidth, openWDBtnHeight, openWDBtnHeight/2)
	ggCtx.Fill()

	openWDBtnRS := g143.RectSpecs{Width: int(openWDBtnWidth), Height: int(openWDBtnHeight),
		OriginX: int(openWDBtnOriginX), OriginY: 10}
	objCoords[OpenWDBtn] = openWDBtnRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(owdStr, 30+float64(openWDBtnRS.OriginX), 10+owdStrHeight+15)

	// Render button
	ifStr := "Info"
	ifStrW, ifStrH := ggCtx.MeasureString(ifStr)
	ifBtnW := ifStrW + 60
	ifBtnH := ifStrH + 30
	ggCtx.SetHexColor("#B19644")
	renderBtnX := openWDBtnRS.OriginX + openWDBtnRS.Width + 20
	ggCtx.DrawRoundedRectangle(float64(renderBtnX), 10, ifBtnW, ifBtnH, ifBtnH/2)
	ggCtx.Fill()

	rbRS := g143.RectSpecs{OriginX: renderBtnX, OriginY: 10, Width: int(ifBtnW),
		Height: int(ifBtnH)}
	objCoords[InfoBtn] = rbRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(ifStr, float64(rbRS.OriginX)+30, 10+ifStrH+15)
	// draw end of topbar demarcation
	ggCtx.SetHexColor("#999")
	ggCtx.DrawRectangle(10, float64(openWDBtnRS.OriginY+openWDBtnRS.Height+10), float64(wWidth)-20, 2)
	ggCtx.Fill()

	return ggCtx
}

func drawFirstUI(window *glfw.Window, page int) {
	currentPage = page
	wWidth, wHeight := window.GetSize()

	ggCtx := drawTopBar(window)

	songFolders := getFolders(page)

	gutter := 40
	currentX := gutter
	currentY := 80

	// album arts
	boxDimension := 250
	for i, songFolder := range songFolders {
		songCoverImg, _ := imaging.Open(songFolder.Cover)
		songCoverImg = imaging.Fit(songCoverImg, boxDimension, boxDimension, imaging.Lanczos)

		ggCtx.DrawImage(songCoverImg, currentX, currentY)
		ggCtx.SetHexColor("#444")
		songCountStr := fmt.Sprintf("(%d songs)", songFolder.NumberOfSongs)
		ggCtx.DrawString(songFolder.Title, float64(currentX)+20, float64(currentY)+fontSize+float64(boxDimension))
		ggCtx.DrawString(songCountStr, float64(currentX)+20, float64(currentY)+fontSize*2+float64(boxDimension))

		aSongRS := g143.NRectSpecs(currentX, currentY, boxDimension, boxDimension+50)
		objCoords[2000+i+1] = aSongRS

		newX := currentX + boxDimension + gutter + 20
		if newX > (wWidth - boxDimension) {
			currentY += boxDimension + gutter + 20
			currentX = gutter
		} else {
			currentX += boxDimension + gutter
		}
	}

	// paging

	aPageCurrentX := 40
	aPageCurrentY := 720
	aPageGutter := 10
	for i := 1; i <= totalPages(); i++ {
		aStr := fmt.Sprintf("%d", i)
		aStrW, aStrH := ggCtx.MeasureString(aStr)
		aPageBtnW := aStrW + 10
		aPageBtnH := aStrH + 10

		ggCtx.SetHexColor("#633232")
		ggCtx.DrawRoundedRectangle(float64(aPageCurrentX), float64(aPageCurrentY), aPageBtnW, aPageBtnH, 5)
		ggCtx.Fill()

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawString(aStr, 5+float64(aPageCurrentX), float64(aPageCurrentY)+fontSize)

		aPageBtnRS := g143.NRectSpecs(aPageCurrentX, aPageCurrentY, int(aPageBtnW), int(aPageBtnH))
		objCoords[3000+i] = aPageBtnRS
		newX := aPageCurrentX + int(aPageBtnW) + aPageGutter
		if newX > (wWidth - int(aPageBtnW)) {
			currentY += int(aPageBtnW) + aPageGutter
			aPageCurrentX = 40
		} else {
			aPageCurrentX += int(aPageBtnW) + aPageGutter
		}
	}

	// send the frame to glfw window
	windowRS := g143.RectSpecs{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()
}

func getDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "s223_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
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

	// for generated folder buttons
	if widgetCode > 2000 && widgetCode < 3000 {
		folderIndex := widgetCode - 2000 - 1
		gottenFolder := getFolders(currentPage)[folderIndex]
		drawFolderUI(window, gottenFolder)
		window.SetMouseButtonCallback(folderUiMouseBtnCallback)
	}

	// for generated page buttons
	if widgetCode > 3000 && widgetCode < 4000 {
		pageNum := widgetCode - 3000
		objCoords = make(map[int]g143.RectSpecs)
		drawFirstUI(window, pageNum)
	}

}

func topBarPartOfMouseCallback(window *glfw.Window, widgetCode int) {
	switch widgetCode {
	case OpenWDBtn:
		rootPath, _ := GetRootPath()
		externalLaunch(rootPath)

	case FoldersViewBtn:
		drawFirstUI(window, currentPage)
		window.SetMouseButtonCallback(mouseBtnCallback)
	}

}
