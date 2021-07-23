package nodes

// TODO general refactor?

// SurveyNode contains groups of questions. Each group name is the Survey key.
type SurveyNode struct {
	node
	ID     string
	Groups []*SurveyGroup
}

// SurveyGroup contains group name/question and possible answers.
type SurveyGroup struct {
	Name    string
	Options []string
}

// NewSurveyNode creates a new survey node with optional questions.
// If survey is nil, a new empty map will be created.
// TODO is "map" above a mistake, or should the code below contain a map?
func NewSurveyNode(id string, groups ...*SurveyGroup) *SurveyNode {
	return &SurveyNode{
		node:   node{typ: NodeSurvey},
		ID:     id,
		Groups: groups,
	}
}

// Empty returns true if each group has 0 options.
func (sn *SurveyNode) Empty() bool {
	for _, g := range sn.Groups {
		if len(g.Options) > 0 {
			return false
		}
	}
	return true
}
