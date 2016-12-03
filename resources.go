package resource

type Source struct {
  User       string `json:"user"`
  Repository string `json:"repository"`

  GitHubAPIURL     string `json:"github_api_url"`
  AccessToken      string `json:"access_token"`
}

type CheckRequest struct {
  Source  Source  `json:"source"`
  Version Version `json:"version"`
}

func NewCheckRequest() CheckRequest {
  res := CheckRequest{}
  return res
}

type Version struct {
  ID string `json:"id,omitempty"`
}
