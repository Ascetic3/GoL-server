package main

// Netcatl - TCP-клиент только для чтения,

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
	for {
		fmt.Println("Enter the data")
		number := 0
		_, err := fmt.Scanf("%d\n", &number)
		if err != nil {
			fmt.Printf("Failed get number from console: %s\n", err)
			continue
		}
		data, err := json.Marshal(number)
		if err != nil {
			fmt.Printf("Failed to marshal number: %s\n", err)
			continue
		}
		sendDataToTcp(data)
	}
}

func sendDataToTcp(data []byte) error {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	n, err := conn.Write(data)
	fmt.Printf("%d byted writed", n)
	return err
}
