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
		val, err := parseValue(reader)
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
	// Get byte length of first statement (command type)
	buf, err := ReadLine(reader)
	if err != nil {
		return -1, err
	}

	if len(buf) < 1 || len(buf) > 2 {
		return -1, errors.New("Invalid command length")
	}

	size, err := strconv.Atoi(string(buf))
	if err != nil {
		return -1, err
	}

	// Read command buffer size
	cmdBuf := make([]byte, size)
	_, err = io.ReadFull(reader, cmdBuf)
	if err != nil {
		return -1, err
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

// Returns the key of the given line and on what index it finished parsing.
func parseKey(reader *bufio.Reader) (string, error) {
	// Get byte length of first statement (command type)
	buf, err := ReadLine(reader)
	if err != nil {
		return "", err
	}

	if len(buf) < 1 {
		return "", errors.New("Invalid key length")
	}

	size, err := strconv.Atoi(string(buf))
	if err != nil {
		return "", err
	}

	// Read command buffer size
	keyBuf := make([]byte, size)
	_, err = io.ReadFull(reader, keyBuf)
	if err != nil {
		return "", err
	}

	keyStr := string(keyBuf)

	return keyStr, nil
}

// Returns the value in bytes.
// Skips whitespaces at the start of given string
func parseValue(reader *bufio.Reader) ([]byte, error) {
	// Get byte length of first statement (command type)
	buf, err := ReadLine(reader)
	if err != nil {
		return nil, err
	}

	if len(buf) < 1 {
		return nil, errors.New("Invalid value length")
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
