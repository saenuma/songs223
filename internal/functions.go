package internal

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "os error")
	}

	dd := os.Getenv("SNAP_USER_COMMON")

	if strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) || dd == "" {
		dd = filepath.Join(hd, "Songs223")
		os.MkdirAll(dd, 0777)
	}

	return dd, nil
}

func DoesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func UntestedRandomString(length int) string {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ExternalLaunch(p string) {
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/C", "start", p).Run()
	} else if runtime.GOOS == "linux" {
		exec.Command("xdg-open", p).Run()
	}
}

func SecondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	secondsStr := fmt.Sprintf("%d", seconds)
	if seconds < 10 {
		secondsStr = "0" + secondsStr
	}
	str := fmt.Sprintf("%d:%s", minutes, secondsStr)
	return str
}

func getAllL8fFolders() []SongFolder {
	rootPath, _ := GetRootPath()
	ret := make([]SongFolder, 0)

	dirFIs, err := os.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err.Error())
		return ret
	}

	noCoverPath := filepath.Join(os.TempDir(), "no_cover.png")
	os.WriteFile(noCoverPath, NoCover, 0777)

	for _, dirFI := range dirFIs {
		if !dirFI.IsDir() || strings.HasPrefix(dirFI.Name(), ".") {
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

		if l8fCount != 0 {
			ret = append(ret, SongFolder{dirFI.Name(), coverPath, l8fCount})
		}
	}

	return ret
}

func GetFolders(page int) []SongFolder {
	allFolders := getAllL8fFolders()
	beginIndex := (page - 1) * PageSize
	endIndex := beginIndex + PageSize

	var trueRet []SongFolder
	if len(allFolders) <= PageSize {
		trueRet = allFolders
	} else if page == 1 {
		trueRet = allFolders[:PageSize]
	} else if endIndex > len(allFolders) {
		trueRet = allFolders[beginIndex+1:]
	} else {
		trueRet = allFolders[beginIndex+1 : endIndex+1]
	}
	return trueRet
}

func TotalPages() int {
	allFolders := getAllL8fFolders()
	return int(math.Ceil(float64(len(allFolders)) / float64(PageSize)))
}
