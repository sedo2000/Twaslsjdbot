package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type Update struct {
	Message *Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int64 `json:"id"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var update Update
	json.NewDecoder(r.Body).Decode(&update)

	if update.Message == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	botToken := os.Getenv("BOT_TOKEN")
	adminID := os.Getenv("ADMIN_ID")

	text := update.Message.Text
	userID := update.Message.Chat.ID

	// رسالة ترسل للأدمن
	if adminID != "" {
		sendMessage(botToken, adminID,
			"📩 رسالة جديدة:\n"+text+"\n\n👤 UserID: "+intToString(userID))
	}

	// رد للمستخدم
	sendMessage(botToken, intToString(userID),
		"تم استلام رسالتك ✅\nسيتم الرد عليك قريبًا.")

	w.WriteHeader(http.StatusOK)
}

func sendMessage(token, chatID, text string) {
	url := "https://api.telegram.org/bot" + token + "/sendMessage"

	payload := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}

	body, _ := json.Marshal(payload)
	http.Post(url, "application/json", bytes.NewBuffer(body))
}

func intToString(id int64) string {
	return json.Number(id).String()
}
