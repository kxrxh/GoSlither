package game

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kxrxh/goslither/terminal"
)

const (
	width  = 20
	height = 20
)

var (
	fpsTicker *time.Ticker
	gameState GameState
)

type GameState struct {
	score      int
	apple      Point
	snake      Snake
	isFinished chan bool
	gameSpeed  int // fps
}

func (gs *GameState) updateApplePosition() {
	gs.apple = PlaceNewApple(gs.snake)
}

func StartNew() {
	gameState = GameState{snake: NewSnake(width/2, height/2), isFinished: make(chan bool), gameSpeed: scoreToSpeed(0)}
	gameState.updateApplePosition()

	// Running game loop
	inputCh := make(chan rune)
	fpsTicker = time.NewTicker(time.Second / time.Duration(gameState.gameSpeed))

	defer fpsTicker.Stop()
	defer close(inputCh)

	go terminal.ReadInput(inputCh, gameState.isFinished)

	for range fpsTicker.C {
		select {
		case char := <-inputCh:
			// Check if the pressed key is a valid direction
			if _, ok := directions[Direction(char)]; ok {
				update(Direction(char))
			} else if char == 'q' {
				gameState.isFinished <- true
				fmt.Println("Shh... bye...")
				return
			} else {
				// Bad input received, continue processing
				isEnd := update(gameState.snake.direction)
				if isEnd {
					gameState.isFinished <- true
					fmt.Println("Game over. Your score: ", gameState.score)
					return
				}
			}
		default:
			if isEnd := update(gameState.snake.direction); isEnd {
				gameState.isFinished <- true
				fmt.Println("Game over. Your score: ", gameState.score)
				return
			}
		}
		draw()
	}
}

func update(dir Direction) bool {
	if dir != gameState.snake.direction.Opposite() {
		gameState.snake.direction = dir
	}

	newHead := gameState.snake.GetNextHeadPos()

	// Check if the new head collides with the snake's body or goes out of bounds
	if contains(gameState.snake.body[1:], newHead) {
		return true
	}

	gameState.snake.body = append([]Point{newHead}, gameState.snake.body...)

	// Check if the snake eats the apple
	if newHead == gameState.apple {
		gameState.updateApplePosition()
		gameState.score += 10
		if gameState.score%50 == 0 {
			gameState.gameSpeed = scoreToSpeed(gameState.score)
			fpsTicker.Reset(time.Second / time.Duration(gameState.gameSpeed))
		}
	} else {
		gameState.snake.body = gameState.snake.body[:len(gameState.snake.body)-1]
	}

	return false
}
func draw() {
	terminal.ClearScreen()

	// Buffer to store the entire game screen
	var screenBuffer strings.Builder

	// Draw top border
	screenBuffer.WriteString("┌")
	for i := 0; i < width; i++ {
		screenBuffer.WriteString("──")
	}
	screenBuffer.WriteString("\033[1D┐\n")

	// Draw game area
	for y := 0; y < height; y++ {
		screenBuffer.WriteString("│")
		for x := 0; x < width; x++ {
			p := Point{x, y}
			cell := "  "

			if p == gameState.snake.body[0] {
				cell = "@ "
			} else if contains(gameState.snake.body[1:], p) {
				cell = "o "
			} else if p == gameState.apple {
				cell = "A "
			}

			screenBuffer.WriteString(cell)
		}
		screenBuffer.WriteString("\033[1D│\n")
	}

	// Draw bottom border
	screenBuffer.WriteString("└")
	for i := 0; i < width; i++ {
		screenBuffer.WriteString("──")
	}
	screenBuffer.WriteString("\033[1D┘\n")

	// Draw score
	screenBuffer.WriteString("Score: ")
	screenBuffer.WriteString(strconv.Itoa(gameState.score))
	screenBuffer.WriteString("\n")

	screenBuffer.WriteString("Tip: Use w, a, s, d to move. Press 'q' to quit.\n")

	// Print the entire screen at once
	fmt.Print(screenBuffer.String())
}
