// +build !linux

package termios

import "syscall"

const (
	getTermiosIoctl = syscall.TIOCGETA
	setTermiosIoctl = syscall.TIOCSETAF
)
