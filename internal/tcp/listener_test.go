package tcp

import (
	"net"
	"sync"
	"testing"
	"time"

	pkgtcp "github.com/ganyariya/tinyserver/pkg/tcp"
)

func TestNewListener(t *testing.T) {
	// Get a free port for testing
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Create TinyServer listener
	tcpListener, err := NewListener("tcp", net.JoinHostPort("localhost", string(rune(port))))
	if err != nil {
		t.Fatalf("NewListener failed: %v", err)
	}
	defer tcpListener.Close()

	if tcpListener == nil {
		t.Fatal("NewListener returned nil")
	}

	if tcpListener.Addr() == nil {
		t.Fatal("Listener address is nil")
	}
}

func TestListenerAccept(t *testing.T) {
	// Get a free port for testing
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	address := net.JoinHostPort("localhost", string(rune(port)))

	// Create TinyServer listener
	tcpListener, err := NewListener("tcp", address)
	if err != nil {
		t.Fatalf("NewListener failed: %v", err)
	}
	defer tcpListener.Close()

	// Connect from a client
	var serverConn pkgtcp.Connection
	var clientConn net.Conn
	var acceptErr error

	done := make(chan struct{})

	// Accept connection in goroutine
	go func() {
		serverConn, acceptErr = tcpListener.Accept()
		close(done)
	}()

	// Give the accept goroutine time to start
	time.Sleep(10 * time.Millisecond)

	// Connect as client
	clientConn, err = net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Client dial failed: %v", err)
	}
	defer clientConn.Close()

	// Wait for accept to complete
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Accept timeout")
	}

	if acceptErr != nil {
		t.Fatalf("Accept failed: %v", acceptErr)
	}

	if serverConn == nil {
		t.Fatal("Accept returned nil connection")
	}

	defer serverConn.Close()

	// Test basic communication
	testData := []byte("Hello, TinyServer!")

	// Send from client
	_, err = clientConn.Write(testData)
	if err != nil {
		t.Fatalf("Client write failed: %v", err)
	}

	// Receive on server
	buffer := make([]byte, len(testData))
	n, err := serverConn.Read(buffer)
	if err != nil {
		t.Fatalf("Server read failed: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Read length mismatch: expected %d, got %d", len(testData), n)
	}

	if string(buffer) != string(testData) {
		t.Errorf("Data mismatch: expected %s, got %s", string(testData), string(buffer))
	}
}

func TestListenerClose(t *testing.T) {
	// Get a free port for testing
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Create TinyServer listener
	tcpListener, err := NewListener("tcp", net.JoinHostPort("localhost", string(rune(port))))
	if err != nil {
		t.Fatalf("NewListener failed: %v", err)
	}

	// Close the listener
	err = tcpListener.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Try to accept after close - should fail
	_, err = tcpListener.Accept()
	if err == nil {
		t.Error("Accept should fail after close")
	}
}

func TestConnectionFactory(t *testing.T) {
	factory := NewConnectionFactory()
	if factory == nil {
		t.Fatal("NewConnectionFactory returned nil")
	}

	// Test CreateDialer
	dialer := factory.CreateDialer()
	if dialer == nil {
		t.Error("CreateDialer returned nil")
	}

	// Test WrapConnection
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	conn := factory.WrapConnection(server)
	if conn == nil {
		t.Error("WrapConnection returned nil")
	}
}

func TestDialer(t *testing.T) {
	// Create a test server
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	address := listener.Addr().String()

	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	// Test dialer
	dialer := NewDialer()
	if dialer == nil {
		t.Fatal("NewDialer returned nil")
	}

	// Test Dial
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Fatal("Dial returned nil connection")
	}
}

func TestDialerWithTimeout(t *testing.T) {
	// Create a test server
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	address := listener.Addr().String()

	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	// Test dialer with timeout
	dialer := NewDialer()

	// Test DialTimeout
	conn, err := dialer.DialTimeout("tcp", address, time.Second)
	if err != nil {
		t.Fatalf("DialTimeout failed: %v", err)
	}
	defer conn.Close()

	if conn == nil {
		t.Fatal("DialTimeout returned nil connection")
	}
}

