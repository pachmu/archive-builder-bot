package main

import (
	"fmt"
	"github.com/pachmu/gojenkins"
	"strings"
)

type JenkinsParams struct {
	JenkinsUrl string
	Username   string
	Password   string
}

type Build struct {
	number int64
	branch string
}

type JenkinsClient interface {
	GetLastBuilds(jobName string) ([]Build, error)
	StartBuild(jobName string) error
}

type jkClient struct {
	jenkins *gojenkins.Jenkins
}

func GetNewJenkinsClinet(jenkinsParams JenkinsParams) (JenkinsClient, error) {
	jenkins := gojenkins.CreateJenkins(nil, jenkinsParams.JenkinsUrl, jenkinsParams.Username, jenkinsParams.Password)
	// Provide CA certificate if server is using self-signed certificate
	// caCert, _ := ioutil.ReadFile("/tmp/ca.crt")
	// jenkins.Requester.CACert = caCert
	_, err := jenkins.Init()

	if err != nil {
		return nil, fmt.Errorf("client init failed, got %s", err.Error())
	}
	return &jkClient{
		jenkins: jenkins,
	}, nil
}

func (jk jkClient) GetLastBuilds(jobName string) ([]Build, error) {
	build, err := jk.jenkins.GetJob(jobName)
	if err != nil {
		return nil, fmt.Errorf("job does not exist, got %s", err.Error())
	}
	builds, err := build.GetAllBuildIds()
	if err != nil {
		return nil, fmt.Errorf("failed to get builds, got %s", err.Error())
	}
	var res []Build
	for _, b := range builds {
		info, err := build.GetBuild(b.Number)
		if err != nil {
			return nil, fmt.Errorf("failed to get build %d, got %s", b.Number, err.Error())
		}
		var branchName string
		if len(info.Raw.Actions) > 0 {
			for _, action := range info.Raw.Actions {
				if len(action.LastBuiltRevision.Branch) > 0 {
					branchName = action.LastBuiltRevision.Branch[0].Name
				}
			}
		}
		res = append(res, Build{
			number: info.Raw.Number,
			branch: trimBranch(branchName),
		})
	}
	return res, nil
}

func (jk jkClient) StartBuild(jobName string) error {
	return nil
}

func trimBranch(b string) string {
	trimmed := strings.Replace(b, "refs/remotes/", "", -1)
	return strings.Replace(trimmed, "origin/", "", -1)
}
