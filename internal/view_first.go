package internal

import (
	"fmt"
	"os"
	"path/filepath"

	g143 "github.com/bankole7782/graphics143"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func DrawTopBar(window *glfw.Window) *gg.Context {
	wWidth, wHeight := window.GetSize()

	// frame buffer
	ggCtx := gg.NewContext(wWidth, wHeight)

	// background rectangle
	ggCtx.DrawRectangle(0, 0, float64(wWidth), float64(wHeight))
	ggCtx.SetHexColor("#ffffff")
	ggCtx.Fill()

	// load font
	fontPath := GetDefaultFontPath()
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
	foldersBtnX := 280
	ggCtx.DrawRectangle(float64(foldersBtnX), 10, foldersBtnW, foldersBtnH)
	ggCtx.Fill()

	foldersBtnRS := g143.NewRect(foldersBtnX, 10, int(foldersBtnW), int(foldersBtnH))
	ObjCoords[FoldersViewBtn] = foldersBtnRS

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
	ggCtx.DrawRectangle(npBtnX, 10, npBtnW, npBtnH)
	ggCtx.Fill()

	npRS := g143.NewRect(int(npBtnX), 10, int(npBtnW), int(npBtnH))
	ObjCoords[NowPlayingViewBtn] = npRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(npStr, 30+npBtnX, 10+npStrH+15)

	ggCtx.SetHexColor("#633260")
	ggCtx.DrawCircle(float64(npRS.OriginX)+npBtnW-30, 10+npBtnH/2, 10)
	ggCtx.Fill()

	// Open Working Directory button
	owdStr := "Open Folder"
	owdStrWidth, owdStrHeight := ggCtx.MeasureString(owdStr)
	openWDBtnWidth := owdStrWidth + 60
	openWDBtnHeight := owdStrHeight + 30
	ggCtx.SetHexColor("#56845A")
	openWDBtnOriginX := float64(npRS.OriginX+npRS.Width) + 20
	ggCtx.DrawRectangle(openWDBtnOriginX, 10, openWDBtnWidth, openWDBtnHeight)
	ggCtx.Fill()

	openWDBtnRS := g143.Rect{Width: int(openWDBtnWidth), Height: int(openWDBtnHeight),
		OriginX: int(openWDBtnOriginX), OriginY: 10}
	ObjCoords[OpenWDBtn] = openWDBtnRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(owdStr, 30+float64(openWDBtnRS.OriginX), 10+owdStrHeight+15)

	// Render button
	ifStr := "Info"
	ifStrW, ifStrH := ggCtx.MeasureString(ifStr)
	ifBtnW := ifStrW + 60
	ifBtnH := ifStrH + 30
	ggCtx.SetHexColor("#B19644")
	renderBtnX := openWDBtnRS.OriginX + openWDBtnRS.Width + 20
	ggCtx.DrawRectangle(float64(renderBtnX), 10, ifBtnW, ifBtnH)
	ggCtx.Fill()

	rbRS := g143.Rect{OriginX: renderBtnX, OriginY: 10, Width: int(ifBtnW),
		Height: int(ifBtnH)}
	ObjCoords[InfoBtn] = rbRS

	ggCtx.SetHexColor("#fff")
	ggCtx.DrawString(ifStr, float64(rbRS.OriginX)+30, 10+ifStrH+15)

	return ggCtx
}

func DrawFirstUI(window *glfw.Window, page int) {
	CurrentPage = page
	wWidth, wHeight := window.GetSize()

	ggCtx := DrawTopBar(window)

	songFolders := GetFolders(page)

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
		ggCtx.DrawString(songFolder.Title, float64(currentX)+20, float64(currentY)+FontSize+float64(boxDimension))
		ggCtx.DrawString(songCountStr, float64(currentX)+20, float64(currentY)+FontSize*2+float64(boxDimension))

		aSongRS := g143.NewRect(currentX, currentY, boxDimension, boxDimension+50)
		ObjCoords[2000+i+1] = aSongRS

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
	for i := 1; i <= TotalPages(); i++ {
		aStr := fmt.Sprintf("%d", i)
		aStrW, aStrH := ggCtx.MeasureString(aStr)
		aPageBtnW := aStrW + 10
		aPageBtnH := aStrH + 10

		if i == CurrentPage {
			ggCtx.SetHexColor("#C5BF56")
		} else {
			ggCtx.SetHexColor("#633260")
		}
		ggCtx.DrawRectangle(float64(aPageCurrentX), float64(aPageCurrentY), aPageBtnW, aPageBtnH)
		ggCtx.Fill()

		ggCtx.SetHexColor("#fff")
		ggCtx.DrawString(aStr, 5+float64(aPageCurrentX), float64(aPageCurrentY)+FontSize)

		aPageBtnRS := g143.NewRect(aPageCurrentX, aPageCurrentY, int(aPageBtnW), int(aPageBtnH))
		ObjCoords[3000+i] = aPageBtnRS
		newX := aPageCurrentX + int(aPageBtnW) + aPageGutter
		if newX > (wWidth - int(aPageBtnW)) {
			currentY += int(aPageBtnW) + aPageGutter
			aPageCurrentX = 40
		} else {
			aPageCurrentX += int(aPageBtnW) + aPageGutter
		}
	}

	// send the frame to glfw window
	windowRS := g143.Rect{Width: wWidth, Height: wHeight, OriginX: 0, OriginY: 0}
	g143.DrawImage(wWidth, wHeight, ggCtx.Image(), windowRS)
	window.SwapBuffers()

	// save the frame
	currentWindowFrame = ggCtx.Image()
}

func GetDefaultFontPath() string {
	fontPath := filepath.Join(os.TempDir(), "s223_font.ttf")
	os.WriteFile(fontPath, DefaultFont, 0777)
	return fontPath
}

func FirstUIScrollCallback(window *glfw.Window, xoff, yoff float64) {

	if scrollEventCount != 5 {
		scrollEventCount += 1
		return
	}

	scrollEventCount = 0

	if xoff == 0 && yoff == -1 && CurrentPage != TotalPages() {
		ObjCoords = make(map[int]g143.Rect)
		DrawFirstUI(window, CurrentPage+1)
	} else if xoff == 0 && yoff == 1 && CurrentPage != 1 {
		ObjCoords = make(map[int]g143.Rect)
		DrawFirstUI(window, CurrentPage-1)
	}

}
