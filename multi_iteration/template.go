package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type array2D [][]int

func newGameOfLife(filename string) array2D {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("error чтения файла", err)
		os.Exit(1)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	size, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	grid := make(array2D, size)
	for i := 0; i < size; i++ {
		line := lines[i+1]
		cells := strings.Fields(line)
		grid[i] = make([]int, size)
		for j := 0; j < size; j++ {
			val, _ := strconv.Atoi(cells[j])
			grid[i][j] = val
		}
	}
	return grid
}

func save(grid array2D, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
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
}

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
					newGrid[i][j] = 1 // Выживает
				} else {
					newGrid[i][j] = 0 // Умирает
				}
			} else {
				// Мертвая клетка
				if neighbors == 3 {
					newGrid[i][j] = 1 // Оживает
				}
				// Иначе остается мертвой (уже 0 по умолчанию)
			}
		}
	}

	return newGrid
}

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
} // periodic boundary condition

// НОВАЯ ВЕРСИЯ: несколько итераций с анимацией
func main() {
	// Количество итераций по умолчанию
	iterations := 20
	patternFile := "state.txt"
	// Если указан аргумент командной строки, используем его
	if len(os.Args) > 1 {
		if num, err := strconv.Atoi(os.Args[1]); err == nil && num > 0 {
			iterations = num
		} else {
			fmt.Println("Ошибка: количество итераций должно быть положительным числом")
			fmt.Println("Использование: go run template.go [количество_итераций]")
			fmt.Printf("Используется значение по умолчанию: %d\n", iterations)
		}
		if len(os.Args) > 2 {
			patternFile = "patterns/" + os.Args[2] + ".txt"
		}
	}

	// Загружаем начальное состояние
	grid := newGameOfLife(patternFile)

	// Папка для сохранения состояний
	statesDir := "states"
	if err := os.MkdirAll(statesDir, 0755); err != nil {
		fmt.Printf("Не удалось создать папку %s: %v\n", statesDir, err)
		os.Exit(1)
	}

	// Сохраняем начальное состояние
	save(grid, filepath.Join(statesDir, "state_0.txt"))
	fmt.Println("Состояние 0 сохранено")

	// Выполняем итерации
	for i := 1; i <= iterations; i++ {
		grid = simulateDay(grid)
		filename := fmt.Sprintf("state_%d.txt", i)
		save(grid, filepath.Join(statesDir, filename))
		fmt.Printf("Итерация %d сохранена\n", i)
		time.Sleep(300 * time.Millisecond) // задержка между итерациями
	}

	fmt.Printf("Симуляция завершена! Выполнено %d итераций.\n", iterations)
}
