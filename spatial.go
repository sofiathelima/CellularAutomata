package main

import (
	"bufio"
	"fmt"
	"gifhelper"
	"image/png"
	"os"
	"strconv"
	"strings"
)

// The data stored in a single cell of a field
type Cell struct {
	strategy string  //represents "C" or "D" corresponding to the type of prisoner in the cell
	score    float64 //represents the score of the cell based on the prisoner's relationship with neighboring cells
}

// The game board is a 2D slice of Cell objects
type GameBoard [][]Cell

// FindNeighbors takes a Gameboard along with r, c indices. It returns a 2D list
// of ints of the r, c indices of all possible neighbors.
func FindNeighbors(board GameBoard, r, c int) [][]int {

	neighbors := make([][]int, 0)
	if r != 0 && r != (len(board)-1) {
		if c != 0 && c != (len(board[0])-1) {
			neighbors = append(neighbors, []int{r - 1, c - 1})
			neighbors = append(neighbors, []int{r - 1, c})
			neighbors = append(neighbors, []int{r - 1, c + 1})
			neighbors = append(neighbors, []int{r, c - 1})
			neighbors = append(neighbors, []int{r, c + 1})
			neighbors = append(neighbors, []int{r + 1, c - 1})
			neighbors = append(neighbors, []int{r + 1, c})
			neighbors = append(neighbors, []int{r + 1, c + 1})
		} else if c == (len(board[0]) - 1) {
			neighbors = append(neighbors, []int{r - 1, c - 1})
			neighbors = append(neighbors, []int{r - 1, c})
			neighbors = append(neighbors, []int{r, c - 1})
			neighbors = append(neighbors, []int{r + 1, c - 1})
			neighbors = append(neighbors, []int{r + 1, c})
		} else { // if c == 0
			neighbors = append(neighbors, []int{r - 1, c})
			neighbors = append(neighbors, []int{r - 1, c + 1})
			neighbors = append(neighbors, []int{r, c + 1})
			neighbors = append(neighbors, []int{r + 1, c})
			neighbors = append(neighbors, []int{r + 1, c + 1})
		}
	} else if r == (len(board) - 1) {
		if c != 0 && c != (len(board[0])-1) {
			neighbors = append(neighbors, []int{r - 1, c - 1})
			neighbors = append(neighbors, []int{r - 1, c})
			neighbors = append(neighbors, []int{r - 1, c + 1})
			neighbors = append(neighbors, []int{r, c - 1})
			neighbors = append(neighbors, []int{r, c + 1})
		} else if c == (len(board[0]) - 1) {
			neighbors = append(neighbors, []int{r - 1, c - 1})
			neighbors = append(neighbors, []int{r - 1, c})
			neighbors = append(neighbors, []int{r, c - 1})
		} else {
			neighbors = append(neighbors, []int{r - 1, c})
			neighbors = append(neighbors, []int{r - 1, c + 1})
			neighbors = append(neighbors, []int{r, c + 1})
		}
	} else { // if r == 0
		if c != 0 && c != (len(board[0])-1) {
			neighbors = append(neighbors, []int{r, c - 1})
			neighbors = append(neighbors, []int{r, c + 1})
			neighbors = append(neighbors, []int{r + 1, c - 1})
			neighbors = append(neighbors, []int{r + 1, c})
			neighbors = append(neighbors, []int{r + 1, c + 1})
		} else if c == (len(board[0]) - 1) {
			neighbors = append(neighbors, []int{r, c - 1})
			neighbors = append(neighbors, []int{r + 1, c - 1})
			neighbors = append(neighbors, []int{r + 1, c})
		} else {
			neighbors = append(neighbors, []int{r, c + 1})
			neighbors = append(neighbors, []int{r + 1, c})
			neighbors = append(neighbors, []int{r + 1, c + 1})
		}
	}

	return neighbors
}

//UpdateCellStrategy takes 2 GameBoards (one with old scores and one with new scores)
// along with r, c indices and returns the new strategy for that particular Cell
// based on the neighbor with the best score.
func UpdateCellStrategy(board, newBoard GameBoard, r, c int) string {
	listNeighbors := FindNeighbors(board, r, c)

	strategy := board[r][c].strategy

	max := newBoard[r][c].score

	for i := range listNeighbors {
		rNeighbor := listNeighbors[i][0]
		cNeighbor := listNeighbors[i][1]

		if newBoard[rNeighbor][cNeighbor].score > max {
			max = newBoard[rNeighbor][cNeighbor].score
			strategy = board[rNeighbor][cNeighbor].strategy
		}
	}
	return strategy
}

// UpdateBoardStrategies takes 2 Gameboards (one with old scores and one with new scores)
// and returns the newBoard with updated the strategy for each cell.
func UpdateBoardStrategies(board, newBoard GameBoard) GameBoard {
	numRows := len(board)
	numCols := len(board[0])

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			newBoard[r][c].strategy = UpdateCellStrategy(board, newBoard, r, c)
		}
	}
	return newBoard
}

