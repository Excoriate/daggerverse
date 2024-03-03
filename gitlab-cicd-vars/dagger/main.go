package main

import (
	"errors"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type GitlabCicdVars struct {
	// token to use for gitlab api
	Token string
}

func New(
	// token is the GitLab API token to use, for information about how to create a token,
	//see https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html
	token string,
) *GitlabCicdVars {
	g := &GitlabCicdVars{}

	g.Token = token

	return g
}

// getGitLabClient returns a new GitLab client
func getGitLabClient(token string) (*gitlab.Client, error) {
	if token == "" {
		return nil, errors.New("the gitlab token is required")
	}

	client, err := gitlab.NewClient(token)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// listVariables returns a list of all the variables in a project
func listVariables(c *gitlab.Client, namespace string) ([]string, error) {
	project, _, err := c.Projects.GetProject(namespace, &gitlab.GetProjectOptions{})
	if err != nil {
		return nil, err
	}

	variables, _, err := c.ProjectVariables.ListVariables(project.ID, &gitlab.ListProjectVariablesOptions{})
	if err != nil {
		return nil, err
	}

	var variablesToReturn []string
	for _, variable := range variables {
		variablesToReturn = append(variablesToReturn, fmt.Sprintf("%s=%s", variable.Key, variable.Value))
	}

	return variablesToReturn, nil
}

// GetAll returns all the variables in a project
func (g *GitlabCicdVars) GetAll(
	// path is the path to the GitLab's project, also known as 'namespace'
	path string,
) (string, error) {
	client, err := getGitLabClient(g.Token)
	if err != nil {
		return "", err
	}

	variablesToReturn, err := listVariables(client, path)
	if err != nil {
		return "", err
	}

	return getInTableFmt(variablesToReturn), nil
}

func (g *GitlabCicdVars) Get(
	// path is the path to the GitLab's project, also known as 'namespace'
	path,
	// varName is the name of the variable to get
	varName string) (string, error) {
	client, err := getGitLabClient(g.Token)
	if err != nil {
		return "", err
	}

	listOfVariablesInStr, err := listVariables(client, path)
	if err != nil {
		return "", err
	}

	for _, variable := range listOfVariablesInStr {
		if strings.HasPrefix(variable, varName) {
			return variable, nil
		}
	}

	return "", fmt.Errorf("variable %s not found", varName)
}

func getInTableFmt(variables []string) string {
	// Define the table header
	table := "\nCI/CD variable 🚀| Value👀\n"
	// Add a separator for clarity
	table += strings.Repeat("-", 50) + "\n"

	for _, variable := range variables {
		parts := strings.SplitN(variable, "=", 2)
		if len(parts) == 2 {
			table += fmt.Sprintf("%-30s | %s\n", parts[0], parts[1])
		}
	}

	return table
}
