package events

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/models"
	"net/http"
	"os"
	"time"
)

func init() {
	//createNotifier := userCreatedNotifier{
	//	adminEmail: "the.boss@example.com",
	//	slackHook: "https://hooks.slack.com/services/...",
	//}

	//UserCreated.Register(createNotifier)
}

type UserCreatedNotifier struct{
	BaseUrl string	`json:"base_url"`
	Schedule models.NewCommandSchedule `json:"schedule"`
}

func (u UserCreatedNotifier) notifyAdmin(email string, time time.Time) {
	// send a message to the admin that a user was created
}

func (t UserCreatedNotifier) sendToSlack(ar models.NewCommandSchedule) {
	data, _ := json.Marshal(ar)

	req, err := http.NewRequest("POST", t.BaseUrl+"/service/schedule", bytes.NewReader(data))
	//os.Exit(1)
	//req.Header.Set("Authorization", "Basic YWRtaW5AZ21haWwuY29tOmFkbWlu")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		//return nil, models.ErrBadParamInput
	}
	user := models.NewCommandSchedule{}
	json.NewDecoder(resp.Body).Decode(&user)

}

func (u UserCreatedNotifier) Handle(payload UserCreatedPayload) {
	// Do something with our payload
	//u.notifyAdmin(payload.Email, payload.Time)
	u.sendToSlack(payload.Schedule)
}
