package components

import (
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
	left          time.Duration
	last          time.Time
	costPerSecond float64
	running       bool
}

func (c *Clock) Controller() moria.Controller {
	duration := mithril.RouteParam("duration").(string)
	c.costPerSecond, _ = strconv.ParseFloat(
		mithril.RouteParam("cost").(string),
		64,
	) // Ignoring error: 0 cost per second is valid if not passed.
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

func (*Clock) View(ctrl moria.Controller) moria.View {
	c := ctrl.(*Clock)

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
				moria.S(strconv.FormatFloat(c.left.Seconds()*c.costPerSecond,
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
