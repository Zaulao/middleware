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

// Note que esta versao NAO sincroniza o acesso a ponte -- batidas irao acontecer
type PonteDesorganizada struct {
	deLaPraCa  chan string
	daquiPraLa chan string
	temVindo   bool
	temIndo    bool
	mux        sync.Mutex
}

func (ponte *PonteDesorganizada) entraPorLa() {
	for {
		ponte.deLaPraCa <- "Entra por la"
		ponte.mux.Lock()
		ponte.temVindo = true
		ponte.mux.Unlock()
	}
}

func (ponte *PonteDesorganizada) saiPorAqui() {
	for {
		carro := <-ponte.deLaPraCa
		ponte.mux.Lock()
		ponte.temVindo = false
		fmt.Println(carro, "e sai por aqui. batida?", ponte.temIndo)
		ponte.mux.Unlock()
	}
}

func (ponte *PonteDesorganizada) entraPorAqui() {
	for {
		ponte.daquiPraLa <- "Entra por aqui"
		ponte.mux.Lock()
		ponte.temIndo = true
		ponte.mux.Unlock()
	}
}

func (ponte *PonteDesorganizada) saiPorLa() {
	for {
		carro := <-ponte.daquiPraLa
		ponte.mux.Lock()
		ponte.temIndo = false
		fmt.Println(carro, "e sai por la. batida?", ponte.temVindo)
		ponte.mux.Unlock()
	}
}

func main() {
	ponte := PonteDesorganizada{
		deLaPraCa:  make(chan string),
		daquiPraLa: make(chan string),
	}
	go ponte.entraPorAqui()
	go ponte.saiPorAqui()
	go ponte.entraPorLa()
	go ponte.saiPorLa()
	time.Sleep(time.Millisecond)
	fmt.Println("fim")
}
