package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	// "io/ioutil"
	// "encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
)

// type device struct {
// 	// Gps_num float32 `json:"gps_num"`
// 	// App string `json:"app"`
// 	// Gps_alt float32 `json:"gps_alt"`
// 	// Fmt_opt int `json:"fmt_opt"`
// 	// Device string `json:"device"`
// 	// S_d2 float32 `json:"s_d2"`
// 	S_d0 float32 `json:"s_d0"`
// 	// S_d1 float32 `json:"s_d1"`
// 	S_h0 float32 `json:"s_h0"`
// 	SiteName string `json:"SiteName"`
// 	// Gps_fix float32 `json:"gps_fix"`
// 	// Ver_app string `json:"ver_app"`
// 	Gps_lat float32 `json:"gps_lat"`
// 	// S_t0 float32 `json:"s_t0"`
// 	Timestamp string `json:"timestamp"`
// 	Gps_lon float32 `json:"gps_lon"`
// 	// Date string `json:"date"`
// 	// Tick float32 `json:"tick"`
// 	Device_id string `json:"device_id"`
// 	// S_1 float32 `json:"s_1"`
// 	// S_0 float32 `json:"s_0"`
// 	// S_3 float32 `json:"s_3"`
// 	// S_2 float32 `json:"s_2"`
// 	// Ver_format string `json:"ver_format"`
// 	// Time string `json:"time"`
// }

// type airbox struct {
// 	Source string `json:"source"`
// 	Feeds []device `json:"feeds"`
// 	Version string `json:"version"`
// 	Num_of_records int `json:"num_of_records"`
// }

var bot *linebot.Client
// var airbox_json airbox

func main() {
	// url := "https://data.lass-net.org/data/last-all-airbox.json"
	// req, _ := http.NewRequest("GET", url, nil)
	// res, _ := http.DefaultClient.Do(req)
	// defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)
	// err := json.Unmarshal(body, &airbox_json)
	// if err != nil {
	// 	fmt.Println(err)
	// }
    // fmt.Println(airbox_json)

	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
