package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Готовые конфигурации (название -> JSON payload)
	presets := []struct {
		name string
		req  map[string]interface{}
	}{
		{"blinker", map[string]interface{}{"size": 10, "live": [][]int{{4, 5}, {5, 5}, {6, 5}}, "steps": 5, "file": "result_blinker.txt"}},
		{"glider", map[string]interface{}{"size": 20, "live": [][]int{{1, 2}, {2, 3}, {3, 1}, {3, 2}, {3, 3}}, "steps": 20, "file": "result_glider.txt"}},
		{"block", map[string]interface{}{"size": 6, "live": [][]int{{1, 1}, {1, 2}, {2, 1}, {2, 2}}, "steps": 10, "file": "result_block.txt"}},
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		// Печатаем меню
		fmt.Println("Выберите конфигурацию для отправки (введите номер):")
		for i, p := range presets {
			// Показать индекс и урезанный JSON-пример
			b, _ := json.Marshal(p.req)
			fmt.Printf("%d) %s — %s\n", i, p.name, string(b))
		}
		fmt.Println("q) выход")
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("read error: %v\n", err)
			continue
		}
		line = strings.TrimSpace(line)
		if line == "q" || line == "Q" {
			fmt.Println("exit")
			return
		}
		idx, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("Введите номер из списка или 'q' для выхода")
			continue
		}
		if idx < 0 || idx >= len(presets) {
			fmt.Println("Неверный индекс")
			continue
		}

		// Сформировать JSON и отправить
		payload, _ := json.Marshal(presets[idx].req)
		if err := sendAndWait(payload); err != nil {
			fmt.Printf("send error: %v\n", err)
		}
	}
}

// sendAndWait открывает соединение, отправляет payload, закрывает сторону записи
// и ждёт ответа сервера, печатая его.
func sendAndWait(data []byte) error {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Write(data); err != nil {
		return err
	}
	// Закрываем запись, чтобы сервер получил EOF и ответил
	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.CloseWrite()
	}

	resp, _ := io.ReadAll(conn)
	fmt.Printf("server response: %s", string(resp))
	fmt.Println(" — request sent and processed")
	return nil
}
