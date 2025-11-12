package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

// СТАРАЯ ВЕРСИЯ: только одна итерация
func main() {
	grid := newGameOfLife("state.txt")
	grid = simulateDay(grid)
	save(grid, "state.txt")
	fmt.Println("Результат сохранен в state.txt")
}
