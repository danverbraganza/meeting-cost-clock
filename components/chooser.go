package components

import (
	"fmt"
	"math"
	"meeting-cost-clock/timefuncs"
	"strconv"
	"strings"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"
)

type Chooser struct {
	Duration   time.Duration
	Selections map[int]int
}

func (c *Chooser) Controller() moria.Controller {
	if c.Selections == nil {
		c.Selections = map[int]int{}
	}
	if c.Duration == 0*time.Second {
		c.Duration, _ = time.ParseDuration("1h")
	}
	return c
}

func FormatDuration(d time.Duration) string {
	rounder := math.Floor
	if math.Signbit(d.Hours()) {
		rounder = math.Ceil
	}

	return fmt.Sprintf("%02.0f:%02d:%02d:%03d",
		rounder(d.Hours()),
		int(math.Abs(d.Minutes()))%60,
		int(math.Abs(d.Seconds()))%60,
		int(math.Abs(float64(d.Nanoseconds()/1e6)))%1000)
}

func (c *Chooser) Cost() timefuncs.Amount {
	return timefuncs.Amount(c.Duration.Seconds()) * timefuncs.CostPerSecond(c.Selections)
}

func (*Chooser) View(x moria.Controller) moria.View {
	c := x.(*Chooser)
	return m("div#wrapper", nil,
		m("h1", nil, moria.S("How much will this meeting cost?")),
		m("div#display", nil,
			m("label.copy[for='totalTime']", nil, moria.S("LENGTH:")),
			m("input#totalTime", js.M{
				"onchange": mithril.WithAttr("value", func(value string) {
					if duration, err := time.ParseDuration(value); err == nil {
						c.Duration = duration
					}
				}),
				"value": FormatDuration(c.Duration),
			}),
			m("hr", nil),
			m("div.copy.costIntro", nil, moria.S("COST:")),
			m("div.cost.money", nil, moria.S(
				c.Cost().String(),
			)),
		),
		m("button#start.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {

				fmt.Println(strings.Join([]string{
					"",
					"clock",
					c.Duration.String(),
					c.Cost().String()},
					"/"))

				mithril.RouteRedirect(
					strings.Join([]string{
						"",
						"clock",
						c.Duration.String(),
						c.Cost().String(),
					},
						"/"),
					js.M{},
					false,
				)
			},
		}, moria.S("Start")),

		m("div.copy#peopleIntro", nil, moria.S("Select the number of attendees:")),
		moria.F(func(children *[]moria.View) {
			for i, cost := range timefuncs.Costs() {
				i := i // Create a copy to escape.
				*children = append(*children, m("div.person", nil,
					m("div.money.salary", nil,
						moria.S(cost.Display)),
					m("br", nil),
					m("button.minus", js.M{"onclick": func() {
						if c.Selections[i] > 0 {
							c.Selections[i]--
						}
					}}, moria.S("\U0001F6B6\u20E0")),
					m("button.plus", js.M{"onclick": func() { c.Selections[i]++ }}, moria.S("\U0001F6B6\U0001F6B6")),
					m("div.count", nil,
						moria.S(strconv.Itoa(c.Selections[i]))),
				))
			}
		}))
}
