package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// точка входа
func main() {
	botToken := ""
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0

	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth wetn wrong: ", err.Error())
		}
		for _, update := range updates {
			err = response(botUrl, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

// запрос обновлений
func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

// ответ на обновления
func response(botUrl string, update Update) error {
	var BotMessage BotMessage
	BotMessage.ChatId = update.Message.Chat.ChatId
	BotMessage.Text = update.Message.Text
	buf, err := json.Marshal(BotMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
