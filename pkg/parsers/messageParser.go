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
	AddCommand     CommandType = "add"
	DeleteCommand  CommandType = "delete"
	ResulstCommand CommandType = "results"
	LoggedCommand  CommandType = "log_entry"
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
			return ParsedMessage{Type: AddCommand, Value: payload}
		case "#delete":
			return ParsedMessage{Type: DeleteCommand, Value: payload}
		case "#results":
			return ParsedMessage{Type: ResulstCommand}
		}
	}
	return ParsedMessage{Type: LoggedCommand, Value: trimmed}
}
