package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"go_text/internal/db" // Adjust to match your go.mod package path
	"go_text/pkg/parsers"
	"go_text/pkg/telegram"
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
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var update TelegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if update.Message != nil && update.Message.Text != "" {
		msg := update.Message.Text
		chatID := update.Message.Chat.ID

		// 1. Set up context with timeout for database queries
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// 2. Fetch / Reuse database pool
		pool, err := db.GetPool(ctx)
		if err != nil {
			log.Printf("[DB ERROR] Failed to connect: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// 3. Parse user message
		parsedMsg := parsers.MessageParser(msg)
		var responseText string

		// 4. Handle commands against database
		switch parsedMsg.Type {

		case parsers.AddCommand:
			if parsedMsg.Value == "" {
				responseText = "⚠️ Please specify a hobby to add. Example: `#add reading`"
			} else {
				created, err := db.AddHobby(ctx, pool, chatID, parsedMsg.Value)
				if err != nil {
					log.Printf("[DB ERROR] AddHobby failed: %v", err)
					responseText = "❌ Failed to save hobby."
				} else if created {
					responseText = fmt.Sprintf("✅ Added '%s' to your tracked hobbies!", parsedMsg.Value)
				} else {
					responseText = fmt.Sprintf("⚠️ '%s' is already in your hobbies list.", parsedMsg.Value)
				}
			}

		case parsers.ResulstCommand:
			results, err := db.GetWeeklyResults(ctx, pool, chatID)
			if err != nil {
				log.Printf("[DB ERROR] GetWeeklyResults failed: %v", err)
				responseText = "❌ Failed to fetch weekly results."
			} else if len(results) == 0 {
				responseText = "📊 No hobbies tracked yet! Use `#add <hobby>` to start."
			} else {
				var sb strings.Builder
				sb.WriteString("📊 Weekly Summary (Last 7 Days):\n\n")
				for _, res := range results {
					sb.WriteString(fmt.Sprintf("• %s: %d completed\n", res.Name, res.Count))
				}
				responseText = sb.String()
			}

		case parsers.LoggedCommand:

			inDb := db.CheckHobby(ctx, pool, parsedMsg.Value)
			if inDb == false {
				err := db.LogHobbyCompletion(ctx, pool, chatID, parsedMsg.Value)
				if err != nil {
					log.Printf("[DB ERROR] LogHobbyCompletion failed: %v", err)
					responseText = "❌ Failed to log entry."
				} else {
					responseText = fmt.Sprintf("🎉 Logged '%s'!", parsedMsg.Value)
				}

			} else {
				responseText = "Hobby has not been added to your tracked hobbies!"

			}
		}

		// 5. Reply to the user on Telegram
		if responseText != "" {
			if err := telegram.SendTelegramMessage(responseText); err != nil {
				log.Printf("[TELEGRAM ERROR] Failed to send reply: %v", err)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
