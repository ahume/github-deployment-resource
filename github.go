package resource

import (
  "net/http"

  "golang.org/x/oauth2"

  "github.com/google/go-github/github"
)

//go:generate counterfeiter . GitHub

type GitHub interface {
  ListDeployments() ([]*github.Deployment, error)
}

type GitHubClient struct {
  client *github.Client

  user       string
  repository string
}

func NewGitHubClient(source Source) (*GitHubClient, error) {
  var client *github.Client

  if source.AccessToken == "" {
    client = github.NewClient(nil)
  } else {
    var err error
    client, err = oauthClient(source)
    if err != nil {
      return nil, err
    }
  }

  return &GitHubClient{
    client:     client,
    user:       source.User,
    repository: source.Repository,
  }, nil
}

func (g *GitHubClient) ListDeployments() ([]*github.Deployment, error) {
  releases, res, err := g.client.Repositories.ListDeployments(g.user, g.repository, nil)
  if err != nil {
    return []*github.Deployment{}, err
  }

  err = res.Body.Close()
  if err != nil {
    return nil, err
  }

  return releases, nil
}

func oauthClient(source Source) (*github.Client, error) {
  ts := oauth2.StaticTokenSource(&oauth2.Token{
    AccessToken: source.AccessToken,
  })

  oauthClient := oauth2.NewClient(oauth2.NoContext, ts)

  githubHTTPClient := &http.Client{
    Transport: oauthClient.Transport,
  }

  return github.NewClient(githubHTTPClient), nil
}
