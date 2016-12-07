package resource

import (
  "io"
  "io/ioutil"
  "errors"
  "fmt"
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
  if request.Params.IDPath == "" {
    return OutResponse{}, errors.New("id is a required parameter")
  }

  id, err := c.fileContents(filepath.Join(sourceDir, request.Params.IDPath))
  if err != nil {
    return OutResponse{}, err
  }

  state := request.Params.State
  if request.Params.StatePath != "" {
    var err error
    state, err = c.fileContents(filepath.Join(sourceDir, request.Params.StatePath))
    if err != nil {
      return OutResponse{}, err
    }
  }

  if state == "" {
    return OutResponse{}, errors.New("state or state_path is a required parameter")
  }

  newStatus := &github.DeploymentStatusRequest{
    State:          github.String(state),
  }

  fmt.Fprintf(c.writer, "creating DeploymentStatus")
  idString, err := strconv.Atoi(id)
  if err != nil {
    return OutResponse{}, err
  }
  status, err := c.github.CreateDeploymentStatus(idString, newStatus)
  if err != nil {
    return OutResponse{}, err
  }

  return OutResponse{
    Version:  Version{ID: id},
    Metadata: metadataFromStatus(status),
  }, nil
}

func (c *OutCommand) fileContents(path string) (string, error) {
  contents, err := ioutil.ReadFile(path)
  if err != nil {
    return "", err
  }

  return strings.TrimSpace(string(contents)), nil
}
