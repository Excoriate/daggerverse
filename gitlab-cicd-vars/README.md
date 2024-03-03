# GitLab CI/CD Vars 🚀

![Dagger Version](https://img.shields.io/badge/dagger%20version-%3E=0.10.0-0f0f19.svg?style=flat-square)


A simple [Dagger](https://dagger.io) module to return and fetch your [GitLab CI/CD variables](https://docs.gitlab.com/ee/ci/variables/).

## Configuration 🛠️

Through the [Dagger CLI](https://docs.dagger.io/cli/465058/install), or by using it directly reusing it within your module, you can configure the following options:

* ⚙️ `token`: The GitLab API token to use to authenticate with the GitLab API. This is a required field.

>NOTE: The token should have the `api` scope. For more information about creating a token, please refer to the [GitLab documentation](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html).

---

## Features 🎨

| Command or functionality                                 | Command     | Example                                                                                                                  | Status |
|----------------------------------------------------------|-------------|--------------------------------------------------------------------------------------------------------------------------|--------|
| Fetch a single CI/CD variable from a GitLab project      | **get**     | `dagger call --token=$GITLAB_PRIVATE_TOKEN get --path="mygroup/subgroup/my-project" --var-name="LOOK_FOR_THIS_CICD_VAR"` | ✅      |
| Fetch all CI/CD variables configured in a GitLab project | **get-all** | `dagger call --token=$GITLAB_PRIVATE_TOKEN get-all --path="group/subgroup/my-project"`                                   | ✅      |

> **NOTE**: The `get` command is used to retrieve a specific CI/CD variable by specifying the `--var-name` argument. The `get-all` command fetches all CI/CD variables configured within the specified project. Both commands require a valid GitLab private token passed through the `--token` argument for authentication.

---

## Usage 🚀

From within your Dagger module's directory:

  ```bash
dagger call --token=$GITLAB_PRIVATE_TOKEN get-all --path=group/subgroup/repo-or-project
```

Or, using the public Dagger module:

  ```bash
dagger -m github.com/Excoriate/daggerverse/gitlab-cicd-vars@v1.5.0 call --token=$GITLAB_PRIVATE_TOKEN \
get-all --path=group/subgroup/repo-or-project
```
