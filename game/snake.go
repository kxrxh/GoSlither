package game

type Snake struct {
	body      []Point
	direction Direction
}

func NewSnake(x, y int) Snake {
	return Snake{body: []Point{{x, y}}, direction: Up}
}

func (s Snake) GetDirection() Point {
	return directions[s.direction]
}

func (s Snake) GetNextHeadPos() Point {
	head := s.body[0]
	newX := head.x + s.GetDirection().x
	newY := head.y + s.GetDirection().y

	if newX < 0 {
		newX = width - 1
	}
	if newX >= width {
		newX = 0
	}

	if newY < 0 {
		newY = height - 1
	}
	if newY >= height {
		newY = 0
	}

	return Point{newX, newY}
}
