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

type PaymentNotifier struct{
	BaseUrl string	`json:"base_url"`
	ConfirmPaymentByDate *models.ConfirmTransactionPayment `json:"confirm_payment_by_date"`
	ConfirmPayment 	*models.ConfirmPaymentIn `json:"confirm_payment"`
}

func (u PaymentNotifier) notifyAdmin(email string, time time.Time) {
	// send a message to the admin that a user was created
}

func (t PaymentNotifier) sendToSlackConfirmPaymentByDate(ar *models.ConfirmTransactionPayment) {
	data, _ := json.Marshal(ar)

	req, err := http.NewRequest("PUT", t.BaseUrl+"/transaction/payments/sending-email-confirm-by-date", bytes.NewReader(data))
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
	//user := models.NewCommandSchedule{}
	//json.NewDecoder(resp.Body).Decode(&user)

}
func (t PaymentNotifier) sendToSlackConfirmPayment(ar *models.ConfirmPaymentIn) {
	data, _ := json.Marshal(ar)

	req, err := http.NewRequest("PUT", t.BaseUrl+"/transaction/payments/sending-email-confirm", bytes.NewReader(data))
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
	//user := models.NewCommandSchedule{}
	//json.NewDecoder(resp.Body).Decode(&user)

}

func (u PaymentNotifier) Handle(payload PaymentPayload) {
	// Do something with our payload
	//u.notifyAdmin(payload.Email, payload.Time)
	if payload.ConfirmPaymentByDate != nil {
		u.sendToSlackConfirmPaymentByDate(payload.ConfirmPaymentByDate)
	}
	if payload.ConfirmPayment != nil {
		u.sendToSlackConfirmPayment(payload.ConfirmPayment)
	}
}
