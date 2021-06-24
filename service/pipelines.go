package service

import (
	"fmt"
	"gitlab-tg-bot/data/provider"
	"gitlab-tg-bot/internal/model"
	"log"
	"strconv"
)

type PipeService struct {
	UserProvider provider.UserProvider
}

func (s *PipeService) GetJobInfo(r model.MergeRequest, id int32) (string, error) {
	user, err := s.UserProvider.Get(id)
	if err != nil {
		return "", err
	}

	if r.ObjectAttributes.Status != "Success" {
		return "", nil
	}

	client := user.GetGitlabClient()

	jobInfo, _, err := client.ProjectJobTrace(strconv.Itoa(r.Project.Id), 8236)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(jobInfo)

	if err != nil {
		log.Fatalln(err)
	}
	return "", nil
}
