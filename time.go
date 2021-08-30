package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	timeT, err := time.Parse("Mon Jan 02 2006 15:04:05 MST-0700", "Fri Aug 27 2021 13:05:20 GMT+0530")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(timeT, timeT.UTC())
}
