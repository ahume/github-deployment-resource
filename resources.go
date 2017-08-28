package resource

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterbourgon/mergemap"
)

type Source struct {
	User         string `json:"user"`
	Repository   string `json:"repository"`
	AccessToken  string `json:"access_token"`
	GitHubAPIURL string `json:"github_api_url"`
	Environment  string `json:"environment"`
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
	Type        string `json:"type"`
	ID          string
	Ref         string
	Environment string
	Task        string
	State       string
	Description string
	Payload     map[string]interface{}
	PayloadPath string `json:"payload_path"`

	RawID          json.RawMessage `json:"id"`
	RawState       json.RawMessage `json:"state"`
	RawRef         json.RawMessage `json:"ref"`
	RawTask        json.RawMessage `json:"task"`
	RawEnvironment json.RawMessage `json:"environment"`
	RawDescription json.RawMessage `json:"description"`
	RawPayload     json.RawMessage `json:"payload"`
}

// Used to avoid recursion in UnmarshalJSON below.
type outParams OutParams

func (p *OutParams) UnmarshalJSON(b []byte) (err error) {
	j := outParams{
		Type: "status",
	}

	if err = json.Unmarshal(b, &j); err == nil {
		*p = OutParams(j)
		if p.RawID != nil {
			p.ID = getStringOrStringFromFile(p.RawID)
		}

		if p.RawState != nil {
			p.State = getStringOrStringFromFile(p.RawState)
		}

		if p.RawRef != nil {
			p.Ref = getStringOrStringFromFile(p.RawRef)
		}

		if p.RawTask != nil {
			p.Task = getStringOrStringFromFile(p.RawTask)
		}

		if p.RawEnvironment != nil {
			p.Environment = getStringOrStringFromFile(p.RawEnvironment)
		}

		if p.RawDescription != nil {
			p.Description = getStringOrStringFromFile(p.RawDescription)
		}

		var payload map[string]interface{}
		json.Unmarshal(p.RawPayload, &payload)

		if p.PayloadPath != "" {
			stringFromFile := fileContents(p.PayloadPath)
			var payloadFromFile map[string]interface{}
			json.Unmarshal([]byte(stringFromFile), &payloadFromFile)

			payload = mergemap.Merge(payloadFromFile, payload)
		}

		p.Payload = payload

		return
	}
	return
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
	return OutRequest{}
}

func getStringOrStringFromFile(field json.RawMessage) string {
	var rawValue interface{}
	if err := json.Unmarshal(field, &rawValue); err == nil {
		switch rawValue := rawValue.(type) {
		case string:
			return rawValue
		case map[string]interface{}:
			return fileContents(rawValue["file"].(string))
		default:
			panic("Could not read string out of Params field")
		}
	}
	return ""
}

func fileContents(path string) string {
	sourceDir := os.Args[1]
	contents, err := ioutil.ReadFile(filepath.Join(sourceDir, path))
	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(contents))
}
