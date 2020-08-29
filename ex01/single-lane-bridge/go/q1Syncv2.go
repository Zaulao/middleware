package main

import (
	"fmt"
	"sync"
	"time"
)

// temVindo e temIndo estao aqui apenas para proposito de visualizacao.
// serve para checar que nao tem um carro passando de la pra ca e um
// daqui pra la ao mesmo tempo. o mux eh apenas para evitar race condition
// no manuseio dessas variaveis ajudantes, o problema seria resolvido igual
// sem temIndo, temVindo e o mux que os protege.

// Esta versao tambem controla acesso a ponte. Diferentemente de q1Sync.go,
// usa Cond para criar um semaforo que fica permite cada lado passar por um
// intervalo de tempo especifico. Uma enxurrada de carros vem, depois volta
// Em q1Sync.go o acesso eh aleatorio: quem conseguir botar a cara primeiro
type PonteOrganizada struct {
	deLaPraCa  chan string
	daquiPraLa chan string
	temIndo    bool
	temVindo   bool
	c          *sync.Cond
	direcao    string
	mux        sync.Mutex
}

func (ponte *PonteOrganizada) entraPorLa() {
	for {
		ponte.c.L.Lock()
		for ponte.direcao != "lapraca" {
			ponte.c.Wait()
		}
		ponte.deLaPraCa <- "Entra por la"

		ponte.mux.Lock()
		ponte.temVindo = true
		ponte.mux.Unlock()
	}
}

func (ponte *PonteOrganizada) saiPorAqui() {
	for {
		carro := <-ponte.deLaPraCa
		ponte.mux.Lock()
		fmt.Println(carro, "e sai por aqui. batida?", ponte.temIndo)
		ponte.temVindo = false
		ponte.mux.Unlock()
		ponte.c.L.Unlock()
	}
}

func (ponte *PonteOrganizada) entraPorAqui() {
	for {
		ponte.c.L.Lock()
		for ponte.direcao != "caprala" {
			ponte.c.Wait()
		}
		ponte.daquiPraLa <- "Entra por aqui"
		ponte.mux.Lock()
		ponte.temIndo = true
		ponte.mux.Unlock()
	}
}

func (ponte *PonteOrganizada) saiPorLa() {
	for {
		carro := <-ponte.daquiPraLa
		ponte.mux.Lock()
		fmt.Println(carro, "e sai por la. batida?", ponte.temVindo)
		ponte.temIndo = false
		ponte.mux.Unlock()
		ponte.c.L.Unlock()
	}
}

func (ponte *PonteOrganizada) semaforo() {
	for {
		time.Sleep(time.Millisecond)
		ponte.c.L.Lock()
		if ponte.direcao == "lapraca" {
			ponte.direcao = "caprala"
		} else {
			ponte.direcao = "lapraca"
		}
		ponte.c.L.Unlock()
		ponte.c.Broadcast()
	}
}

func main() {
	var mux sync.Mutex
	ponte := PonteOrganizada{
		deLaPraCa:  make(chan string),
		daquiPraLa: make(chan string),
		direcao:    "lapraca",
		c:          sync.NewCond(&mux),
	}
	go ponte.entraPorAqui()
	go ponte.saiPorAqui()
	go ponte.entraPorLa()
	go ponte.saiPorLa()
	go ponte.semaforo()

	time.Sleep(time.Second)
	fmt.Println("fim")
}
