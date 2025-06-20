package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganyariya/tinyserver/internal/common"
	"github.com/ganyariya/tinyserver/internal/tcp"
	pkgtcp "github.com/ganyariya/tinyserver/pkg/tcp"
)

func main() {
	// Parse command line flags
	var (
		port    = flag.Int("port", pkgtcp.DefaultEchoPort, "Port to listen on")
		host    = flag.String("host", "localhost", "Host to bind to")
		verbose = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Set up logger
	logger := common.NewDefaultLogger()
	if *verbose {
		logger.SetLevel(common.LogLevelDebug)
	}

	// Create server address
	address := fmt.Sprintf("%s:%d", *host, *port)

	// Create TCP server
	server, err := tcp.NewServer("tcp", address)
	if err != nil {
		logger.Error("Failed to create server: %v", err)
		os.Exit(1)
	}

	// Set up echo handler
	server.SetHandler(echoHandler(logger))

	// Start server
	logger.Info("Starting TCP Echo Server on %s", address)
	if err := server.Start(); err != nil {
		logger.Error("Failed to start server: %v", err)
		os.Exit(1)
	}

	logger.Info("TCP Echo Server is running...")
	logger.Info("Press Ctrl+C to stop the server")

	// Set up graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-signalChan

	logger.Info("Shutting down server...")
	if err := server.Stop(); err != nil {
		logger.Error("Error during server shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Server stopped successfully")
}

// echoHandler creates a connection handler that echoes back received data
func echoHandler(logger *common.Logger) pkgtcp.ConnectionHandler {
	return func(conn pkgtcp.Connection) {
		defer conn.Close()

		remoteAddr := conn.RemoteAddr().String()
		logger.Info("New client connected: %s", remoteAddr)

		// Set connection timeout
		if err := conn.SetDeadline(time.Now().Add(5 * time.Minute)); err != nil {
			logger.Warn("Failed to set connection deadline: %v", err)
		}

		buffer := make([]byte, 1024)

		for {
			// Read data from client
			n, err := conn.Read(buffer)
			if err != nil {
				if err.Error() != "EOF" {
					logger.Debug("Read error from %s: %v", remoteAddr, err)
				}
				break
			}

			if n == 0 {
				continue
			}

			receivedData := buffer[:n]
			logger.Debug("Received from %s: %q", remoteAddr, string(receivedData))

			// Echo back the data
			_, err = conn.Write(receivedData)
			if err != nil {
				logger.Debug("Write error to %s: %v", remoteAddr, err)
				break
			}

			logger.Debug("Echoed back to %s: %q", remoteAddr, string(receivedData))

			// Reset deadline for next operation
			if err := conn.SetDeadline(time.Now().Add(5 * time.Minute)); err != nil {
				logger.Warn("Failed to reset connection deadline: %v", err)
			}
		}

		logger.Info("Client disconnected: %s", remoteAddr)
	}
}
