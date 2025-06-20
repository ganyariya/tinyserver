package tcp

import (
	"net"
	"testing"
	"time"
)

func TestNewConnection(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	conn := NewConnection(server)
	if conn == nil {
		t.Fatal("NewConnection returned nil")
	}

	// Test that the connection implements the interface properly
	if conn.LocalAddr() != server.LocalAddr() {
		t.Error("LocalAddr mismatch")
	}

	if conn.RemoteAddr() != server.RemoteAddr() {
		t.Error("RemoteAddr mismatch")
	}
}

func TestConnectionReadWrite(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewConnection(server)
	clientConn := NewConnection(client)

	// Test data
	testData := []byte("Hello, TinyServer!")

	// Write from client
	go func() {
		n, err := clientConn.Write(testData)
		if err != nil {
			t.Errorf("Write failed: %v", err)
		}
		if n != len(testData) {
			t.Errorf("Write length mismatch: expected %d, got %d", len(testData), n)
		}
	}()

	// Read from server
	buffer := make([]byte, len(testData))
	n, err := serverConn.Read(buffer)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if n != len(testData) {
		t.Errorf("Read length mismatch: expected %d, got %d", len(testData), n)
	}

	if string(buffer) != string(testData) {
		t.Errorf("Data mismatch: expected %s, got %s", string(testData), string(buffer))
	}
}

func TestConnectionClose(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer client.Close()

	conn := NewConnection(server)

	// Close the connection
	err := conn.Close()
	if err != nil {
		t.Fatalf("Close failed: %v", err)
	}

	// Try to write after close - should fail
	_, err = conn.Write([]byte("test"))
	if err == nil {
		t.Error("Write should fail after close")
	}

	// Try to read after close - should fail
	buffer := make([]byte, 10)
	_, err = conn.Read(buffer)
	if err == nil {
		t.Error("Read should fail after close")
	}
}

func TestConnectionDeadlines(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	conn := NewConnection(server)

	// Test SetDeadline
	deadline := time.Now().Add(time.Second)
	err := conn.SetDeadline(deadline)
	if err != nil {
		t.Errorf("SetDeadline failed: %v", err)
	}

	// Test SetReadDeadline
	err = conn.SetReadDeadline(deadline)
	if err != nil {
		t.Errorf("SetReadDeadline failed: %v", err)
	}

	// Test SetWriteDeadline
	err = conn.SetWriteDeadline(deadline)
	if err != nil {
		t.Errorf("SetWriteDeadline failed: %v", err)
	}
}

func TestBufferedConnection(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewBufferedConnection(server)
	clientConn := NewBufferedConnection(client)

	// Test WriteLine and ReadLine
	testLine := "Hello, TinyServer!"

	// Write a line from client
	go func() {
		err := clientConn.WriteLine([]byte(testLine))
		if err != nil {
			t.Errorf("WriteLine failed: %v", err)
		}
	}()

	// Read the line from server
	line, err := serverConn.ReadLine()
	if err != nil {
		t.Fatalf("ReadLine failed: %v", err)
	}

	if string(line) != testLine {
		t.Errorf("Line mismatch: expected %s, got %s", testLine, string(line))
	}
}

func TestBufferedConnectionFlush(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	conn := NewBufferedConnection(server)

	// Write some data
	writer := conn.BufferedWriter()
	_, err := writer.Write([]byte("test"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Flush the data
	err = conn.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}
}

func TestMessageConnection(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewConnection(server)
	clientConn := NewConnection(client)

	serverMsgConn := NewMessageConnection(serverConn)
	clientMsgConn := NewMessageConnection(clientConn)

	// Test message sending
	testMessage := []byte("Hello, TinyServer!")

	// Send message from client
	go func() {
		err := clientMsgConn.WriteMessage(testMessage)
		if err != nil {
			t.Errorf("WriteMessage failed: %v", err)
		}
	}()

	// Receive message on server
	message, err := serverMsgConn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}

	if string(message) != string(testMessage) {
		t.Errorf("Message mismatch: expected %s, got %s", string(testMessage), string(message))
	}
}

func TestMessageConnectionWithCustomDelimiter(t *testing.T) {
	// Create a test connection using a pipe
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewConnection(server)
	clientConn := NewConnection(client)

	serverMsgConn := NewMessageConnection(serverConn)
	clientMsgConn := NewMessageConnection(clientConn)

	// Set custom delimiter
	delimiter := []byte("||")
	serverMsgConn.SetMessageDelimiter(delimiter)
	clientMsgConn.SetMessageDelimiter(delimiter)

	// Test message sending with custom delimiter
	testMessage := []byte("Hello, TinyServer!")

	// Send message from client
	go func() {
		err := clientMsgConn.WriteMessage(testMessage)
		if err != nil {
			t.Errorf("WriteMessage failed: %v", err)
		}
	}()

	// Receive message on server
	message, err := serverMsgConn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage failed: %v", err)
	}

	if string(message) != string(testMessage) {
		t.Errorf("Message mismatch: expected %s, got %s", string(testMessage), string(message))
	}
}

func TestFindDelimiter(t *testing.T) {
	tests := []struct {
		buffer    []byte
		delimiter []byte
		expected  int
	}{
		{[]byte("hello\nworld"), []byte("\n"), 5},
		{[]byte("hello||world"), []byte("||"), 5},
		{[]byte("helloworld"), []byte("\n"), -1},
		{[]byte(""), []byte("\n"), -1},
		{[]byte("hello"), []byte(""), -1},
	}

	for i, test := range tests {
		result := findDelimiter(test.buffer, test.delimiter)
		if result != test.expected {
			t.Errorf("Test %d: expected %d, got %d", i, test.expected, result)
		}
	}
}

func TestMatchDelimiter(t *testing.T) {
	tests := []struct {
		buffer    []byte
		delimiter []byte
		expected  bool
	}{
		{[]byte("hello\nworld"), []byte("hello"), true},
		{[]byte("hello\nworld"), []byte("world"), false},
		{[]byte("hi"), []byte("hello"), false},
		{[]byte(""), []byte("hello"), false},
		{[]byte("hello"), []byte(""), true},
	}

	for i, test := range tests {
		result := matchDelimiter(test.buffer, test.delimiter)
		if result != test.expected {
			t.Errorf("Test %d: expected %t, got %t", i, test.expected, result)
		}
	}
}

// Benchmark tests
func BenchmarkConnectionReadWrite(b *testing.B) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewConnection(server)
	clientConn := NewConnection(client)

	data := make([]byte, 1024)
	buffer := make([]byte, 1024)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		go func() {
			clientConn.Write(data)
		}()
		serverConn.Read(buffer)
	}
}

func BenchmarkBufferedConnectionReadWrite(b *testing.B) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	serverConn := NewBufferedConnection(server)
	clientConn := NewBufferedConnection(client)

	data := []byte("test line")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		go func() {
			clientConn.WriteLine(data)
		}()
		serverConn.ReadLine()
	}
}
