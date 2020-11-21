package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const SearchUrl = "https://www.googleapis.com/youtube/v3/search"
const YoutubeToken = ""
const YoutubeVideoUrl = "https://www.youtube.com/watch?v="
const TelegramToken = ""
const TelegramApi = "https://api.telegram.org/bot"

func main() {
	botUrl := TelegramApi + TelegramToken
	offset := 0

	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = response(botUrl, update)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

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

func response(botUrl string, update Update) error {
	var BotMessage BotMessage
	BotMessage.ChatId = update.Message.Chat.ChatId
	videoUrls, err := getLastVideos(update.Message.Text)
	if err != nil {
		return err
	}
	for _, element := range videoUrls {
		BotMessage.Text = YoutubeVideoUrl + element.Id.VideoId
		buf, err := json.Marshal(BotMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
	}
	return nil
}

func getLastVideos(Query string) ([]Item, error) {
	items, err := retriveVideos(Query)
	if err != nil {
		return nil, err
	}
	if len(items) < 1 {
		return nil, errors.New("No videos found")
	}
	return items, nil
}

func retriveVideos(Query string) ([]Item, error) {
	req, err := makeRequests(Query, 5)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var restResponse RestYoutubeResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Items, nil
}

func makeRequests(Query string, maxResults int) (*http.Request, error) {
	lastSlashIndex := strings.LastIndex(Query, "/")
	videoId := Query[lastSlashIndex+1:]
	req, err := http.NewRequest("GET", SearchUrl, nil)
	if err != nil {
		return nil, err
	}
	queryStr := req.URL.Query()
	queryStr.Add("part", "id")
	queryStr.Add("maxResults", strconv.Itoa(maxResults))
	queryStr.Add("order", "relevance")
	queryStr.Add("q", videoId)
	queryStr.Add("key", YoutubeToken)
	req.URL.RawQuery = queryStr.Encode()
	return req, nil
}
