package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
)

type API struct {
}

func (a API) GetCurrentBTCUAH() (float64, error) {
	url := "https://api.coinbase.com/v2/exchange-rates?currency=BTC"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(res.Body)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var result map[string]map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}
	price := result["data"]["rates"].(map[string]interface{})["UAH"]
	priceUAH, err := strconv.ParseFloat(fmt.Sprint(price), 64)
	if err != nil {
		return 0, err
	}
	return priceUAH, nil
}

func (a API) SendMail(toEmails []string, message string) []string {
	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	//from := "timofeyka.com.03@gmail.com"
	//password := "wbpczxponwnswjro"
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := "Subject: BTC current rate\r\n\r\n"
	body := fmt.Sprintf("1 BTC = %s UAH", message)
	msg := []byte(subject + body)
	auth := smtp.PlainAuth("", from, password, host)
	var invalidEmails []string
	for _, v := range toEmails {
		err := smtp.SendMail(address, auth, from, []string{v}, msg)
		if err != nil {
			invalidEmails = append(invalidEmails, v)
		}
		log.Printf("Sending to %s ...", v)
	}
	return invalidEmails
}

func (a API) GetRate() (int, error) {
	url := "http://localhost:8000/api/rate"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if res.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(res.Body)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var result int
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}
	return result, nil
}
