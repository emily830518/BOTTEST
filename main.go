package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/go-redis/redis"
	"strconv"
	// "time"
)

type device struct {
	// Gps_num float32 `json:"gps_num"`
	// App string `json:"app"`
	// Gps_alt float32 `json:"gps_alt"`
	// Fmt_opt int `json:"fmt_opt"`
	// Device string `json:"device"`
	// S_d2 float32 `json:"s_d2"`
	S_d0 float32 `json:"s_d0"`
	// S_d1 float32 `json:"s_d1"`
	S_h0 float32 `json:"s_h0"`
	SiteName string `json:"SiteName"`
	// Gps_fix float32 `json:"gps_fix"`
	// Ver_app string `json:"ver_app"`
	Gps_lat float32 `json:"gps_lat"`
	S_t0 float32 `json:"s_t0"`
	Timestamp string `json:"timestamp"`
	Gps_lon float32 `json:"gps_lon"`
	// Date string `json:"date"`
	// Tick float32 `json:"tick"`
	Device_id string `json:"device_id"`
	// S_1 float32 `json:"s_1"`
	// S_0 float32 `json:"s_0"`
	// S_3 float32 `json:"s_3"`
	// S_2 float32 `json:"s_2"`
	// Ver_format string `json:"ver_format"`
	// Time string `json:"time"`
}

type airbox struct {
	Source string `json:"source"`
	Feeds []device `json:"feeds"`
	Version string `json:"version"`
	Num_of_records int `json:"num_of_records"`
}

type subscribeid struct {
	Device_id []string `json:"device_id"`
}

var bot *linebot.Client
var airbox_json airbox
var lass_json airbox
var maps_json airbox
var all_device []device
var history_json subscribeid
var	client=redis.NewClient(&redis.Options{
		Addr:"hipposerver.ddns.net:6379",
		Password:"",
		DB:0,
	})

func main() {
	url := "https://data.lass-net.org/data/last-all-airbox.json"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	errs := json.Unmarshal(body, &airbox_json)
	if errs != nil {
		fmt.Println(errs)
	}

	url = "https://data.lass-net.org/data/last-all-lass.json"
	req, _ = http.NewRequest("GET", url, nil)
	res, _ = http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	errs = json.Unmarshal(body, &lass_json)
	if errs != nil {
		fmt.Println(errs)
	}

	url = "https://data.lass-net.org/data/last-all-maps.json"
	req, _ = http.NewRequest("GET", url, nil)
	res, _ = http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	errs = json.Unmarshal(body, &maps_json)
	if errs != nil {
		fmt.Println(errs)
	}

	all_device=append(maps_json.Feeds,lass_json.Feeds...)
	all_device=append(all_device,airbox_json.Feeds...)

	url = "https://data.lass-net.org/data/airbox_list.json"
	req, _ = http.NewRequest("GET", url, nil)
	res, _ = http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ = ioutil.ReadAll(res.Body)
	errs = json.Unmarshal(body, &history_json)
	if errs != nil {
		fmt.Println(errs)
	}
	// pushmessage()
	// fmt.Println(airbox_json)
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	// _,_=bot.PushMessage("U3617adbdd46283d7e859f36302f4f471", linebot.NewTextMessage("hi!")).Do()
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
				var txtmessage string
				inText := strings.ToLower(message.Text)
				if strings.Contains(inText,"訂閱"){
					userID:=event.Source.UserID
					// pong, _ := client.Ping().Result()
					// txtmessage=pong
					for i:=0; i<len(history_json.Device_id); i++ {
						if strings.Contains(inText,strings.ToLower(history_json.Device_id[i])) {
							val, err:=client.Get(history_json.Device_id[i]).Result()
							if err!=nil{
								client.Set(history_json.Device_id[i],userID,0)
								txtmessage="訂閱成功!"
								break
							}
							if strings.Contains(inText,"取消"){
								stringSlice:=strings.Split(val,",")
								if stringInSlice(userID,stringSlice){
									if len(stringSlice)==1{
										client.Del(history_json.Device_id[i])
										txtmessage="取消訂閱成功!"
										break
									}else{
										var s []string
										s = removeStringInSlice(stringSlice, userID)
										client.Set(history_json.Device_id[i],s,0)
										txtmessage="取消訂閱成功!"
										break
									}
								}else{
									txtmessage="你並沒有訂閱此ID。"
									break
								}
							}
							stringSlice:=strings.Split(val,",")
							if stringInSlice(userID,stringSlice){
								txtmessage="您已訂閱過此ID!"
								break
							} else{
								val=val+","+userID
								client.Set(history_json.Device_id[i],val,0)
								txtmessage="訂閱成功!"
								break
							}
						}
					}
				} else{
					for i:=0; i<len(all_device); i++ {
						if strings.Contains(inText,strings.ToLower(all_device[i].Device_id)) {
							txtmessage="Device_id: "+all_device[i].Device_id+"\n"
							txtmessage=txtmessage+"Site Name: "+all_device[i].SiteName+"\n"
							txtmessage=txtmessage+"Location: ("+strconv.FormatFloat(float64(all_device[i].Gps_lon),'f',3,64)+","+strconv.FormatFloat(float64(all_device[i].Gps_lat),'f',3,64)+")"+"\n"
							txtmessage=txtmessage+"Timestamp: "+all_device[i].Timestamp+"\n"
							txtmessage=txtmessage+"PM2.5: "+strconv.FormatFloat(float64(all_device[i].S_d0),'f',0,64)+"\n"
							txtmessage=txtmessage+"Humidity: "+strconv.FormatFloat(float64(all_device[i].S_h0),'f',0,64)+"\n"
							txtmessage=txtmessage+"Temperature: "+strconv.FormatFloat(float64(all_device[i].S_t0),'f',0,64)
							break
						}
					}
				}
				if len(txtmessage)==0{
					txtmessage="Sorry! No this device ID, please check again."
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(txtmessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func removeStringInSlice(s []string, r string) []string {
    for i, v := range s {
        if v == r {
            return append(s[:i], s[i+1:]...)
        }
    }
    return s
}

