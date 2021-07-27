package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

//Go script for vaccine slot availability
type serve interface {
	Address() string
	item() int
}
type smtpServer struct {
	host string
	port string
}

func (s smtpServer) Address() string {
	return s.host + ":" + s.port
}
func sendEmail() {
	from := "******42@gmail.com"
	pass := "password"
	to := []string{"****@gmail.com"}
	smtServer := smtpServer{host: "smtp.gmail.com", port: "587"}
	message := []byte("Hi!\n Vaccine Slots are now available")
	auth := smtp.PlainAuth("", from, pass, smtServer.host)
	err := smtp.SendMail(smtServer.Address(), auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email sent with slot details")
}
func main() {
	date := time.Now().Format("01-02-2006")
	dater := strings.Split(date, "-")
	date = dater[1] + "-" + dater[0] + "-" + dater[2]
	fmt.Println("Checking for vaccination slot avialability for " + date + " ...\n")
	resp, err := http.Get("https://cdn-api.co-vin.in/api/v2/appointment/sessions/public/findByPin?pincode=457001&date=" + date)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var found bool = false
	message_adder := ""
	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)
	centres := result["sessions"].([]interface{})
	for i := range centres {
		centre := centres[i].(map[string]interface{})
		if centre["available_capacity_dose1"].(float64) > 0 {
			found = true
			message_adder += "Available at centre :  " + centre["address"].(string) + "\n"
		}
	}
	if found {
		sendEmail()
	} else {
		fmt.Println("No vaccination slot is available. Try again later !")
	}

}
