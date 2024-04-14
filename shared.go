package main

import (
	"fmt"
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

func externalLaunch(p string) {
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
