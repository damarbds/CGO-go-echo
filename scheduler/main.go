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
	baseUrlLocalPRD := "https://api-cgo-prod.azurewebsites.net"
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
	//dev
	go ChangeStatusJob(baseUrlLocal)
	go UpdateStatusExpiredPaymentJob(baseUrlLocal)
	go RemainingPaymentJob(baseUrlLocal)
	go CreateExChangeJob(baseUrlLocal)

	//prd
	// go ChangeStatusJobPRD(baseUrlLocalPRD)
	go UpdateStatusExpiredPaymentJobPRD(baseUrlLocalPRD)
	go RemainingPaymentJobPRD(baseUrlLocalPRD)
	go CreateExChangeJobPRD(baseUrlLocalPRD)
	log.Fatal(e.Start(":9090"))

}
func CreateExChangeJob(baseUrl string) {
	done := make(chan bool)
	//ticker := time.NewTicker(time.Hour * 3)
	ticker := time.NewTicker(time.Second * 1)
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
			time.Sleep(3 * time.Hour)
			req, err := http.NewRequest("POST", baseUrl+"/misc/exchange-rate", nil)

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

func ChangeStatusJob(baseUrl string) {
	done := make(chan bool)
	//ticker := time.NewTicker(time.Hour)
	ticker := time.NewTicker(time.Second * 1)
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
			now := time.Now().Format("15:04")
			if now == "01:00" {
				req, err := http.NewRequest("POST", baseUrl+"/booking/changes-status-scheduler", nil)

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
}

func UpdateStatusExpiredPaymentJob(baseUrl string) {
	done := make(chan bool)
	//ticker := time.NewTicker(time.Hour)
	ticker := time.NewTicker(time.Second * 1)
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
			time.Sleep(1 * time.Hour) // wait for 10 seconds
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
	//ticker := time.NewTicker(time.Hour * 24)
	ticker := time.NewTicker(time.Second * 1)
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
			time.Sleep(24 * time.Hour) // wait for 10 seconds
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

func ChangeStatusJobPRD(baseUrl string) {
	done := make(chan bool)
	//ticker := time.NewTicker(time.Hour)
	ticker := time.NewTicker(time.Second * 1)
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
			now := time.Now().Format("15:04")
			if now == "01:00" {
				req, err := http.NewRequest("POST", baseUrl+"/booking/changes-status-scheduler", nil)

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
}

func CreateExChangeJobPRD(baseUrl string) {
	done := make(chan bool)
	//ticker := time.NewTicker(time.Hour * 4)
	ticker := time.NewTicker(time.Second * 1)
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
			time.Sleep(4 * time.Hour) // wait for 10 seconds
			req, err := http.NewRequest("POST", baseUrl+"/misc/exchange-rate", nil)

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

func UpdateStatusExpiredPaymentJobPRD(baseUrl string) {
	done := make(chan bool)
	//ticker := time.NewTicker(time.Hour * 2)
	ticker := time.NewTicker(time.Second * 1)

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
			time.Sleep(2 * time.Hour) // wait for 10 seconds
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

func RemainingPaymentJobPRD(baseUrl string) {
	done := make(chan bool)
	//ticker := time.NewTicker(time.Hour * 25)
	ticker := time.NewTicker(time.Second * 1)
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
			time.Sleep(25 * time.Hour) // wait for 10 seconds
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
