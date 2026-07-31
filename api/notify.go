package handler

import (
	"encoding/json"
	"fmt"
	"go_text/pkg/parsers"
	"go_text/pkg/telegram"
	"net/http"
)

type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message  *struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}

func NotifyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var update TelegramUpdate

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fmt.Printf("%+v\n", r.Body)
	/*
		parse body
		design how we want this to be used.

		#add <hobby>
		#delete <hobby>
		#results weekly total for all hobbies completed

		simple text no # to add completed hobby to db
	*/

	if update.Message != nil && update.Message.Text != "" {
		msg := update.Message.Text
		chatID := update.Message.Chat.ID
		fmt.Printf("%+v\n", chatID)

		parsedMsg := parsers.MessageParser(msg)

		switch parsedMsg.Type {

		case parsers.AddCommand:
			err := telegram.SendTelegramMessage(parsedMsg.Value)
			if err != nil {
				return
			}
		case parsers.DeleteCommand:
			err := telegram.SendTelegramMessage(parsedMsg.Value)
			if err != nil {
				return
			}

		case parsers.ResulstCommand:
			err := telegram.SendTelegramMessage(parsedMsg.Value)
			if err != nil {
				return
			}
		case parsers.LoggedCommand:
			err := telegram.SendTelegramMessage(parsedMsg.Value)
			if err != nil {
				return
			}

		}

	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "message sent to telegram"})
}
