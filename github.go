package resource

import (
	"net/http"

	"golang.org/x/oauth2"

	"github.com/ahume/go-github/github"
)

//go:generate counterfeiter -o fakes/fake_git_hub.go . GitHub

type GitHub interface {
	ListDeployments() ([]*github.Deployment, error)
	ListDeploymentStatuses(ID int) ([]*github.DeploymentStatus, error)
	GetDeployment(ID int) (*github.Deployment, error)
	CreateDeployment(request *github.DeploymentRequest) (*github.Deployment, error)
	CreateDeploymentStatus(ID int, request *github.DeploymentStatusRequest) (*github.DeploymentStatus, error)
}

type GitHubClient struct {
	client *github.Client

	user       string
	repository string
}

func NewGitHubClient(source Source) (*GitHubClient, error) {
	var client *github.Client

	client, err := oauthClient(source)
	if err != nil {
		return nil, err
	}

	return &GitHubClient{
		client:     client,
		user:       source.User,
		repository: source.Repository,
	}, nil
}

func (g *GitHubClient) ListDeployments() ([]*github.Deployment, error) {
	deployments, res, err := g.client.Repositories.ListDeployments(g.user, g.repository, nil)
	if err != nil {
		return []*github.Deployment{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

func (g *GitHubClient) GetDeployment(ID int) (*github.Deployment, error) {
	deployment, res, err := g.client.Repositories.GetDeployment(g.user, g.repository, ID)
	if err != nil {
		return &github.Deployment{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (g *GitHubClient) CreateDeployment(request *github.DeploymentRequest) (*github.Deployment, error) {
	deployment, res, err := g.client.Repositories.CreateDeployment(g.user, g.repository, request)
	if err != nil {
		return &github.Deployment{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (g *GitHubClient) ListDeploymentStatuses(ID int) ([]*github.DeploymentStatus, error) {
	statuses, res, err := g.client.Repositories.ListDeploymentStatuses(g.user, g.repository, ID, nil)
	if err != nil {
		return []*github.DeploymentStatus{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (g *GitHubClient) CreateDeploymentStatus(ID int, request *github.DeploymentStatusRequest) (*github.DeploymentStatus, error) {
	status, res, err := g.client.Repositories.CreateDeploymentStatus(g.user, g.repository, ID, request)
	if err != nil {
		return &github.DeploymentStatus{}, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return status, nil
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
