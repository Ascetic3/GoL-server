package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

// ------------------------------------------------------------------
// Простой TCP сервер «Игра Жизнь»
// ------------------------------------------------------------------
// клиент присылает один JSON (в байтах) и закрывает соединение.
// В обработчике соединения я НЕ парсю JSON — только читаю байты и кладу
// строку в канал jobs. Одна фоновая горутина (worker) последовательно
// читает строки из канала, парсит JSON, запускает симуляцию и записывает
// файл с результатом. Архитектура: handleConn -> jobs chan -> worker.
// ------------------------------------------------------------------
// Request — ожидаемые поля JSON от клиента.
// Пример JSON:
// {"size":10,"live":[[4,5],[5,5],[6,5]],"steps":5,"file":"result.txt"}
type Request struct {
	Size  int     `json:"size"`
	Live  [][]int `json:"live"`
	Steps int     `json:"steps"`
	File  string  `json:"file"`
}

// array2D — тип сетки, использую ваш существующий формат [][]int
type array2D [][]int

// simulateDay вычисляет следующее поколение сетки по правилам Game of Life.
func simulateDay(grid array2D) array2D {
	size := len(grid)
	newGrid := make(array2D, size)
	for i := range newGrid {
		newGrid[i] = make([]int, size)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			neighbors := countNeighborsPeriodic(grid, i, j)
			currentCell := grid[i][j]

			if currentCell == 1 {
				if neighbors == 2 || neighbors == 3 {
					newGrid[i][j] = 1
				} else {
					newGrid[i][j] = 0
				}
			} else {
				if neighbors == 3 {
					newGrid[i][j] = 1
				}
			}
		}
	}
	return newGrid
}

// countNeighborsPeriodic считает соседей с периодическими границами

func countNeighborsPeriodic(grid array2D, row, col int) int {
	size := len(grid)
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			r := pbc(row+i, size)
			c := pbc(col+j, size)
			count += grid[r][c]
		}
	}
	return count
}

func pbc(x, l int) int {
	if x < 0 {
		return x + l
	}
	if x < l {
		return x
	}
	return x - l
}

// save сохраняет сетку в файл
func save(grid array2D, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	size := len(grid)
	fmt.Fprintf(file, "%d\n", size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Fprintf(file, "%d ", grid[i][j])
		}
		fmt.Fprintln(file)
	}
	return nil
}

// newGridFromLive строит начальную сетку из списка живых клеток.

func newGridFromLive(size int, live [][]int) array2D {
	grid := make(array2D, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]int, size)
	}
	for _, p := range live {
		if len(p) >= 2 {
			r, c := p[0], p[1]
			if r >= 0 && r < size && c >= 0 && c < size {
				grid[r][c] = 1
			}
		}
	}
	return grid
}

// ------------------------------------------------------------------
// main: запускает TCP listener, канал задач и одного фонового worker'а
// ------------------------------------------------------------------
func main() {
	// канал задач: сюда я помещаю  JSON-строки (строка = bytes->string).

	jobs := make(chan string)

	// запускаю одного worker'а, который последовательно обрабатывает задачи
	go worker(jobs)

	// запускаю TCP listener на localhost:8000
	ln, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("listen:", err)
	}
	defer ln.Close()
	fmt.Println("gol server listening on localhost:8000")

	// главный цикл: принимаю соединения и запускаю обработчик
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("connection error", err)
			continue
		}
		go handleConn(conn, jobs)
	}
}

// handleConn читает все байты из соединения (клиент должен закрыть
// соединение после отправки) и кладёт полученную строку в канал jobs.

func handleConn(conn net.Conn, jobs chan<- string) {
	defer conn.Close()

	// Ограничиваю время на чтение
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	data, err := io.ReadAll(conn)
	if err != nil {
		log.Println("read error:", err)
		return
	}
	msg := string(data)

	// отправляю строку в канал — worker  распарсит JSON
	jobs <- msg

	// подтверждаю приём клиенту
	conn.Write([]byte("accepted\n"))
}

// worker: единственная горутина, последовательно обрабатывающая задачи.
// Она парсит JSON, выполняет указанное число шагов и записывает файл.
func worker(jobs <-chan string) {
	for msg := range jobs {
		var req Request
		if err := json.Unmarshal([]byte(msg), &req); err != nil {
			log.Println("invalid json:", err)
			continue
		}
		// минимальная валидация полей запроса
		if req.Size <= 0 || req.Steps < 0 || req.File == "" {
			log.Println("invalid request fields")
			continue
		}

		// строю стартовую сетку по списку живых клеток
		grid := newGridFromLive(req.Size, req.Live)

		// прогоняю симуляцию нужное число шагов
		for i := 0; i < req.Steps; i++ {
			grid = simulateDay(grid)
		}

		// сохраняю результат в файл (перезаписываю, если есть)
		if err := save(grid, req.File); err != nil {
			log.Println("save error:", err)
		} else {
			log.Println("wrote result to", req.File)
		}
	}
}
