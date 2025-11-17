package protocol

import (
	"errors"
	"strings"
	"unicode"
)

type CommandType int

const (
	Get CommandType = iota
	Set
)

type Command struct {
	CommandType CommandType
	Key         string
	Value       []byte
}

// Parse interprets a single protocol line such as "GET x" or "SET x 10".
// It returns a Command struct or a parse error if the input is malformed.
func Parse(line string) (*Command, error) {
	command, next, err := parseCommand(line)
	if err != nil {
		// Return err
		return nil, err
	}

	// Parse key
	key, next := parseKey(line, next)

	cmd := Command{
		CommandType: command,
		Key:         key,
	}

	// If SET command parse value
	if command == Set {
		value := parseValue(line, next)
		cmd.Value = value
	}

	return &cmd, nil
}

// Returns the type of command and on what index it finished parsing.
func parseCommand(line string) (CommandType, int, error) {
	var builder strings.Builder
	index := 0

	// Write byte into string builder until whitespace
	for i := 0; i < len(line); i++ {
		index = i
		char := line[i]
		if char == ' ' {
			index = i + 1
			break
		}

		builder.WriteByte(char)
	}

	// Switch case matching input string with possible commands
	commandType := Get
	switch strings.ToUpper(builder.String()) {
	case "GET":
		commandType = Get
	case "SET":
		commandType = Set
	default:
		return 0, index, errors.New("Parsed line is not a valid command")
	}

	return commandType, index, nil
}

// Returns the key of the given line and on what index it finished parsing.
func parseKey(line string, from int) (string, int) {
	var builder strings.Builder
	index := 0

	for from < len(line) && unicode.IsSpace(rune(line[from])) {
		from++
	}

	for i := from; i < len(line); i++ {
		char := line[i]
		if unicode.IsSpace(rune(char)) {
			// Pass it to the next char
			index = i + 1
			break
		}

		builder.WriteByte(char)
	}

	return builder.String(), index
}

// Returns the value in bytes.
// Skips whitespaces at the start of given string
func parseValue(line string, from int) []byte {
	for from < len(line) && line[from] == ' ' {
		from++
	}

	return []byte(line[from:])
}
