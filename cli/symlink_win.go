package cli

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"syscall"

	"github.com/charmbracelet/log"
	"golang.org/x/sys/windows"
)

// winElevatedRun elevates a process to run as an administrator on Windows.
func winElevatedRun(name string, arg ...string) (bool, error) {
	log.Debugf("Elevated process: %q", name)
	ok, err := run("cmd", nil, append([]string{"/C", name}, arg...)...)
	if err != nil {
		ok, err = run("elevate.cmd", nil, append([]string{"cmd", "/C", name}, arg...)...)
	}
	return ok, err
}

// run actually constructs the command to be ran from ZVM
func run(name string, dir *string, arg ...string) (bool, error) {
	c := exec.Command(name, arg...)
	if dir != nil {
		c.Dir = *dir
	}
	var stderr bytes.Buffer
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return false, err
	}

	return true, nil
}

func becomeAdmin() error {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		return err
	}
	
	return nil
}

func checkAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")

	return err == nil
}