package protocol

import "testing"

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
