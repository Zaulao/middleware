package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// temVindo e temIndo estao aqui apenas para proposito de visualizacao.
// serve para checar que nao tem um carro passando de la pra ca e um
// daqui pra la ao mesmo tempo. o mux eh apenas para evitar race condition
// no manuseio dessas variaveis ajudantes, o problema seria resolvido igual
// sem temIndo, temVindo e o mux que os protege.

// Esta versao controla o acesso a ponte, batidas nao devem acontecer.
type PonteOrganizada struct {
	deLaPraCa  chan string
	daquiPraLa chan string
	temIndo    bool
	temVindo   bool
	sem        *semaphore.Weighted
	ctx        context.Context
	mux        sync.Mutex
}

func (ponte *PonteOrganizada) entraPorLa() {
	for {
		ponte.sem.Acquire(ponte.ctx, 1)
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
		ponte.sem.Release(1)
	}
}

func (ponte *PonteOrganizada) entraPorAqui() {
	for {
		ponte.sem.Acquire(ponte.ctx, 1)
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
		ponte.sem.Release(1)
	}
}

func main() {
	ponte := PonteOrganizada{
		deLaPraCa:  make(chan string),
		daquiPraLa: make(chan string),
		ctx:        context.TODO(),
		sem:        semaphore.NewWeighted(1),
	}
	go ponte.entraPorAqui()
	go ponte.saiPorAqui()
	go ponte.entraPorLa()
	go ponte.saiPorLa()

	time.Sleep(time.Millisecond)
	fmt.Println("fim")
}
