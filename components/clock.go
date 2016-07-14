package components

import (
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
	s     = moria.S
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
	return (c.totalCost / timefuncs.Amount(c.totalTime)) * timefuncs.Amount(c.timeSpent)
}

func (*Clock) View(ctrl moria.Controller) moria.View {
	c := ctrl.(*Clock)

	maybeRed := js.M{}
	if c.timeSpent > c.totalTime {
		maybeRed["style"] = "color:darkred;"
	}

	pauseSigil := s("\u23F8")
	if !c.running {
		pauseSigil = s("\u25B6")
	}

	return m("div#wrapper", nil,
		m("h1", nil, s("How much is this meeting costing?")),
		m("table#display", nil,
			m("tr", nil,
				m("label.copy[for='totalTime']", nil, s("Time Elapsed")),
				m("input#totalTime", js.M{
					"value": timefuncs.FormatDuration(c.timeSpent),
					"style": maybeRed["style"],
				})),
			m("tr", nil,
				m("label.copy.costIntro", nil, s("Cost")),
				m("span.cost.money", maybeRed,
					s(c.MoneySpent().String())),
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
			s("\u25a0")))
}
