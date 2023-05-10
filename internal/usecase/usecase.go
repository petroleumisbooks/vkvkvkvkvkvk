package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	h "internbot/internal/handler"
	"net/http"
	"time"
)

type UseCase struct {
	Token   string
	Runing  bool
	Handler h.Handler
}

func (us *UseCase) sendMessage(messageBytes []byte) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", us.Token)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (us *UseCase) getUpdates(offset int64) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", us.Token, offset)
}

func (us *UseCase) Run() {
	var lastId int64 = 0

	for us.Runing {
		time.Sleep(time.Duration(1) * time.Second)

		url := us.getUpdates(lastId)

		response, err := http.Get(url)
		if err != nil {
			continue
		}
		defer response.Body.Close()

		message, err := us.Handler.Control(response, &lastId)
		if err != nil {
			continue
		}

		messageBytes, err := json.Marshal(message)
		if err != nil {
			continue
		}

		err = us.sendMessage(messageBytes)
		if err != nil {
			continue
		}
	}

}

func (us *UseCase) Shutdown() {
	us.Runing = false
}
