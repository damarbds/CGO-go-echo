package paypal

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PayOrderDetail struct {
	ID            string            `json:"id"`
	Status        string            `json:"status"`
	Intent        string            `json:"intent"`
	PurchaseUnits []PayPurchaseUnit `json:"purchase_units"`
	CreateTime    string            `json:"create_time"`
}

type PayPurchaseUnit struct {
	ReferenceID string    `json:"reference_id"`
	Amount      PayAmount `json:"amount"`
}

type PayAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PaypalToken struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   uint64 `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

type PaypalConfig struct {
	OAuthUrl string
	OrderUrl string
}

const (
	defHTTPClientTimeout = 15 * time.Second
	defClientID          = "AUQh6fyxaSlxELSff1NUlBaBskgh5emI7MnGwnm68xF1lUJ5jnmPXFcKNDp4D5ZmNaXRyA1ONKnzavKt"
	defSecretKey         = "ENH-GqjLFOmJUyHKwNA_eMs1s8N5qbKcSCl2zGOvyIYf-1M_G4fK_Yk18WOcrD15J2B_d1mzeVKlb0zt"
	defClientIDProd      = "Aer1l0Q_pqnw_S0Opw-Qd2VMEPWFuDD05xW6hGC1VhAv4NnKvnSi6JS93B5ZGDO1HYqRHNFTerjo4Aos"
	defSecretKeyProd     = "EIhYAOGV5EIHWfQ85aRo0KopugZmDyUTx3fZp9nWXGJ5EYwZ3VzRirm6BDphz-U7PGzBSrEDYpNd5MBZ"
)

var (
	PaypalOauthUrl = "https://api.sandbox.paypal.com/v1/oauth2/token/"
	PaypalOrderUrl = "https://api.sandbox.paypal.com/v2/checkout/orders/"
)

func AuthPaypal(paypalOauthUrl string) (string, error) {
	client := &http.Client{Timeout: defHTTPClientTimeout}
	res := PaypalToken{}

	input := []byte(defClientID + ":" + defSecretKey)

	AUTH_STRING := b64.StdEncoding.EncodeToString(input)

	data := url.Values{}
	data.Add("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, paypalOauthUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return res.AccessToken, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Basic "+AUTH_STRING)

	rsp, err := client.Do(req)
	if err != nil {
		return res.AccessToken, err
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return res.AccessToken, err
	}
	if rsp.StatusCode >= 300 {
		return res.AccessToken, errors.New("status_code=" + rsp.Status + " body=" + string(rspBody) + "input=" + string(input) + "auth string=" + AUTH_STRING)
	}

	// decode response
	err = json.Unmarshal(rspBody, &res)
	if err != nil {
		return res.AccessToken, err
	}

	return res.AccessToken, nil
}

func PaypalSetup(cfg PaypalConfig, orderID string) (PayOrderDetail, error) {
	client := &http.Client{Timeout: defHTTPClientTimeout}
	res := PayOrderDetail{}

	accessToken, err := AuthPaypal(cfg.OAuthUrl)
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest(http.MethodGet, cfg.OrderUrl+orderID, nil)
	if err != nil {
		return res, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	rsp, err := client.Do(req)
	if err != nil {
		return res, err
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return res, err
	}
	if rsp.StatusCode >= 300 {
		return res, errors.New("status_code=" + rsp.Status + " body=" + string(rspBody) + "msg= paypal setup")
	}

	// decode response
	err = json.Unmarshal(rspBody, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
