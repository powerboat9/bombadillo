package telnet

import (
	"fmt"
	"os"
	"os/exec"
)

func StartSession(host string, port string) (string, error) {
	// Case for telnet links
	c := exec.Command("telnet", host, port)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	// Clear the screen and position the cursor at the top left
	fmt.Print("\033[2J\033[0;0H")
	err := c.Run()
	if err != nil {
		return "", fmt.Errorf("Telnet error response: %s", err.Error())
	}

	return "Telnet session terminated", nil
}

