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
	Type					  	string `json:"type"`
	IDPath					  string `json:"id,omitempty"`
	Ref 							string `json:"ref,omitempty"`
	RefPath      		  string `json:"ref_path,omitempty"`
	Environment       string `json:"environment,omitempty"`
	EnvironmentPath   string `json:"env_path,omitempty"`
	Task              string `json:"task,omitempty"`
	TaskPath          string `json:"task_path,omitempty"`
	State             string `json:"state,omitempty"`
	StatePath         string `json:"state_path,omitempty"`
	Description   		string `json:"description,omitempty"`
	DescriptionPath   string `json:"description_path,omitempty"`
	PayloadPath				string `json:"payload_path,omitempty"`
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
