package model

import (
	"fmt"
	"net"
)

type XGitlabEvent string

const (
	HookTypePush             XGitlabEvent = "Push Hook"
	HookTypePushTag          XGitlabEvent = "Tag Push Hook"
	HookTypePushIssue        XGitlabEvent = "Issue Hook"
	HookTypePushNote         XGitlabEvent = "Note Hook"
	HookTypePushMergeRequest XGitlabEvent = "Merge Request Hook"
	HookTypePushWiki         XGitlabEvent = "Wiki Page Hook"
	HookTypePushPipeline     XGitlabEvent = "Pipeline Hook"
	HookTypePushJob          XGitlabEvent = "Job Hook"
	HookTypePushDeployment   XGitlabEvent = "Deployment Hook"
	HookTypePushMember       XGitlabEvent = "Member Hook"
	HookTypePushSubgroup     XGitlabEvent = "Subgroup Hook"
	HookTypePushFeatureFlag  XGitlabEvent = "Feature Flag Hook"
	HookTypePushRelease      XGitlabEvent = "Release Hook"
)

type Hook struct {
	RepoId                   string
	PushEvents               bool
	IssuesEvents             bool
	ConfidentialIssuesEvents bool
	MergeRequestsEvents      bool
	TagPushEvents            bool
	NoteEvents               bool
	JobEvents                bool
	PipelineEvents           bool
	WikiPageEvents           bool
}

func GetMyInterfaceAddr() (net.IP, error) {

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	addresses := []net.IP{}
	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			addresses = append(addresses, ip)
		}
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no address Found, net.InterfaceAddrs: %v", addresses)
	}
	//only need first
	return addresses[0], nil
}
