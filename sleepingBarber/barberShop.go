package main

import (
	"fmt"
	"time"
)

type BarberShop struct {
	ShopCapacity    int
	HaircutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (shop *BarberShop) addBarber(barber string) {
	shop.NumberOfBarbers++

	go func() {
		isSleeping := false
		fmt.Printf("%s goes to the waiting room to check for clients\n", barber)

		for {
			// if there are no clients, then the barber goes to sleep
			if len(shop.ClientsChan) == 0 {
				fmt.Printf("There is nothing to do, so %s goes to sleep\n", barber)
				isSleeping = true
			}

			client, shopOpen := <-shop.ClientsChan
			if shopOpen {
				if isSleeping {
					fmt.Printf("%s wakes %s up", client, barber)
					isSleeping = false
				}
				// cut hair
				shop.cutHair(barber, client)
			} else {
				// barber goes home
				shop.sendBarberHome(barber)
				return
			}
		}
	}()
}

func (shop *BarberShop) cutHair(barber, client string) {
	fmt.Printf("%s is cutting %s's hair\n", barber, client)
	time.Sleep(shop.HaircutDuration)
	fmt.Printf("%s has finished cutting %s's hair\n", barber, client)
}

func (shop *BarberShop) sendBarberHome(barber string) {
	fmt.Printf("%s goes home\n", barber)
	shop.BarbersDoneChan <- true
}

func (shop *BarberShop) closeShopForDay() {
	fmt.Println("Closing the barbershop for the day")

	// close the clients channel
	close(shop.ClientsChan)

	shop.Open = false

	// wait for all barbers to finish
	for a := 1; a <= shop.NumberOfBarbers; a++ {
		<-shop.BarbersDoneChan
	}
	close(shop.BarbersDoneChan)

	fmt.Println("The barbershop is now closed for the day")
	fmt.Println("-----------------------------------------")
}

func (shop *BarberShop) addClient(client string) {
	// print out a message
	fmt.Printf("*** client %s has arrived\n", client)

	if shop.Open {
		select {
		case shop.ClientsChan <- client:
			fmt.Printf("%s takes a seat in the waiting room\n", client)
		default:
			fmt.Printf("%s sees that the waiting room is full and leaves\n", client)
		}
	} else {
		fmt.Printf("%s sees that the barbershop is closed and leaves\n", client)
	}
}
