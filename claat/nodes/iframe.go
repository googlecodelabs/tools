package nodes

// iframe allowlist - set of domains allow to embed iframes in a codelab.
// TODO make this configurable somehow
var IframeAllowlist = []string{
	"carto.com",
	"codepen.io",
	"dartlang.org",
	"dartpad.dev",
	"demo.arcade.software",
	"github.com",
	"glitch.com",
	"google.com",
	"google.dev",
	"observablehq.com",
	"repl.it",
	"web.dev",
}

// NewIframeNode creates a new embedded iframe.
func NewIframeNode(url string) *IframeNode {
	return &IframeNode{
		node: node{typ: NodeIframe},
		URL:  url,
	}
}

// IframeNode is an embeddes iframe.
type IframeNode struct {
	node
	URL string
}

// Empty returns true if iframe's URL field is empty.
func (iframe *IframeNode) Empty() bool {
	return iframe.URL == ""
}
