package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v28/github"
)

type DeploymentOutCommand struct {
	github GitHub
	writer io.Writer
}

func NewDeploymentOutCommand(github GitHub, writer io.Writer) *DeploymentOutCommand {
	return &DeploymentOutCommand{
		github: github,
		writer: writer,
	}
}

func (c *DeploymentOutCommand) Run(sourceDir string, request OutRequest) (OutResponse, error) {
	if request.Params.Ref == nil {
		return OutResponse{}, errors.New("ref is a required parameter")
	}

	newDeployment := &github.DeploymentRequest{
		Ref:              request.Params.Ref,
		RequiredContexts: &[]string{},
	}

	concoursePayload := map[string]interface{}{
		"build_id":            os.Getenv("BUILD_ID"),
		"build_name":          os.Getenv("BUILD_NAME"),
		"build_job_name":      os.Getenv("BUILD_JOB_NAME"),
		"build_pipeline_name": os.Getenv("BUILD_PIPELINE_NAME"),
		"build_team_name":     os.Getenv("BUILD_TEAM_NAME"),
		"build_url": fmt.Sprintf("%v/teams/%v/pipelines/%v/jobs/%v/builds/%v",
			os.Getenv("ATC_EXTERNAL_URL"), os.Getenv("BUILD_TEAM_NAME"), os.Getenv("BUILD_PIPELINE_NAME"), os.Getenv("BUILD_JOB_NAME"), os.Getenv("BUILD_NAME")),
		"atc_external_url": os.Getenv("ATC_EXTERNAL_URL"),
	}

	if request.Params.Payload != nil {
		payload := *request.Params.Payload
		payload["concourse_payload"] = concoursePayload
	} else {
		request.Params.Payload = &map[string]interface{}{
			"concourse_payload": concoursePayload,
		}
	}
	p, err := json.Marshal(request.Params.Payload)
	newDeployment.Payload = github.String(string(p))

	if request.Params.Task != nil {
		newDeployment.Task = request.Params.Task
	}
	if request.Params.Environment != nil {
		newDeployment.Environment = request.Params.Environment
	}
	if request.Params.Description != nil {
		newDeployment.Description = request.Params.Description
	}
	if request.Params.AutoMerge != nil {
		newDeployment.AutoMerge = request.Params.AutoMerge
	}

	fmt.Fprintln(c.writer, "creating deployment")
	deployment, err := c.github.CreateDeployment(newDeployment)
	if err != nil {
		return OutResponse{}, err
	}

	return OutResponse{
		Version:  Version{ID: strconv.FormatInt(*deployment.ID, 10)},
		Metadata: metadataFromDeployment(deployment, []*github.DeploymentStatus{}),
	}, nil
}

func (c *DeploymentOutCommand) fileContents(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
