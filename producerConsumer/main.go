package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= numberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d, it will take %d seconds\n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d ***\n", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d ***\n", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza #%d is ready!\n", pizzaNumber)
		}

		p := &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
		return p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	i := 0

	// run forever or until we receive a quit notification
	// try to make pizzas
	for {
		// try to make a pizza
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// we tried to make a pizza
			case pizzaMaker.data <- *currentPizza:
				// fmt.Println("here")
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}

}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	fmt.Println("The pizzeria is now open for business")
	fmt.Println("-------------------------------------")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzeria((pizzaJob))

	// create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= numberOfPizzas {
			if i.success {
				fmt.Println(i.message)
				fmt.Printf("Order #%d is out for delivery:\n", i.pizzaNumber)
			} else {
				fmt.Println(i.message)
				fmt.Printf("Order %d failed,the customer is really mad\n", i.pizzaNumber)
			}
		} else {
			fmt.Println("done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				fmt.Println("Error closing channel", err)
			}
		}
	}

	// print out the ending message
	fmt.Println("-------------------------------------")
	fmt.Println("Done for the day")

	fmt.Printf("We made %d pizzas, %d failed, with %d attempts in total\n", pizzasMade, pizzasFailed, total)

	switch {
	case pizzasFailed > 9:
		fmt.Println("It was an awful day")
	case pizzasFailed >= 5:
		fmt.Println("It was not a very good day")
	case pizzasFailed >= 2:
		fmt.Println("It was a pretty good day")
	default:
		fmt.Println("It was a great day")
	}
}
