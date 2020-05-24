// +build linux

package termios

import "syscall"

const (
	getTermiosIoctl = syscall.TCGETS
	setTermiosIoctl = syscall.TCSETS
)
