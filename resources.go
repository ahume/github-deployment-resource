package resource

import (
	"encoding/json"
)

type Source struct {
	User         string `json:"user"`
	Repository   string `json:"repository"`
	AccessToken  string `json:"access_token"`
	GitHubAPIURL string `json:"github_api_url"`
}

type Version struct {
	ID       string `json:"id"`
	Statuses string `json:"status"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

type InResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

type OutResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

type OutParams struct {
	Type        string      `json:"type"`
	ID          interface{} `json:"id"`
	Ref         interface{}
	Environment interface{}
	Task        interface{}
	State       interface{}
	Description interface{}
	Payload     json.RawMessage
	PayloadPath string `json:"payload_path"`
}

type File struct {
	File string `json:"file"`
}

type Filer struct {
	File interface{} `json:"file"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewCheckRequest() CheckRequest {
	return CheckRequest{}
}

func NewInRequest() InRequest {
	return InRequest{}
}

func NewOutRequest() OutRequest {
	return OutRequest{
		Params: OutParams{
			Type: "status",
		},
	}
}
