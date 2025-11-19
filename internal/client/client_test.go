package client

import "testing"

func TestCommandEncoding(t *testing.T) {
	result := encodeProtocol("GET", "John", []byte("test"))
	expected := "3\r\nGET4\r\nJohn4\r\ntest"

	if string(result) != expected {
		t.Fatalf("expected encoding to be %s, but was: %s", expected, result)
	}
}
