package main

import (
	"github.com/deckarep/gosx-notifier"
	coinmarketcapApi "github.com/miguelmota/go-coinmarketcap"
	"strings"
	"time"
	"log"
)

//List of coins check
var coins []string = []string{
	"stellar",
	"ripple",
	"cardano",
	"siacoin",
	"ubiq",
	"electra",
	"myriad"}


func main() {
	resultCh := make(chan string)
  errorCh := make(chan error)

	for _, s := range coins {
		go checkCoin(resultCh, errorCh, s)
	}
	for {
		select {
		case result := <-resultCh:
			if strings.HasPrefix(result, "n") {
				s := strings.Split(strings.Trim(result, "n"), ",")
				sendNotif(s[0], "Change is " + s[1])
			}
		case error:= <- errorCh:
			log.Println(error);
		}
	}
}

func sendNotif(title string, message string) {

	note := gosxnotifier.NewNotification(message)
	note.Title = "Crypto: " + title
	note.Sound = gosxnotifier.Blow

	note.Push()
}

func checkCoin(resultCh chan string, errorCh chan error, coin string) {
	var lastChangeNotif float64 = 0

	for {
		for {
				coinInfo, err := coinmarketcapApi.GetCoinData(coin)
				if err != nil {
					errorCh <- err
				} else {
					if ((lastChangeNotif == 0 || lastChangeNotif < coinInfo.PercentChange1h) && coinInfo.PercentChange1h > 1) {
						resultCh <- "n " + coin + "," + "UP"
						lastChangeNotif = coinInfo.PercentChange1h
					} else if ((lastChangeNotif == 0 || lastChangeNotif > coinInfo.PercentChange1h) && coinInfo.PercentChange1h < -1) {
						resultCh <- "n " + coin + "," + "DOWN"
						lastChangeNotif = coinInfo.PercentChange1h
					}
				}
				time.Sleep(20 * time.Second)
			}
	}
}
