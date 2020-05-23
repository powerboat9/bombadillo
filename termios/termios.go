package termios

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

var fd = os.Stdin.Fd()

func ioctl(fd, request, argp uintptr) error {
	if _, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, request, argp); e != 0 {
		return e
	}
	return nil
}

func GetWindowSize() (int, int) {
	var value winsize
	ioctl(fd, syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&value)))
	return int(value.Col), int(value.Row)
}

func getTermios() syscall.Termios {
	var value syscall.Termios
	err := ioctl(fd, getTermiosIoctl, uintptr(unsafe.Pointer(&value)))
	if err != nil {
		panic(err)
	}
	return value
}

func setTermios(termios syscall.Termios) {
	err := ioctl(fd, setTermiosIoctl, uintptr(unsafe.Pointer(&termios)))
	if err != nil {
		panic(err)
	}
	runtime.KeepAlive(termios)
}

func SetCharMode() {
	t := getTermios()
	t.Lflag = t.Lflag ^ syscall.ICANON
	t.Lflag = t.Lflag ^ syscall.ECHO
	setTermios(t)
}

func SetLineMode() {
	var t = getTermios()
	t.Lflag = t.Lflag | (syscall.ICANON | syscall.ECHO)
	setTermios(t)
}
