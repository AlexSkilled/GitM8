package test

import (
	"gitlab-tg-bot/service/model"
	"gitlab-tg-bot/test/utils"
	"gitlab-tg-bot/transport/gitlab"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

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

}

func TestGitlabMergeRequest(t *testing.T) {
	gitlabMergeRequestSuite := &GitlabMergeRequestSuite{tgMessageMock: utils.TgMessengerMock{}}
	handler := gitlab.NewHandler(testEnvironment.application.ServiceStorage, &gitlabMergeRequestSuite.tgMessageMock)

	gitlabMergeRequestSuite.server = httptest.NewServer(handler)

	dir, _ := os.Getwd()
	dir += "/../etc/gitlab/merge-request"

	filesMigrations, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Fatal(err)
	}
	files := make([]string, 0, len(filesMigrations))
	for _, f := range filesMigrations {
		name := f.Name()
		if strings.Contains(name, ".http") {
			files = append(files, name)
		}
	}

	suite.Run(t, gitlabMergeRequestSuite)
}
