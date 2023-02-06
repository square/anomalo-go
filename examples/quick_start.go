package main

import (
	"fmt"
	"log"

	"github.com/square/anomalo-go/anomalo"
)

func main() {
	client, err := anomalo.CreateClient()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	fmt.Println(client.Ping())
}
