package game

type Direction rune

const (
	Up    Direction = 'w'
	Down  Direction = 's'
	Left  Direction = 'a'
	Right Direction = 'd'
)

var directions = map[Direction]Point{
	'w': {0, -1},
	's': {0, 1},
	'a': {-1, 0},
	'd': {1, 0},
}

func (r Direction) Opposite() Direction {
	switch r {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	}
	return r
}
