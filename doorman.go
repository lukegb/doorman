package main

import (
	"github.com/lukegb/doorman/gpio"
	"net/http"
	"fmt"
	"log"
	"time"
	"sync"
)

const (
	GPIO_PIN = 4
)

func runAndServe(pin gpio.Pin) error {
	pinLock := sync.NewCond(new(sync.Mutex))
	pinOpen := false

	http.HandleFunc("/door/front/open", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Add("Allow", "POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `must be POST`)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		pinLock.L.Lock()
		defer pinLock.L.Unlock()
		if pinOpen {
			pinLock.Wait()
		} else {
			pinOpen = true
			pinLock.L.Unlock()
			pin.SetValue(true)
			time.Sleep(5 * time.Second)
			pin.SetValue(false)
			pinLock.L.Lock()
			pinOpen = false
			pinLock.Broadcast()
		}
		fmt.Fprintf(w, `{"message": "Door opened"}`)
	})
	return http.ListenAndServe(":8080", nil)
}

func main() {
	pin, err := gpio.BuildGpioPin(GPIO_PIN)
	if err != nil {
		log.Fatalln(err)
	}

	pin.SetDirection(gpio.GPIO_OUT)

	log.Println("SERVING!")
	log.Fatal(runAndServe(pin))
}
