package test

import (
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/test/utils"
	"gitlab-tg-bot/transport/gitlab"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GitlabMergeRequestSuite struct {
	suite.Suite
	tgMessageMock utils.TgMessengerMock
	server        *httptest.Server
	messages      map[httpGitHookTypeKey]string
}

type httpGitHookTypeKey struct {
	model.GitHookType
	model.GitHookSubtype
}

func (g *GitlabMergeRequestSuite) Test_OpenMergeRequest() {
	_, err := http.Post(g.server.URL+model.Gitlab.GetUri(), "", nil)
	if err != nil {
		g.T().Fatal(err)
	}
}

func TestGitlabMergeRequest(t *testing.T) {
	gitlabMergeRequestSuite := &GitlabMergeRequestSuite{tgMessageMock: utils.TgMessengerMock{}}
	handler := gitlab.NewHandler(testEnvironment.application.ServiceStorage, &gitlabMergeRequestSuite.tgMessageMock)

	gitlabMergeRequestSuite.server = httptest.NewServer(handler)

	dir, _ := os.Getwd()
	dir += "/../etc/gitlab/merge-request/"

	files := utils.GetFiles(dir)

	for i, item := range files {
		files[i] = item[strings.Index(item, "{"):]
	}

	suite.Run(t, gitlabMergeRequestSuite)
}
