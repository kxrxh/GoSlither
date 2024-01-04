package game

import "math/rand"

func PlaceNewApple(snake Snake) Point {
	newApple := Point{}
	// Searching place for new apple
	for {
		newApple.x = rand.Intn(width)
		newApple.y = rand.Intn(height)
		if !contains(snake.body, newApple) {
			return newApple
		}
	}
}
