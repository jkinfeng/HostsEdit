package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
	"fmt"
	"golang.org/x/sys/windows"
)

const (
	MB_OK                = 0x00000000
	MB_ICONERROR         = 0x00000010
)

func getAppAndFile() ([]string, error) {
	systemRootDir := os.Getenv("SystemRoot")
	if systemRootDir == "" {
		return nil, fmt.Errorf("fatal error: unable to retrieve system environment variables")
	}

	notepadPath := filepath.Join(systemRootDir, "System32", "notepad.exe")
	hostsPath := filepath.Join(systemRootDir, "System32", "drivers", "etc", "hosts")

	if _, err := os.Stat(notepadPath); err != nil {
		return nil, err
	}
	
	if _, err := os.Stat(hostsPath); err != nil {
		return nil, err
	}

	return []string{notepadPath, hostsPath}, nil
}

func setFileAttributes(filename string, attrs uint32) error {
	pathPtr, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	
	err = windows.SetFileAttributes(pathPtr, attrs)
	if err != nil {
		return err
	}

	return nil
}

func MessageBox(title, text string, flags uint32) int {
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0
	}
	textPtr, err := windows.UTF16PtrFromString(text)
	if err != nil {
		return 0
	}

	ret, _, _ := windows.NewLazySystemDLL("user32.dll").NewProc("MessageBoxW").Call(
		0,
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(unsafe.Pointer(titlePtr)),
		uintptr(flags),
	)

	return int(ret)
}

func main() {
	paths, err := getAppAndFile()
	if err != nil {
		os.Exit(1)
	}

	err = setFileAttributes(paths[1], windows.FILE_ATTRIBUTE_NORMAL)
	if err != nil {
		MessageBox("Error", "Please run as administrator", MB_OK|MB_ICONERROR)
		os.Exit(1)
	}

	cmd := exec.Command(paths[0], paths[1])
	_ = cmd.Run()
	_ = setFileAttributes(paths[1], windows.FILE_ATTRIBUTE_READONLY | windows.FILE_ATTRIBUTE_SYSTEM | windows.FILE_ATTRIBUTE_HIDDEN)

	os.Exit(0)
}
