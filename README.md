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

* `environments`: *Optional.* A list of environments to get versions for. [Not implemented]

### Example

``` yaml
- name: deployment
  type: github-deployment
  source:
    user: BrandwatchLtd
    repository: analytics
    access_token: abcdef1234567890
```

``` yaml
- get: deployment
```

``` yaml
- put: deployment
  params:
    id:
      file: deployment/id # path to a file containing the deployment ID
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

* `type`: *Optional.* Either `deployment` or `status`. Defaults to `status`.

##### If type=status

* `id`: *Required.* The ID of the deployment to update with the new status.
  NB: You'll most likely want to reference a file with this ID stored in (see below).

* `state`: *Required.*  The state of the new deployment status.
  Must be one of `pending`, `success`, `error`, `inactive`, or `failure`.

##### If type=deployment

* `ref`: *Optional.* The ref of the deployment. A branch name, a tag, or SHA.

* `environment`: *Optional.* The name of the environment that is being deployed to.

* `description`: *Optional.* The description of the deployment.

* `payload`: *Optional.* Additional data about the deployment.

* `task`: *Optional.* The name of the task for the deployment.

##### Reading values from files

All of the above parameters can be used to pass the name of a file to read the applicable value
from. For example...

```yaml
- put: deployment
  params:
    id:
      file: path/to/the/id/file
    state: success
    description:
      file: path/to/the/description
```
The above configuration, would read in the `id` and `description` values from files, but use
the `state` value which has been passed in directly as a string.