func TestTCPServer(t *testing.T) {
	// Get a free port for testing
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	address := net.JoinHostPort("localhost", string(rune(port)))

	// Create TCP server
	server, err := NewServer("tcp", address)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// Test initial state
	if server.IsRunning() {
		t.Error("Server should not be running initially")
	}

	// Test starting without handler - should fail
	err = server.Start()
	if err == nil {
		t.Error("Start should fail without handler")
	}

	// Set a simple echo handler
	var handlerCalled bool
	var mu sync.Mutex

	server.SetHandler(func(conn pkgtcp.Connection) {
		mu.Lock()
		handlerCalled = true
		mu.Unlock()

		// Simple echo
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		conn.Write(buffer[:n])
	})

	// Start the server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Check that server is running
	if !server.IsRunning() {
		t.Error("Server should be running")
	}

	// Give server time to start
	time.Sleep(10 * time.Millisecond)

	// Connect and test
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Client dial failed: %v", err)
	}
	defer conn.Close()

	// Send test data
	testData := []byte("Hello, TinyServer!")
	_, err = conn.Write(testData)
	if err != nil {
		t.Fatalf("Client write failed: %v", err)
	}

	// Read echo response
	buffer := make([]byte, len(testData))
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Client read failed: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Read length mismatch: expected %d, got %d", len(testData), n)
	}

	if string(buffer) != string(testData) {
		t.Errorf("Echo mismatch: expected %s, got %s", string(testData), string(buffer))
	}

	// Give handler time to be called
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if !handlerCalled {
		t.Error("Handler was not called")
	}
	mu.Unlock()

	// Stop the server
	err = server.Stop()
	if err != nil {
		t.Fatalf("Stop failed: %v", err)
	}

	// Check that server is not running
	if server.IsRunning() {
		t.Error("Server should not be running after stop")
	}
}

func TestServerMultipleConnections(t *testing.T) {
	// Get a free port for testing
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to get free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	address := net.JoinHostPort("localhost", string(rune(port)))

	// Create TCP server
	server, err := NewServer("tcp", address)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}
	defer server.Stop()

	// Set handler that counts connections
	var connectionCount int32
	var mu sync.Mutex

	server.SetHandler(func(conn pkgtcp.Connection) {
		mu.Lock()
		connectionCount++
		mu.Unlock()

		// Keep connection open for a bit
		time.Sleep(50 * time.Millisecond)
	})

	// Start the server
	err = server.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Give server time to start
	time.Sleep(10 * time.Millisecond)

	// Create multiple connections
	numConnections := 5
	var wg sync.WaitGroup

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", address)
			if err != nil {
				t.Errorf("Client dial failed: %v", err)
				return
			}
			defer conn.Close()

			// Keep connection open briefly
			time.Sleep(10 * time.Millisecond)
		}()
	}

	wg.Wait()

	// Give handlers time to complete
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if connectionCount != int32(numConnections) {
		t.Errorf("Expected %d connections, got %d", numConnections, connectionCount)
	}
	mu.Unlock()
}

// Benchmark tests
func BenchmarkListenerAccept(b *testing.B) {
	// Create listener
	listener, err := NewListener("tcp", ":0")
	if err != nil {
		b.Fatalf("NewListener failed: %v", err)
	}
	defer listener.Close()

	address := listener.Addr().String()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Connect as client
		go func() {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return
			}
			conn.Close()
		}()

		// Accept on server
		conn, err := listener.Accept()
		if err != nil {
			b.Errorf("Accept failed: %v", err)
			continue
		}
		conn.Close()
	}
}

func BenchmarkDialer(b *testing.B) {
	// Create test server
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		b.Fatalf("Failed to create test server: %v", err)
	}
	defer listener.Close()

	address := listener.Addr().String()

	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	dialer := NewDialer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		conn, err := dialer.Dial("tcp", address)
		if err != nil {
			b.Errorf("Dial failed: %v", err)
			continue
		}
		conn.Close()
	}
}
