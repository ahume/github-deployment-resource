package resource

type Source struct {
	User         string `json:"user"`
	Repository   string `json:"repository"`
	GitHubAPIURL string `json:"github_api_url"`
	AccessToken  string `json:"access_token"`
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

func NewCheckRequest() CheckRequest {
	res := CheckRequest{}
	return res
}

func NewInRequest() InRequest {
	return InRequest{}
}

type InResponse struct {
	Version  Version        `json:"version"`
	Metadata []MetadataPair `json:"metadata"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
