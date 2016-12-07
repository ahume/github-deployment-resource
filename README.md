# GitHub Deployment Resource

Fetches and creates GitHub Deployments and Deployment Statuses. For more about Github Deployments
see the [API documentation](https://developer.github.com/v3/repos/deployments/).

## Source Configuration

* `user`: *Required.* The GitHub username or organization name for the
  repository that you are deploying.

* `repository`: *Required.* The repository name that you are deploying.

* `access_token`: *Required.* Used for accessing deployment data and creating deployments
  and deployment statuses.

* `github_api_url`: *Optional.* If you use a non-public GitHub deployment then
  you can set your API URL here.

* `environments`: *Optional.* A list of environments to get versions for.

### Example

``` yaml
- name: gh-deployment
  type: github-deployment
  source:
    user: BrandwatchLtd
    repository: analytics
    access_token: abcdef1234567890
```

``` yaml
- get: gh-deployment
```

``` yaml
- put: gh-deployment
  params:
    id: path/to/id/file
    state: success
```

## Behavior

### `check`: Check for Deployments

`/check` always returns the single latest deployment. It assumes that any preceding deployments
are invalidated by the existence of a later deployment.

### `in`: Fetch Deployment

Fetches the latest deployment and creates the following files:

* `id` containing the `id` of the deployment being fetched.
* `ref` containting the name of the `ref` the deployment is relating to. A branch, tag, or SHA.
* `sha` containg the SHA that was recorded at deployment creation time.
* `task` containing the name of the task for the deployment.
* `environment` containing the name of the environment that is being deployed to.
* `description` containing the description of the deployment
* `deploymentJSON` containing the full JSON of the deployment as received from the API.


### `out`: Create a Deployment or DeploymentStatus

Create a new Deployment, or update a given Deployment with a new DeploymentStatus

#### Parameters

* `type`: *Required.* Either `deployment` or `status`.

##### If type=status

* `id`: *Required.* A path to a file containing the ID of the deployment to update
  with the new status.

* `state`: *Required.*  A path to a file containing the state of the new deployment status.
  Must be one of `pending`, `success`, `error`, `inactive`, or `failure`.

##### If type=deployment

* `ref`: *Optional.* A path to a file containing the ref of the deployment. A branch name, a tag,
  or SHA.

* `environment`: *Optional.* A path to a file containing the name of the environment that is being
  deployed to.

* `description`: *Optional.* A path to a file containing the description of the deployment.

* `payload`: *Optional.* A path to a JSON file containing any additional data about the deployment.

* `task`: *Optional.* A path to a file containing the name of the task for the deployment.
