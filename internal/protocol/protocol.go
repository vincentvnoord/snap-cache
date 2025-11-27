package protocol

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type CommandType int

const (
	Ping CommandType = iota
	Get
	Set
)

type Command struct {
	CommandType CommandType
	Key         string
	Value       []byte
}

// Reads line until CRLF (\r\n) returns in bytes.
func ReadLine(reader *bufio.Reader) ([]byte, error) {
	var buf []byte

	// Read until \r\n or error
	for {
		byte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if byte == '\r' {
			next, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			// Finished reading line
			if next == '\n' {
				return buf, nil
			}

			buf = append(buf, byte, next)
		}

		buf = append(buf, byte)
	}
}

// Parse reads bytes from the provided bufio.Reader and parses them into a Command.
//
// The reader should contain bytes in the protocol format expected by this server,
// e.g., a command type, key, and optional value, properly encoded with line endings(\r\n).
// Parse will read until a full command is received or return an error if the input
// is malformed or incomplete.
//
// Returns a pointer to a Command struct representing the parsed command,
// or an error if parsing fails.
//
// Example usage:
//
//	reader := bufio.NewReader(conn)
//	cmd, err := Parse(reader)
//	if err != nil {
//		// handle error
//	}
//	fmt.Println("Command type:", cmd.CommandType)
func Parse(reader *bufio.Reader) (*Command, error) {
	// Return err if not valid command
	cmdType, err := parseCommand(reader)
	if err != nil {
		return nil, err
	}

	key, err := parseKey(reader)
	if err != nil {
		return nil, err
	}

	value := []byte{}
	if cmdType == Set {
		val, err := captureBytes(reader, 1, -1)
		value = val
		if err != nil {
			return nil, err
		}
	}

	return &Command{
		CommandType: cmdType,
		Key:         key,
		Value:       value,
	}, nil
}

// Returns the type of command from a reader (first incoming bytes should be the command type).
func parseCommand(reader *bufio.Reader) (CommandType, error) {
	cmdBuf, err := captureBytes(reader, 1, 10)
	if err != nil {
		return 0, err
	}

	cmdStr := string(cmdBuf)

	// Switch case matching input string with possible commands
	commandType := Ping
	switch strings.ToUpper(cmdStr) {
	case "PING":
		commandType = Ping
	case "GET":
		commandType = Get
	case "SET":
		commandType = Set
	default:
		return 0, errors.New("Parsed line is not a valid command")
	}

	return commandType, nil
}

// Returns the key from reader.
func parseKey(reader *bufio.Reader) (string, error) {
	buf, err := captureBytes(reader, 1, -1)
	if err != nil {
		return "", err
	}

	keyStr := string(buf)

	return keyStr, nil
}

// Captures bytes considering captured byte length in first bytes
//
// Example read input: "3\r\nGET"
//
// # This reads 3 bytes, parses into int and takes next set of bytes with length of 3, returning GET
//
// maxLen <= 0 disables the max constraint
func captureBytes(reader *bufio.Reader, minLen int, maxLen int) ([]byte, error) {
	// Read buffer that should be length of next bytes
	buf, err := ReadLine(reader)
	if err != nil {
		return nil, err
	}

	if len(buf) < minLen {
		return nil, errors.New("Byte length too short")
	}

	if maxLen > 0 && len(buf) > maxLen {
		return nil, errors.New("Byte length too long")
	}

	size, err := strconv.Atoi(string(buf))
	if err != nil {
		return nil, err
	}

	// Read command buffer size
	value := make([]byte, size)
	_, err = io.ReadFull(reader, value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
