package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusOK)
		return
	}

	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if update.Message == nil || update.Message.Text == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	botToken := os.Getenv("BOT_TOKEN")
	adminID := os.Getenv("ADMIN_ID")

	userID := update.Message.Chat.ID
	userText := update.Message.Text

	// إرسال رسالة للأدمن
	if botToken != "" && adminID != "" {
		sendMessage(
			botToken,
			adminID,
			"📩 رسالة جديدة:\n\n"+userText+"\n\n👤 User ID: "+strconv.FormatInt(userID, 10),
		)
	}

	// رد للمستخدم
	sendMessage(
		botToken,
		strconv.FormatInt(userID, 10),
		"✅ تم استلام رسالتك.\nسيتم الرد عليك قريبًا.",
	)

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
