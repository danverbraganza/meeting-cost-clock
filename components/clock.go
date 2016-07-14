package components

import (
	"fmt"
	"meeting-cost-clock/timefuncs"
	"strconv"
	"sync"
	"time"

	"github.com/danverbraganza/go-mithril"
	"github.com/danverbraganza/go-mithril/moria"
	"github.com/gopherjs/gopherjs/js"
)

var (
	m     = moria.M
	fps30 = time.Tick(time.Second / 30)
)

type Clock struct {
	sync.Mutex
	totalTime, timeSpent time.Duration
	last                 time.Time
	totalCost            timefuncs.Amount
	running              bool
}

func (c *Clock) Controller() moria.Controller {
	duration := mithril.RouteParam("duration").(string)

	totalCost, _ := strconv.ParseFloat(
		mithril.RouteParam("cost").(string),
		64,
	) // Ignoring error: 0 cost per second is valid if not passed.
	c.totalCost = timefuncs.Amount(totalCost)

	fmt.Println(mithril.RouteParam("cost").(string), c.totalCost)

	var err error
	c.totalTime, err = time.ParseDuration(duration)
	if err != nil {
		c.totalTime, _ = time.ParseDuration("1h")
	}
	c.timeSpent = 0 * time.Second
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
			c.timeSpent += now.Sub(c.last)
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

func (c *Clock) MoneySpent() timefuncs.Amount {
	fmt.Printf("Spent %v, Total %v, Cost %v\n",
		c.timeSpent,
		c.totalTime,
		c.totalCost,
	)
	return (c.totalCost / timefuncs.Amount(c.totalTime)) * timefuncs.Amount(c.timeSpent)
}

func (*Clock) View(ctrl moria.Controller) moria.View {
	c := ctrl.(*Clock)

	fmt.Println(c.totalCost)

	styleRed := js.M{}
	if c.timeSpent > c.totalTime {
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
				"value": FormatDuration(c.timeSpent),
				"style": styleRed["style"],
			}),
			m("hr", nil),
			m("div.copy.costIntro", nil, moria.S("COST:")),
			m("div.cost.money", styleRed,
				moria.S(c.MoneySpent().String())),
		),
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
