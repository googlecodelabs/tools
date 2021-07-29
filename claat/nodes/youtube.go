package nodes

// NewYouTubeNode creates a new YouTube video node.
func NewYouTubeNode(vid string) *YouTubeNode {
	return &YouTubeNode{
		node:    node{typ: NodeYouTube},
		VideoID: vid,
	}
}

// YouTubeNode is a YouTube video.
type YouTubeNode struct {
	node
	VideoID string
}

// Empty returns true if yt's VideoID field is zero.
func (yt *YouTubeNode) Empty() bool {
	return yt.VideoID == ""
}
