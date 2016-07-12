package main

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

var (
	m          = moria.M
	fps30      = time.Tick(time.Second / 30)
	Selections = map[int]int{}
	Costs      []float64
)

type Chooser struct {
	Duration time.Duration
}

func init() {
	var currentAmount, currentDiff float64 = 20000, 5000
	for i := 0; i < 23; i++ {
		currentAmount += currentDiff * float64(i/10+1)
		Costs = append(Costs, currentAmount)
	}
}

func (c *Chooser) Controller() moria.Controller {
	*c = Chooser{}
	c.Duration, _ = time.ParseDuration("1h")

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

// Cost returns the cost per second.
func CostPerSecond() (cumulative float64) {
	for i, cost := range Costs {
		cumulative += float64(Selections[i]) * cost
	}
	return cumulative / (2000 * 60 * 60)
}

// TODO(danver): Use a Controller PER tier.
func (c *Chooser) View(x moria.Controller) moria.View {

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
				strconv.FormatFloat(
					c.Duration.Seconds()*CostPerSecond(),
					'f', 2, 64,
				),
			)),
		),
		m("button#start.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				mithril.RouteRedirect(
					"/clock/"+c.Duration.String(),
					js.M{},
					false,
				)
			},
		}, moria.S("Start")),

		m("div.copy#peopleIntro", nil, moria.S("Select the number of attendees:")),
		moria.F(func(children *[]moria.View) {
			for i, cost := range Costs {
				i := i // Create a copy to escape.
				*children = append(*children, m("div.person", nil,
					m("div.money.salary", nil,
						moria.S(strconv.FormatFloat(cost, 'f', 0, 64))),
					m("br", nil),
					m("button.minus", js.M{"onclick": func() {
						if Selections[i] > 0 {
							Selections[i]--
						}
					}}, moria.S("\U0001F6B6\u20E0")),
					m("button.plus", js.M{"onclick": func() { Selections[i]++ }}, moria.S("\U0001F6B6\U0001F6B6")),
					m("div.count", nil,
						moria.S(strconv.Itoa(Selections[i]))),
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
	duration := mithril.RouteParam("duration").(string)
	c.left, _ = time.ParseDuration(duration)
	c.Start()
	return c
}

func (c *Clock) Start() {
	c.Lock()
	defer c.Unlock()
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
	c.Lock()
	defer c.Unlock()
	c.running = false
}

func (c *Clock) View(ctrl moria.Controller) moria.View {
	styleRed := js.M{}
	if c.left.Seconds() < 0 {
		styleRed["style"] = "color:darkred;"
	}

	pauseSigil := moria.S("\u23F8")
	if !c.running {
		pauseSigil = moria.S("\u25B6")
	}

	return m("div#wrapper", nil,
		m("h1", nil, moria.S("How much will this meeting cost?")),
		m("div#display", nil,
			m("label.copy[for='totalTime']", nil, moria.S("LENGTH:")),
			m("input#totalTime", js.M{
				"value": FormatDuration(c.left),
				"style": styleRed["style"],
			}),
			m("hr", nil),
			m("div.copy.costIntro", nil, moria.S("COST:")),
			m("div.cost.money", styleRed,
				moria.S(strconv.FormatFloat(c.left.Seconds()*CostPerSecond(),
					'f', 2, 64)),
			)),
		m("button#pause.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				if c.running {
					c.Stop()
				} else {
					c.Start()
				}
			},
		}, pauseSigil),
		m("button#stop.control", js.M{
			"config": mithril.RouteConfig,
			"onclick": func() {
				c.Stop()
				mithril.RouteRedirect(
					"/",
					js.M{},
					false,
				)
			},
		},
			moria.S("\u25a0")))
}

func main() {
	myComponent := &Chooser{}
	myClock := &Clock{}

	moria.Route(
		dom.GetWindow().Document().QuerySelector("body"), "/",
		map[string]moria.Component{
			"/":                myComponent,
			"/clock/:duration": myClock,
		})
}
