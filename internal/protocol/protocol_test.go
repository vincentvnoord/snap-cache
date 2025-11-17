package protocol

import (
	"bufio"
	"bytes"
	"strconv"
	"testing"
)

func buildBufReader(input string) *bufio.Reader {
	buf := bytes.NewBufferString(input)
	return bufio.NewReader(buf)
}

// helper: wrap input as a length-prefixed command
func buildLengthPrefixed(input string) *bufio.Reader {
	cmd := []byte(input)
	sizeStr := strconv.Itoa(len(cmd))
	return buildBufReader(sizeStr + "\r\n" + input)
}

func TestParseGetCommandSuccess(t *testing.T) {
	reader := buildLengthPrefixed("GET")

	cmdType, err := parseCommand(reader)
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}
	if cmdType != Get {
		t.Fatalf("expected cmd type Get, got: %d", cmdType)
	}
}

func TestParseSetCommandSuccess(t *testing.T) {
	reader := buildLengthPrefixed("SET")

	cmdType, err := parseCommand(reader)
	if err != nil {
		t.Fatalf("expected no error, got: %s", err)
	}
	if cmdType != Set {
		t.Fatalf("expected cmd type Set, got: %d", cmdType)
	}
}

func TestParseCommandError(t *testing.T) {
	reader := buildLengthPrefixed(" UNKNOWN")

	_, err := parseCommand(reader)
	if err == nil {
		t.Fatalf("expected parse error for unknown command")
	}
}

func TestParseSetCommandReturnsKeyValue(t *testing.T) {
	input := "3\r\nSET5\r\nalice2\r\n19"
	reader := buildBufReader(input)

	parsed, err := Parse(reader)
	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	if parsed.Key != "alice" {
		t.Fatalf("expected key 'alice', got: %s", parsed.Key)
	}
	if !bytes.Equal(parsed.Value, []byte("19")) {
		t.Fatalf("expected value '19', got: %s", parsed.Value)
	}
}

func TestParseSetCommandVariousInputs(t *testing.T) {
	tests := []struct {
		input       string
		expectedKey string
		expectedVal []byte
	}{
		{"3\r\nSET3\r\nbob20\r\n12345678901234567890", "bob", []byte("12345678901234567890")},
		{"3\r\nSET4\r\nuser34\r\nsome long value with spaces inside", "user", []byte("some long value with spaces inside")},
		{"3\r\nSET5\r\nimage4\r\n\x00\xFF\x10\x20", "image", []byte("\x00\xFF\x10\x20")},
		{"3\r\nSET13\r\n    spacedkey21\r\n    spaced   value   ", "    spacedkey", []byte("    spaced   value   ")},
		{"3\r\nSET5\r\nempty0\r\n", "empty", []byte("")},
	}

	for _, tc := range tests {
		reader := buildBufReader(tc.input)
		parsed, err := Parse(reader)
		if err != nil {
			t.Fatalf("input %q: expected no error, got: %s", tc.input, err)
		}

		if parsed.Key != tc.expectedKey {
			t.Fatalf("input %q: expected key %q, got %q", tc.input, tc.expectedKey, parsed.Key)
		}
		if !bytes.Equal(parsed.Value, tc.expectedVal) {
			t.Fatalf("input %q: expected value %v, got %v", tc.input, tc.expectedVal, parsed.Value)
		}
	}
}
