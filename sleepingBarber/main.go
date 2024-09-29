package main

import (
	"fmt"
	"math/rand"
	"time"
)

// variables
var (
	seatingCapacity = 2
	arrivalRate     = 100
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10 * time.Second
)

func main() {
	// seed our random number generator
	rand.Seed(time.Now().UnixNano())

	// print our welcome message
	fmt.Println("The sleeping barber problem")
	fmt.Println("---------------------------")

	// create our channels
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HaircutDuration: cutDuration,
		BarbersDoneChan: doneChan,
		ClientsChan:     clientChan,
		Open:            true,
		NumberOfBarbers: 0,
	}

	fmt.Println("The barbershop is open for business")

	// add barbers
	shop.addBarber("Frank")
	shop.addBarber("Gerrard")
	shop.addBarber("Romeo")

	// start the barbershop as a goroutine
	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	// add clients
	i := 1

	go func() {
		for {
			// get a random number for average arrival rate
			randomMilliSeconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Duration(randomMilliSeconds) * time.Millisecond):
				shop.addClient(fmt.Sprintf("Client%d", i))
				i++
			}

		}
	}()

	// block until the barbershop is closed
	<-closed
	// time.Sleep(5 * time.Second)
}
