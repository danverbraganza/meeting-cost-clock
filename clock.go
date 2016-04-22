package main

import (
	"strconv"
	"sync"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	m     = moria.M
	fps30 = time.Tick(time.Second / 30)
)

type Chooser struct {
	Costs      []float64
	Selections map[int]int
	Duration   time.Duration
}

func (c *Chooser) Controller() moria.Controller {
	*c = Chooser{}
	var currentAmount, currentDiff float64 = 20000, 5000

	for i := 0; i < 20; i++ {
		currentAmount += currentDiff * float64(i/10+1)
		c.Costs = append(c.Costs, currentAmount)
	}
	c.Selections = map[int]int{}
	c.Duration, _ = time.ParseDuration("1h")

	return c
}

// Cost returns the cost per second.
func (c Chooser) CostPerSecond() (cumulative float64) {
	for i, cost := range c.Costs {
		cumulative += float64(c.Selections[i]) * cost
	}
	return cumulative / (2000 * 60 * 60)
}

// TODO(danver): Use a Controller PER tier.
func (c *Chooser) View(x moria.Controller) moria.View {
	// TODO(danver): Do NOT use c, but x.

	return m("div#wrapper", nil,
		m("h1", nil, moria.S("How much will this meeting cost?")),
		m("label.copy[for='totalTime']", nil, moria.S("How long is this meeting?")),
		m("br", nil),
		m("input#totalTime", js.M{
			"onchange": mithril.WithAttr("value", func(value string) {
				if duration, err := time.ParseDuration(value); err == nil {
					c.Duration = duration
				}
			}),
			"value": c.Duration.String(),
		}),
		m("a#start", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				mithril.RouteRedirect(
					"/clock",
					js.M{},
					false,
				)
			},
		}, moria.S("Start the meeting")),
		m("div.copy.costIntro", nil, moria.S("This is how much this meeting will cost you:")),
		m("div.cost.money", nil, moria.S(
			strconv.FormatFloat(
				c.Duration.Seconds()*c.CostPerSecond(),
				'f', 2, 64,
			),
		)),
		m("div.copy#peopleIntro", nil, moria.S("Select the number of attendees:")),
		moria.F(func(children *[]moria.View) {
			for i, cost := range c.Costs {
				i_ := i // Create a copy to escape.
				*children = append(*children, m("div.person", nil,
					m("div.money.salary", js.M{"onclick": func() { c.Selections[i_]++ }},

						moria.S(strconv.FormatFloat(cost, 'f', 0, 64))),
					m("div.count", js.M{"onclick": func() {
						if c.Selections[i_] > 0 {
							c.Selections[i_]--
						}
					}},
						moria.S(strconv.Itoa(c.Selections[i]))),
				))
			}
		}))
}

type Clock struct {
	sync.Mutex
	left    time.Duration
	last    time.Time
	running bool
}

func (c *Clock) Controller() moria.Controller {
	c.left, _ = time.ParseDuration("1h")
	c.Start()
	return c
}

func (c *Clock) Start() {
	defer c.Unlock()
	c.Lock()
	c.last = time.Now()
	c.running = true
	go func() {
		for c.running {
			<-fps30
			now := time.Now()
			c.left -= now.Sub(c.last)
			c.last = now
			mithril.Redraw(false)
		}
	}()
}

func (c *Clock) Stop() {
	defer c.Unlock()
	c.Lock()
	c.running = false
}

func (c *Clock) View(ctrl moria.Controller) moria.View {
	return m("div#wrapper", nil,
		m("div#clock", js.M{
			"onclick": func() {
				c.running = !c.running
			},
		},
			moria.S(c.left.String())),

		m("a#pause", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				if c.running {
					c.Stop()
				} else {
					c.Start()
				}
			},
		}, moria.S("Pause/Unpause")),

		m("a#stop", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				c.Stop()
				mithril.RouteRedirect(
					"/",
					js.M{},
					false,
				)
			},
		}, moria.S("Stop the meeting")))

}

func main() {
	myComponent := &Chooser{}
	myClock := &Clock{}

	moria.Route(
		dom.GetWindow().Document().QuerySelector("body"), "/",
		map[string]moria.Component{
			"/":      myComponent,
			"/clock": myClock,
		})
}
