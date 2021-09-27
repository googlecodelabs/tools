package nodes

import "strings"

type NewImageNodeOptions struct {
	Src   string
	Width float32
	Alt   string
	Title string
}

// NewImageNode creates a new ImageNode with the given options.
// TODO this API is inconsistent with button
func NewImageNode(opts NewImageNodeOptions) *ImageNode {
	return &ImageNode{
		node:  node{typ: NodeImage},
		Src:   opts.Src,
		Width: opts.Width,
		Alt:   opts.Alt,
		Title: opts.Title,
	}
}

// ImageNode represents a single image.
type ImageNode struct {
	node
	Src   string
	Width float32
	Alt   string
	Title string
}

// Empty returns true if its Src is zero, excluding space runes.
func (in *ImageNode) Empty() bool {
	return strings.TrimSpace(in.Src) == ""
}

// ImageNodes extracts everything except NodeImage nodes, recursively.
// TODO rename
func ImageNodes(nodes []Node) []*ImageNode {
	var imgs []*ImageNode
	for _, n := range nodes {
		switch n := n.(type) {
		case *ImageNode:
			imgs = append(imgs, n)
		case *ListNode:
			imgs = append(imgs, ImageNodes(n.Nodes)...)
		case *ItemsListNode:
			for _, i := range n.Items {
				imgs = append(imgs, ImageNodes(i.Nodes)...)
			}
		case *HeaderNode:
			imgs = append(imgs, ImageNodes(n.Content.Nodes)...)
		case *URLNode:
			imgs = append(imgs, ImageNodes(n.Content.Nodes)...)
		case *ButtonNode:
			imgs = append(imgs, ImageNodes(n.Content.Nodes)...)
		case *InfoboxNode:
			imgs = append(imgs, ImageNodes(n.Content.Nodes)...)
		case *GridNode:
			for _, r := range n.Rows {
				for _, c := range r {
					imgs = append(imgs, ImageNodes(c.Content.Nodes)...)
				}
			}
		}
	}
	return imgs
}
