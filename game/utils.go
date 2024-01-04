package game

// contains checks if a given point is present in a list of points.
//
// The function takes in two parameters:
// - points: a slice of Point structs representing a list of points.
// - p: a Point struct representing the point to be checked.
//
// The function returns a boolean value indicating whether the point is present in the list or not.
func contains(points []Point, p Point) bool {
	for _, v := range points {
		if v == p {
			return true
		}
	}
	return false
}
