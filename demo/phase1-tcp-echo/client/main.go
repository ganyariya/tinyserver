package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ganyariya/tinyserver/internal/common"
	"github.com/ganyariya/tinyserver/internal/tcp"
	pkgtcp "github.com/ganyariya/tinyserver/pkg/tcp"
)

func main() {
	// Parse command line flags
	var (
		port    = flag.Int("port", pkgtcp.DefaultEchoPort, "Server port to connect to")
		host    = flag.String("host", "localhost", "Server host to connect to")
		verbose = flag.Bool("verbose", false, "Enable verbose logging")
		message = flag.String("message", "", "Single message to send (non-interactive mode)")
	)
	flag.Parse()

	// Set up logger
	logger := common.NewDefaultLogger()
	if *verbose {
		logger.SetLevel(common.LogLevelDebug)
	}

	// Create server address
	address := fmt.Sprintf("%s:%d", *host, *port)

	// Create dialer
	dialer := tcp.NewDialer()

	// Connect to server
	logger.Info("Connecting to TCP Echo Server at %s", address)
	conn, err := dialer.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		logger.Error("Failed to connect to server: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	logger.Info("Connected to server successfully!")

	// Check if we're in single message mode
	if *message != "" {
		sendSingleMessage(conn, *message, logger)
		return
	}

	// Interactive mode
	runInteractiveMode(conn, logger)
}

// sendSingleMessage sends a single message and prints the response
func sendSingleMessage(conn pkgtcp.Connection, message string, logger *common.Logger) {
	// Send message
	_, err := conn.Write([]byte(message))
	if err != nil {
		logger.Error("Failed to send message: %v", err)
		os.Exit(1)
	}

	logger.Debug("Sent: %q", message)

	// Set read timeout
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		logger.Warn("Failed to set read deadline: %v", err)
	}

	// Read response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.Error("Failed to read response: %v", err)
		os.Exit(1)
	}

	response := string(buffer[:n])
	logger.Info("Echo response: %q", response)

	// Verify echo
	if response == message {
		logger.Info("✓ Echo successful!")
	} else {
		logger.Error("✗ Echo mismatch! Expected: %q, Got: %q", message, response)
		os.Exit(1)
	}
}

// runInteractiveMode runs the client in interactive mode
func runInteractiveMode(conn pkgtcp.Connection, logger *common.Logger) {
	logger.Info("Interactive mode started. Type messages to echo. Type 'quit' to exit.")
	fmt.Println()
	fmt.Println("TCP Echo Client - Interactive Mode")
	fmt.Println("=================================")
	fmt.Println("Type your message and press Enter. Type 'quit' to exit.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Prompt for input
		fmt.Print("> ")

		// Read user input
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Check for exit command
		if input == "quit" || input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		// Skip empty input
		if input == "" {
			continue
		}

		// Send message to server
		_, err := conn.Write([]byte(input))
		if err != nil {
			logger.Error("Failed to send message: %v", err)
			break
		}

		logger.Debug("Sent: %q", input)

		// Set read timeout
		if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			logger.Warn("Failed to set read deadline: %v", err)
		}

		// Read echo response
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			logger.Error("Failed to read response: %v", err)
			break
		}

		response := string(buffer[:n])
		fmt.Printf("Echo: %s\n", response)

		// Verify echo in verbose mode
		if logger.GetLevel() == common.LogLevelDebug {
			if response == input {
				logger.Debug("✓ Echo verified")
			} else {
				logger.Debug("✗ Echo mismatch! Expected: %q, Got: %q", input, response)
			}
		}

		fmt.Println()
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		logger.Error("Input error: %v", err)
	}
}
