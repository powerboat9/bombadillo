// Package telnet provides a function that starts a telnet session in a subprocess.
package telnet

import (
	"fmt"
	"os"
	"os/exec"

	"tildegit.org/sloum/bombadillo/cui"
)

// StartSession starts a telnet session as a subprocess, connecting to the host
// and port specified. Telnet is run interactively as a subprocess until the
// process ends. It returns any errors from the telnet session.
func StartSession(host string, port string) (string, error) {
	c := exec.Command("telnet", host, port)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	// Clear the screen and position the cursor at the top left
	fmt.Print("\033[2J\033[0;0H")
	// Defer reset and reinit of the terminal to prevent any changes from
	// telnet carrying over to the client (or beyond...)
	defer func() {
		cui.Tput("reset")
		cui.InitTerm()
	}()

	err := c.Run()
	if err != nil {
		return "", fmt.Errorf("Telnet error response: %s", err.Error())
	}

	return "Telnet session terminated", nil
}
