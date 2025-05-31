/*
Copyright © 2025 Kodo Robotics
*/
package core

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func RunShellProvisionOverSSH(username, password string, hostPort int, scriptPath string) error {
	// Load script
	scriptBytes, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to read script: %v", err)
	}
	script := string(scriptBytes)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	addr := fmt.Sprintf("127.0.0.1:%d", hostPort)

	var client *ssh.Client
	err = waitForSSHWithDots(func() error {
		var err error
		client, err = ssh.Dial("tcp", addr, config)
		return err
	}, 30*time.Second)
	if err != nil {
		return fmt.Errorf("ssh connection failed: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Pipe stdout and stderr to terminal
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdin pipe: %v", err)
	}

	fmt.Println("📜 Running provision script...")
	if err := session.Start("bash -s"); err != nil {
		return fmt.Errorf("failed to start shell: %v", err)
	}

	_, err = stdin.Write([]byte(script))
	if err != nil {
		return fmt.Errorf("failed to send script: %v", err)
	}
	stdin.Close()

	err = session.Wait()
	if err != nil {
		return fmt.Errorf("script execution failed: %v", err)
	}

	fmt.Println("✅ Provisioning complete.")
	return nil
}

func waitForSSHWithDots(connectFn func() error, timeout time.Duration) error {
	done := make(chan error)

	// Run connect function in background
	go func() {
		err := connectFn()
		done <- err
	}()

	// Spinner Loop
	fmt.Print("🔐 Connecting to VM via SSH")
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	timeoutTimer := time.NewTicker(timeout)
	defer timeoutTimer.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Print(".")
		case err := <-done:
			fmt.Println()
			return err
		case <-timeoutTimer.C:
			fmt.Println()
			return fmt.Errorf("SSH connection timed out after %s", timeout)
		}
	}
}
