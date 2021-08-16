package mergereq

type MergeRequest struct {
	Name string
	Link string

	SourceBranch string
	TargetBranch string

	AssignedTo string
}
