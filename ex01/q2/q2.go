package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Produtor struct {
	creq chan (chan int)
}

type Consumidor struct {
	creq chan (chan int)
	id   int
}

func (prod *Produtor) produzir(n int) {
	fmt.Println("P: comidinha", n)
	if sendTo := len(prod.creq); sendTo > 0 {
		for i := 0; i < sendTo; i++ {
			chreq := <-prod.creq
			chreq <- n
		}
		fmt.Println("P: enviou", n, "para", sendTo, "consumidores")
	} else {
		fmt.Println("P: bloqueado")
		chreq := <-prod.creq
		chreq <- n
		fmt.Println("P: enviou", n)
	}
}

func (cons *Consumidor) consumir() int {
	chreq := make(chan int)
	cons.creq <- chreq
	fmt.Printf("C%d: bloqueado\n", cons.id)
	n := <-chreq
	close(chreq)
	return n
}

func produz(prod *Produtor, maxProduzidas int) {
	for i := 0; i < maxProduzidas; i++ {
		prod.produzir(rand.Int())
		// esse temporizador serve para ilustracoes
		//time.Sleep(time.Millisecond * 500)
	}
	fmt.Println("enviando sinal de finalizacao de programa")
	finish <- 1
}

func invocaConsumidores(creq chan (chan int)) {
	consumers := 1 + rand.Intn(10)
	fmt.Println("qtd de consumidores:", consumers)
	for i := 0; i < consumers; i++ {
		go func(id int) {
			cons := Consumidor{creq: creq, id: id}
			for i := 0; ; i++ {
				n := cons.consumir()
				fmt.Printf("C%d pegou %d (vez: %d)\n", id, n, i)
				// esse temporizador serve para ilustracoes
				//time.Sleep(time.Millisecond * 300)
			}
		}(i)
	}
}

var finish = make(chan byte)

func main() {
	/* */
	rand.Seed(time.Now().Unix())
	creq := make(chan (chan int), 10)
	prod := Produtor{creq}

	maxProduzidas := 10

	invocaConsumidores(creq)
	time.Sleep(1500 * time.Millisecond)
	go produz(&prod, maxProduzidas)
	/* /
	exemplo1()
	time.Sleep(time.Second)
	exemplo2()
	/* */
	<-finish
	fmt.Println("fim")
}

func exemplo1() {
	fmt.Println("EXEMPLO 1")

	creq := make(chan (chan int), 10)
	prod := Produtor{creq}
	cons1 := Consumidor{creq, 1}
	cons2 := Consumidor{creq, 2}

	go func() {
		fmt.Println("C1 consumiu", cons1.consumir())
	}()
	go func() {
		fmt.Println("C2 consumiu", cons2.consumir())
	}()
	go prod.produzir(13)
}

func exemplo2() {
	fmt.Println("EXEMPLO 2")

	creq := make(chan (chan int), 10)
	prod := Produtor{creq}
	cons1 := Consumidor{creq, 1}
	cons2 := Consumidor{creq, 2}
	go prod.produzir(23)

	go func() {
		fmt.Println("C1 consumiu", cons1.consumir())
	}()
	go func() {
		fmt.Println("C2 consumiu", cons2.consumir())
	}()
}
