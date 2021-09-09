package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

type Memory struct {
	f     Function
	cache map[int]FunctionResult
	lock  sync.Mutex
}

type Function func(key int) (interface{}, error)

type FunctionResult struct {
	value interface{}
	err   error
}

func NewCache(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]FunctionResult),
	}
}

func (m *Memory) Get(key int) (interface{}, error) {
	m.lock.Lock()
	result, exist := m.cache[key]
	m.lock.Unlock()
	if !exist {
		m.lock.Lock()
		result.value, result.err = m.f(key)
		m.cache[key] = result
		m.lock.Unlock()
	}
	return result.value, result.err
}

func GetFibonacci(n int) (interface{}, error) { //Retornar interface{} quiere decir que puede devolver cualquier tipo de valor
	return Fibonacci(n), nil
}

func main() {
	cache := NewCache(GetFibonacci)
	fibo := []int{42, 41, 38, 40, 38, 42, 42, 42, 42, 42, 42}
	var wg sync.WaitGroup
	maxGoRoutines := 2
	channel := make(chan int, maxGoRoutines)
	startTotal := time.Now()
	for _, n := range fibo {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			channel <- 1
			start := time.Now()
			value, err := cache.Get(index)
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("%d, %s, %d\n", index, time.Since(start), value)
			<-channel
		}(n)
	}
	wg.Wait()
	fmt.Printf("Tiempo total: %s\n", time.Since(startTotal))
}
