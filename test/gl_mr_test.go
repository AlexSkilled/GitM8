package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GitlabMergeRequestSuite struct {
	suite.Suite
}

func (g *GitlabMergeRequestSuite) Test_OpenMergeRequest() {

}

func TestGitlabMergeRequest(t *testing.T) {
	suite.Run(t, new(GitlabMergeRequestSuite))
}
