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
		m("table#display", nil,
			m("tr", nil,
				m("td", nil, m("label.copy[for='totalTime']", nil, moria.S("Length"))),
				m("td", nil, m("input#totalTime", js.M{
					"onchange": mithril.WithAttr("value", func(value string) {
						if duration, err := time.ParseDuration(value); err == nil {
							c.Duration = duration
						}
					}),
					"value": c.Duration.String(),
				})),
			),
			m("tr", nil,
				m("td", nil,
					m("label.copy.costIntro", nil, moria.S("Cost"))),
				m("td", nil, m("span.cost.money", nil, moria.S(
					c.Cost().String()))),
			),
		),
		m("button#start.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {

				mithril.RouteRedirect(
					strings.Join([]string{
						"",
						"clock",
						c.Duration.String(),
						c.Cost().FloatStr(),
					},
						"/"),
					js.M{},
					false,
				)
			},
		}, moria.S("Start")),

		m("div.copy#peopleIntro", nil, moria.S("Add or remove paid attendees")),
		m("table.chooser", nil,
			m("tr", nil,
				m("th[tablewidth='1']", nil, moria.S("Salary")),
				m("th[tablewidth='1']", nil, moria.S("Count")),
			),
			moria.F(func(children *[]moria.View) {
				for i, cost := range timefuncs.Costs() {
					i := i // Create a copy to escape.
					*children = append(*children, m("tr.person", nil,
						m("td.money.salary", nil,
							moria.S(cost.Display)),
						m("td.count", nil,
							m("input[type='number'][min='0'].count", js.M{
								"value": c.Selections[i],
								"onchange": mithril.WithAttr("value", func(value string) {
									if intValue, err := strconv.Atoi(value); err == nil {
										c.Selections[i] = intValue
									}
								})}))))
				}
			})))
}
