package protocol

import (
	"bytes"
	"testing"
)

func TestParseGetCommandSuccess(t *testing.T) {
	cmdType, _, err := parseCommand("get")
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	if cmdType != Get {
		t.Fatalf("expected cmd type to be 0, got: %d", cmdType)
	}

	cmdType, _, err = parseCommand("get   ")
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	if cmdType != Get {
		t.Fatalf("expected cmd type to be 0, got: %d", cmdType)
	}
}

func TestParseSetCommandSuccess(t *testing.T) {
	cmdType, _, err := parseCommand("set")
	if err != nil {
		t.Fatalf("expected no error, got: %s", err.Error())
	}

	if cmdType != Set {
		t.Fatalf("expected cmd type to be 0, got: %d", cmdType)
	}
}

func TestParseCommandError(t *testing.T) {
	cmdType, _, err := parseCommand(" set")
	if err == nil {
		t.Fatalf("expected parse error, got cmd: %d", cmdType)
	}
}

func TestParseSetCommandReturnsKeyValue(t *testing.T) {
	parsed := Parse("set alice 19")

	if parsed.Key != "alice" {
		t.Fatalf("expected key to be alice, but was: %s", parsed.Key)
	}

	if !bytes.Equal(parsed.Value, []byte("19")) {
		t.Fatalf("expected key to be alice, but was: %s", parsed.Key)
	}
}

func TestParseSetCommandVariousInputs(t *testing.T) {
	tests := []struct {
		input       string
		expectedKey string
		expectedVal []byte
	}{
		{
			input:       "set bob 12345678901234567890",
			expectedKey: "bob",
			expectedVal: []byte("12345678901234567890"),
		},
		{
			input:       "SET user some long value with spaces inside",
			expectedKey: "user",
			expectedVal: []byte("some long value with spaces inside"),
		},
		{
			input:       "set image \x00\xFF\x10\x20", // binary bytes
			expectedKey: "image",
			expectedVal: []byte("\x00\xFF\x10\x20"),
		},
		{
			input:       "   set   spacedkey    spaced   value   ",
			expectedKey: "spacedkey",
			expectedVal: []byte("spaced   value   "),
		},
		{
			input:       "set empty ",
			expectedKey: "empty",
			expectedVal: []byte(""), // value is empty string
		},
	}

	for _, tc := range tests {
		parsed := Parse(tc.input)

		if parsed.Key != tc.expectedKey {
			t.Fatalf("for input %q expected key %q but got %q",
				tc.input, tc.expectedKey, parsed.Key)
		}

		if !bytes.Equal(parsed.Value, tc.expectedVal) {
			t.Fatalf("for input %q expected value %v but got %v",
				tc.input, tc.expectedVal, parsed.Value)
		}
	}
}
