// Package timefuncs contains generalized functions for working with time and
// money.
package timefuncs

import "strconv"

var costs []Amount

func init() {
	var currentAmount, currentDiff Amount = 20000, 5000
	for i := 0; i < 23; i++ {
		currentAmount += currentDiff * Amount(i/10+1)
		costs = append(costs, currentAmount)
	}
}

type Amount float64

func (a Amount) String() string {
	return strconv.FormatFloat(float64(a), 'f', 2, 64)

}

// Cost returns the cost per second of a given set of selections.
func CostPerSecond(selections map[int]int) (cumulative Amount) {
	for i, cost := range costs {
		cumulative += Amount(selections[i]) * cost
	}
	return cumulative / (2000 * 60 * 60)
}

func Costs() (result []struct{ Title, Display string }) {
	for _, c := range costs {
		result = append(result,
			struct{ Title, Display string }{
				Title:   "Unused",
				Display: c.String(),
			})
	}
	return
}
