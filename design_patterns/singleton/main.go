package main

import (
	"fmt"
	"sync"
	"time"
)

type Database struct{}

//CreateSingleConnection emula conexión a base de datos
func (Database) CreateSingleConnection() {
	fmt.Println("Creating Singleton fo Database")
	time.Sleep(2 * time.Second)
	fmt.Println("Creation Done")
}

var db *Database
var lock sync.Mutex

func getDatabaseInstance() *Database {
	lock.Lock() // Lock para evitar condición de carrera
	defer lock.Unlock()
	if db == nil {
		fmt.Println("Creating DB Connection")
		db = &Database{}
		db.CreateSingleConnection()
	} else {
		fmt.Println("DB Already Created")
	}
	return db
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getDatabaseInstance()
		}()
	}
	wg.Wait()
}
