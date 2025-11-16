package protocol

import (
	"errors"
	"strings"
)

type CommandType int

const (
	Get CommandType = iota
	Set
)

type Command struct {
	commandType CommandType
	key         string
}

func Parse(line string) Command {
	command, next, err := parseCommand(line)
	if err != nil {
		// Return err
	}

	// Parse key
	key, next := parseKey(line, next)

	// If SET command parse value

	return Command{
		commandType: command,
		key:         key,
	}
}

// Returns the type of command and on what index it finished parsing.
func parseCommand(line string) (CommandType, int, error) {
	var builder strings.Builder
	index := 0

	for i := 0; i < len(line); i++ {
		char := line[i]
		if char == ' ' {
			index = i
			break
		}

		builder.WriteByte(char)
	}

	commandType := Get

	switch builder.String() {
	case "get":
		commandType = Get
	case "GET":
		commandType = Get
	case "set":
		commandType = Set
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

	for i := from; i < len(line); i++ {
		char := line[i]
		if char == ' ' {
			// Pass it to the next char
			index = i + 1
			break
		}

		builder.WriteByte(char)
	}
	return builder.String(), index
}
