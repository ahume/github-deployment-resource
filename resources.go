package resource

type Source struct {
	User         string `json:"user"`
	Repository   string `json:"repository"`
	AccessToken  string `json:"access_token"`
	GitHubAPIURL string `json:"github_api_url"`
}

type Version struct {
	ID string `json:"id,omitempty"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type InRequest struct {
	Source  Source   `json:"source"`
	Version *Version `json:"version"`
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
	Type					  	string `json:"type"`
	IDPath					  string `json:"id,omitempty"`
	RefPath      		  string `json:"ref,omitempty	"`
	EnvironmentPath   string `json:"env,omitempty"`
	TaskPath          string `json:"task,omitempty"`
	StatePath         string `json:"state,omitempty"`
	DescriptionPath   string `json:"description,omitempty"`
	PayloadPath				string `json:"payload,omitempty"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewCheckRequest() CheckRequest {
	res := CheckRequest{}
	return res
}

func NewInRequest() InRequest {
	return InRequest{}
}

func NewOutRequest() OutRequest {
	return OutRequest{}
}
