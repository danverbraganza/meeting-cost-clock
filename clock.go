package main

import (
	"github.com/danverbraganza/go-mithril/moria"
	"honnef.co/go/js/dom"

	"meeting-cost-clock/components"
)

func main() {
	myComponent := &components.Chooser{}
	myClock := &components.Clock{}

	moria.Route(
		dom.GetWindow().Document().QuerySelector("body"), "/",
		map[string]moria.Component{
			"/":                      myComponent,
			"/clock/:duration/:cost": myClock,
		})
}
