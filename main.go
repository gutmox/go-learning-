package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type HelloWorld struct {
	value string
}

type Greet interface {
	Welcome() string
}

func (h *HelloWorld) Welcome() string {
	return "Hello " + h.value
}

var wg sync.WaitGroup

func main() {

	wg.Add(2)
	resultMin := make(chan int, 1)
	resultMax := make(chan int, 1)

	helloWorld := &HelloWorld{
		value: "Pablo",
	}

	fmt.Println(helloWorld.Welcome())

	dList := []int{59, 2, 3, 6, 7, 8, 90, 18, 19}

	go minOf(resultMin, dList...)
	go maxOf(resultMax, dList...)

	wg.Wait()

	fmt.Println(<-resultMin)
	fmt.Println(<-resultMax)

	rand.Seed(time.Now().UnixNano())
	if cluster, err := loadBalancer("/accounts/subscriptions"); err == nil {
		fmt.Println(cluster)
	} else {
		fmt.Println(err)
	}
}

func minOf(result chan int, vars ...int) {
	defer wg.Done()
	min := vars[0]
	for _, i := range vars {
		if min > i {
			min = i
		}
	}
	result <- min
	return
}

func maxOf(result chan int, vars ...int) {
	defer wg.Done()
	max := vars[0]
	for _, i := range vars {
		if max < i {
			max = i
		}
	}
	result <- max
	return
}

var router = map[string]string{
	"node0": "127.0.0.1",
	"node1": "localhost",
	"node2": "0.0.0.0",
}

func loadBalancer(uri string) (url string, err error) {

	index := rand.Intn(len(router))

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			url = "none"
			err = errors.New("Cluster is down")
		}
	}()

	if index == 0 {
		panic(fmt.Sprintf("Panicing %v", index))
	}

	url = fmt.Sprintf("http://%s/%s", router[fmt.Sprintf("node%d", index)], uri)
	return
}
