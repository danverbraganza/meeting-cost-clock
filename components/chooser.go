package components

import (
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
		c.Duration = 60 * time.Minute
	}
	return c
}

func (c *Chooser) Cost() timefuncs.Amount {
	return timefuncs.Amount(c.Duration.Seconds()) * timefuncs.CostPerSecond(c.Selections)
}

func (*Chooser) View(x moria.Controller) moria.View {
	c := x.(*Chooser)
	return m("div#wrapper", nil,
		m("h1", nil, moria.S("How much will this meeting cost?")),

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
		}, moria.S("Start the meeting")),

		m("div#display", nil,
			m("label.copy[for='totalTime']", nil, moria.S("Meeting length")),
			m("br", nil),
			moria.S("Enter the meeting length in the following format: 1h20m30s is 1 hour, 20 minutes and 30 seconds."),
			m("br", nil),
			m("input#totalTime", js.M{
				"onchange": mithril.WithAttr("value", func(value string) {
					if duration, err := time.ParseDuration(value); err == nil {
						c.Duration = duration
					}
				}),
				"value": c.Duration.String(),
			}),
			m("br", nil),
			m("label.copy.costIntro[for='totalCost']", nil, moria.S("Estimated cost")),
			m("br", nil),
			m("span.cost.money#totalCost", nil, moria.S(c.Cost().String()))),

		m("div.copy#peopleIntro", nil, moria.S("Pick attendees")),
		m("span.subHeading", nil, moria.S("Enter the number of attendees based on their approximate salary.")),
		m("table.chooser", nil,
			m("tr", nil,
				m("th", nil, moria.S("Salary")),
				m("th", nil, moria.S("Count")),
			),
			moria.F(func(children *[]moria.View) {
				for i, cost := range timefuncs.Costs() {
					i := i // Create a copy to escape.
					*children = append(*children, m("tr.person", nil,
						m("td.money.salary", nil,
							moria.S(cost.Display)),
						m("td.rightCell", nil,
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
