package components

import (
	g "maragu.dev/gomponents"
	"maragu.dev/gomponents/html"
)

func HomePage() g.Node {
	return html.H1(g.Text("Hello, World"))
}
