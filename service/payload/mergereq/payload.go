package mergereq

type MergeRequest struct {
	Name string
	Link string

	SourceBranch string
	TargetBranch string

	AssignedTo string
}
type Change struct {
	Old  string
	New  string
	Type ChangeType
}

type ChangeType string

const (
	Rename     ChangeType = "rename"
	ReAssignee ChangeType = "reAssignee"
	Update     ChangeType = "update"
)
