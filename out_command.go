package resource

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ahume/go-github/github"
)

type OutCommand struct {
	github GitHub
	writer io.Writer
}

func NewOutCommand(github GitHub, writer io.Writer) *OutCommand {
	return &OutCommand{
		github: github,
		writer: writer,
	}
}

func (c *OutCommand) Run(sourceDir string, request OutRequest) (OutResponse, error) {
	if request.Params.ID == nil {
		return OutResponse{}, errors.New("id is a required parameter")
	}

	id, ok := request.Params.ID.(string)
	if ok != true {
		var err error
		v := request.Params.ID.(map[string]interface{})
		id, err = c.fileContents(filepath.Join(sourceDir, v["file"].(string)))
		if err != nil {
			return OutResponse{}, errors.New("id or id.file is a required param")
		}
	}

	state, ok := request.Params.State.(string)
	if ok != true {
		var err error
		v := request.Params.State.(map[string]interface{})
		state, err = c.fileContents(filepath.Join(sourceDir, v["file"].(string)))
		if err != nil {
			return OutResponse{}, errors.New("state or state.file is a required parameter")
		}
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return OutResponse{}, err
	}
	deployment, err := c.github.GetDeployment(idInt)
	if err != nil {
		return OutResponse{}, err
	}

	newStatus := &github.DeploymentStatusRequest{
		State: github.String(state),
	}

	fmt.Fprintf(c.writer, "creating DeploymentStatus")
	status, err := c.github.CreateDeploymentStatus(*deployment.ID, newStatus)
	if err != nil {
		return OutResponse{}, err
	}

	return OutResponse{
		Version:  Version{ID: strconv.Itoa(*deployment.ID)},
		Metadata: metadataFromDeployment(deployment, status),
	}, nil
}

func (c *OutCommand) fileContents(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(contents)), nil
}
