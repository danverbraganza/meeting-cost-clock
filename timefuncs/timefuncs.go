// Package timefuncs contains generalized functions for working with time and
// money.
package timefuncs

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

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
	result := &bytes.Buffer{}
	amt := float64(a)
	groups := []string{}

	round := math.Floor
	if amt < 0 {
		round = math.Ceil
		io.WriteString(result, "-")
	}

	whole := round(amt)
	frac := math.Abs(amt - whole)
	fracStr := fmt.Sprintf("%02.0f", frac*100)

	wholeStr := strconv.FormatInt(int64(amt), 10)
	idx := len(wholeStr)
	for idx > 0 {
		endx := idx - 3
		if endx < 0 {
			endx = 0
		}
		groups = append(groups, wholeStr[endx:idx])
		idx -= 3
	}

	l := len(groups) - 1
	for i := 0; i < len(groups)/2; i++ {
		groups[i], groups[l-i] = groups[l-i], groups[i]
	}

	io.WriteString(result, strings.Join(groups, ","))

	io.WriteString(result, ".")

	io.WriteString(result, fracStr)

	return result.String()
}

func (a Amount) FloatStr() string {
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
