package ports

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

func CheckPort(port int) error {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("port %d is already in use", port)
	}

	// Port is free, close the listener
	ln.Close()
	fmt.Printf("Port %d is not in use.\n", port)
	return nil
}
func CheckAndKillPort(port int) error {
	address := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		// Port is already in use, try to find the process using it and kill it
		fmt.Printf("Port %d is already in use, attempting to close it...\n", port)
		return KillPortProcess(port)
	}

	// Port is free, close the listener
	ln.Close()
	fmt.Printf("Port %d is not in use.\n", port)
	return nil
}

// KillPortProcess kills the process using the specified port
func KillPortProcess(port int) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", fmt.Sprintf("netstat -ano | findstr :%d", port))
	} else {
		cmd = exec.Command("sh", "-c", fmt.Sprintf("lsof -ti :%d", port))
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to find process using port %d: %v", port, err)
	}

	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return fmt.Errorf("no process found using port %d", port)
	}

	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/F", "/PID", pidStr)
	} else {
		cmd = exec.Command("kill", "-9", pidStr)
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill process %s: %v", pidStr, err)
	}

	fmt.Printf("Successfully killed process %s using port %d.\n", pidStr, port)
	return nil
}
