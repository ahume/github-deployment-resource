package main

import (
  "encoding/json"
  "errors"
  "os"
  "fmt"

  "github.com/ahume/github-deployment-resource"
)

func main() {
  if len(os.Args) < 2 {
    resource.Sayf("usage: %s <sources directory>\n", os.Args[0])
    os.Exit(1)
  }

  request := resource.NewOutRequest()
  inputRequest(&request)

  sourceDir := os.Args[1]

  github, err := resource.NewGitHubClient(request.Source)
  if err != nil {
    resource.Fatal("constructing github client", err)
  }

  fmt.Println(request.Params)

  command := resource.NewOutCommand(github, os.Stderr)

  if request.Params.Type == "deployment" {
    response, err := command.RunForDeployment(sourceDir, request)
    if err != nil {
      resource.Fatal("running command", err)
    }
    outputResponse(response)
  } else if request.Params.Type == "status" {
    response, err := command.RunForStatus(sourceDir, request)
    if err != nil {
      resource.Fatal("running command", err)
    }
    outputResponse(response)
  } else {
    resource.Fatal("Bad params", errors.New("params.type must be one of `deployment` or `status`"))
  }
}

func inputRequest(request *resource.OutRequest) {
  if err := json.NewDecoder(os.Stdin).Decode(request); err != nil {
    resource.Fatal("reading request from stdin", err)
  }
}

func outputResponse(response resource.OutResponse) {
  if err := json.NewEncoder(os.Stdout).Encode(response); err != nil {
    resource.Fatal("writing response to stdout", err)
  }
}