// UpdateCell takes 2 Gameboards (one with the old scores and one with new scores)
// along with r, c indices and calculates the new score for that particular Cell
// based on the interactions with its neighbors. It returns the updated Cell in the newBoard.
func UpdateCellScore(board, newBoard GameBoard, r, c int, b float64) Cell {

	listNeighbors := FindNeighbors(board, r, c)

	for i := range listNeighbors {
		rNeighbor := listNeighbors[i][0]
		cNeighbor := listNeighbors[i][1]

		if board[r][c].strategy == "C" && board[rNeighbor][cNeighbor].strategy == "C" {
			newBoard[r][c].score += 1.0
			newBoard[rNeighbor][cNeighbor].score += 1.0
		} else if board[r][c].strategy == "C" && board[rNeighbor][cNeighbor].strategy == "D" {
			newBoard[rNeighbor][cNeighbor].score += b
		} else if board[r][c].strategy == "D" && board[rNeighbor][cNeighbor].strategy == "C" {
			newBoard[r][c].score += b
		}
	}

	return newBoard[r][c]
}

// UpdateBoardScores takes a Gameboard and the b parameter.
// It creates a GameBoard with all the current strategies and updates the scores for each Cell.
func UpdateBoardScores(board GameBoard, b float64) GameBoard {
	numRows := len(board)
	numCols := len(board[0])
	newBoard := EmptyBoard(numRows, numCols)

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			newBoard[r][c].strategy = board[r][c].strategy
		}
	}

	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			newBoard[r][c] = UpdateCellScore(board, newBoard, r, c, b)
		}
	}

	return newBoard
}

// PlayTournament takes the initial Gameboard, b, and numGens as input
// and generates a list of boards resulting from each round in the tournament.
func PlayTournament(initialBoard GameBoard, b float64, numGens int) []GameBoard {
	boards := make([]GameBoard, numGens+1)
	boards[0] = initialBoard

	for i := 0; i < numGens; i++ {
		boards[i+1] = UpdateBoardStrategies(boards[i], UpdateBoardScores(boards[i], b))
	}

	return boards
}

// EmptyBoard takes numRows and numCols as inputs and returns a gameboard with
// appropriate number of rows and colums, where all values = 0.
func EmptyBoard(numRows, numCols int) GameBoard {
	var board GameBoard
	board = make(GameBoard, numRows)
	for r := range board {
		board[r] = make([]Cell, numCols)
	}

	return board
}

// ReadFileAndInitializeBoard reads the filename input and scans the file
// and returns a GameBoard with initialized strategies and 0.0 scores.
func ReadFileAndInitializeBoard(filename string) GameBoard {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: something went wrong opening the file.")
		fmt.Println("Probably you gave the wrong filename.")
	}

	defer file.Close()

	// scan text file and create nBoard with nRows and nCols in the first line
	scanner := bufio.NewScanner(file)
	nBoard := make([]string, 0)
	for scanner.Scan() {
		currentLine := scanner.Text()
		nBoard = append(nBoard, currentLine)
	}
	if scanner.Err() != nil {
		panic("Error: There was a problem reading the file")
	}

	// read nRows and nCols from first line of file
	firstLine := strings.Split(nBoard[0], " ")
	nRow, err2 := strconv.Atoi(firstLine[0])
	if err2 != nil {
		panic(err2)
	}
	nCol, err3 := strconv.Atoi(firstLine[1])
	if err3 != nil {
		panic(err3)
	}

	// read characters from strings in nBoard and initialize board
	board := EmptyBoard(nRow, nCol)
	for i := 0; i < nRow; i++ {
		line := nBoard[i+1]

		for j := range line {
			char := line[j : j+1]
			if char == "C" {
				board[i][j].strategy = "C"
			} else {
				board[i][j].strategy = "D"
			}
		}
	}

	return board
}

func main() {
	filename := os.Args[1]
	b, _ := strconv.ParseFloat(os.Args[2], 64)
	numGens, _ := strconv.Atoi(os.Args[3])

	initialBoard := ReadFileAndInitializeBoard(filename)

	fmt.Println("Playing the tournament.")

	boards := PlayTournament(initialBoard, float64(b), numGens)

	fmt.Println("Automaton played. Now, drawing images.")

	// we need a slice of image objects
	imglist := DrawGameBoards(boards)
	fmt.Println("Boards drawn to images! Now, convert to animated GIF.")

	// convert images to a GIF
	gifhelper.ImagesToGIF(imglist, "out")

	fmt.Println("Success! GIF produced.")

	f, _ := os.Create("Prisoners.png")
	png.Encode(f, imglist[len(imglist)-1])
}
