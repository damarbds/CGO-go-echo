package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	_echoMiddleware "github.com/labstack/echo/middleware"
)

func main() {
	//dev
	baseUrlLocal := "http://cgo-web-api.azurewebsites.net"
	//prd
	// baseUrlLocal := "https://api-cgo-prod.azurewebsites.net"
	//local
	//baseUrlLocal := "http://localhost:9090"
	e := echo.New()
	//middL := middleware.InitMiddleware()

	//e.Use(middL.CORS)
	e.Use(_echoMiddleware.CORSWithConfig(_echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	go UpdateStatusExpiredPaymentJob(baseUrlLocal)
	go RemainingPaymentJob(baseUrlLocal)
	log.Fatal(e.Start(":9090"))

}
func UpdateStatusExpiredPaymentJob(baseUrl string) {
	done := make(chan bool)
	ticker := time.NewTicker(time.Hour)

	go func() {
		//time.Sleep(10 * time.Second) // wait for 10 seconds
		//done <- true
	}()

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case t := <-ticker.C:
			req, err := http.NewRequest("POST", baseUrl+"/booking/update-expired-payment", nil)

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
			fmt.Println(resp.Body)
			fmt.Println("Current time: ", t.String())
		}
	}
}

func RemainingPaymentJob(baseUrl string) {
	done := make(chan bool)
	ticker := time.NewTicker(time.Hour * 24)

	go func() {
		//time.Sleep(10 * time.Second) // wait for 10 seconds
		//done <- true
	}()

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case t := <-ticker.C:
			req, err := http.NewRequest("POST", baseUrl+"/booking/remaining-payment-booking", nil)

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
			fmt.Println(resp.Body)
			fmt.Println("Current time: ", t.String())
		}
	}
}
