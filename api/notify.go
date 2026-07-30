package handler

import (
	"encoding/json"
	"fmt"
	"go_text/pkg/telegram"
	"net/http"
)

type TelegramMessage struct {
}

func NotifyHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var update TelegramMessage

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fmt.Println(update)
	/*
		parse body
		design how we want this to be used.

		#add <hobby>
		#delete <hobby>
		#results weekly total for all hobbies completed

		simple text no # to add completed hobby to db
	*/

	err := telegram.SendTelegramMessage("received")
	if err != nil {
		fmt.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "message sent to telegram"})
}
