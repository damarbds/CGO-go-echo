package usecase

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/identityserver"
	"github.com/models"
	"net/http"
	"net/url"
	"os"
)

type identityserverUsecase struct {
	baseUrl				string
	basicAuth 			string
}

// NewidentityserverUsecase will create new an identityserverUsecase object representation of identityserver.Usecase interface
func NewidentityserverUsecase(baseUrl string,basicAuth string) identityserver.Usecase {
	return &identityserverUsecase{
		baseUrl: 			baseUrl,
		basicAuth:			basicAuth,
	}
}
func (m identityserverUsecase) GetUserInfo(token string) (*models.GetUserInfo, error) {

	req, err := http.NewRequest("POST", m.baseUrl + "/connect/userinfo", nil)
	//os.Exit(1)
	req.Header.Set("Authorization", "Bearer " + token)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrUsernamePassword
	}
	user := models.GetUserInfo{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

func (m identityserverUsecase) GetToken(username string, password string) (*models.GetToken, error) {

	var param = url.Values{}
	param.Set("grant_type", "password")
	param.Set("username", username)
	param.Set("password", password)
	param.Set("scope", "openid")
	var payload = bytes.NewBufferString(param.Encode())

	req, err := http.NewRequest("POST", m.baseUrl + "/connect/token", payload)
	//os.Exit(1)
	req.Header.Set("Authorization", "Basic " + m.basicAuth)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil  {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrUsernamePassword
	}
	user := models.GetToken{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}


func (m identityserverUsecase) UpdateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser,error) {

	data, _:= json.Marshal(ar)

	req, err := http.NewRequest("POST", m.baseUrl + "/connect/update-user", bytes.NewReader(data))
	//os.Exit(1)
	//req.Header.Set("Authorization", "Basic YWRtaW5AZ21haWwuY29tOmFkbWlu")
	req.Header.Set("Content-Type","application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrBadParamInput
	}
	user := models.RegisterAndUpdateUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

func (m identityserverUsecase) CreateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser, error) {
	data, _:= json.Marshal(ar)
	req, err := http.NewRequest("POST", m.baseUrl + "/connect/register", bytes.NewReader(data))
	//os.Exit(1)
	//req.Header.Set("Authorization", "Basic YWRtaW5AZ21haWwuY29tOmFkbWlu")
	req.Header.Set("Content-Type","application/json")
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error : ", err.Error())
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		return nil, models.ErrBadParamInput
	}
	user := models.RegisterAndUpdateUser{}
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
