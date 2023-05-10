package handler

import (
	"encoding/json"
	"fmt"
	m "internbot/internal/model"
	"net/http"
	"strings"
)

type Handler struct{}

var keyboard = m.ReplyMarkup{
	Keyboard: [][]m.KeyboardButton{
		{
			{
				Text: "Button 1",
			},
			{
				Text: "Button 2",
			},
			{
				Text: "Button 3",
			},
			{
				Text: "Button 4",
			},
		},
	},
	ResizeKeyboard:  true,
	OneTimeKeyboard: true,
	IsPersistent:    true,
}

func (h *Handler) Control(response *http.Response, lastId *int64) (message *m.Message, err error) {

	if response.StatusCode != http.StatusOK {
		err := fmt.Errorf("wtf")
		return nil, err
	}

	data := m.Data{}

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}

	for _, update := range data.Result {

		if *lastId < update.UpdateID {
			switch {
			case update.Message.Text == "/start":
				message = &m.Message{
					ChatId:      update.Message.From.Id,
					Text:        update.Message.From.Username + ", welcome! Please, use '/menu'",
					ReplyMarkup: &keyboard,
				}
			case strings.HasPrefix(update.Message.Text, "Button"):
				message = &m.Message{
					ChatId:      update.Message.From.Id,
					Text:        "You pressed: " + update.Message.Text,
					ReplyMarkup: &keyboard,
				}
			}
			*lastId = update.UpdateID
		}
	}
	return message, nil
}
