package parsers

import (
	"strings"
)

type CommandType string

type ParsedMessage struct {
	Type  CommandType
	Value string
}

const (
	addCommand     CommandType = "add"
	deleteCommand  CommandType = "delete"
	resulstCommand CommandType = "results"
	loggedCommand  CommandType = "log_entry"
)

func MessageParser(text string) ParsedMessage {

	trimmed := strings.TrimSpace(text)

	if strings.HasPrefix(trimmed, "#") {
		parts := strings.Split(trimmed, " ")
		cmd := parts[0]

		var payload string

		if len(parts) > 1 {
			payload = parts[1]
		}

		switch cmd {
		case "#add":
			return ParsedMessage{Type: addCommand, Value: payload}
		case "#delete":
			return ParsedMessage{Type: deleteCommand, Value: payload}
		case "#results":
			return ParsedMessage{Type: resulstCommand}
		}
	}
	return ParsedMessage{Type: loggedCommand, Value: trimmed}
}
