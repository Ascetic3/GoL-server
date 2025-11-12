package main

// Clockl является TCP-сервером, периодически выводящим время,
import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	strings := make(chan string)
	stopChan := make(chan interface{})
	go func() {
// for x := range <-strings {
// 	fmt.Println("String:", x)
// }

		for {
			x, ok := <-strings
			if !ok {
				break // Канал закрыт и опустошен
			}
			fmt.Println("String:", x)
		}

	}()

	go listener()

	go func() {

		time.Sleep(30 * time.Second)
		stopChan <- stopChan

	}()
	

	<-stopChan
}

func listener() {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // Например, обрыв соединения
			continue
		}
		go handleConn(conn, strings) // Обработка единственного подключения
	}
}
func handleConn(c net.Conn, strings chan string) {
	defer c.Close()
	err := c.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Println("SetReadDeadline failed:", err)

		return
	}

	recvBuf := make([]byte, 1024) //io.ReadAll

	_, err = c.Read(recvBuf[:]) // recv data
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			log.Println("read timeout:", err)
			// time out
		} else {
			if err == io.EOF {
				//return
			}
			log.Println("read error: ", err)
			// some error else, do something else, for example create new conn

		}
		return
	}

	strings <- string(recvBuf) // канал небуферезированный, значение не будет отправлено до тех пор, пока оно не разблокируются.
}
