package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type githubStatus struct {
	State       string `json:"state"`
	TargetURL   string `json:"target_url"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

func sendGithubStatus(githubOrgProject string, accessToken string, job effectiveJob) error {
	status := githubStatus{State: "pending", TargetURL: job.URL, Description: "Waiting for a build node", Context: job.Context}
	url := fmt.Sprintf("https://api.github.com/repos/%v/statuses/%v?access_token=%v", githubOrgProject, job.Commit, accessToken)

	jsonBody, err := json.Marshal(status)
	if err != nil {
		return err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Non OK http response code (%v)", response.StatusCode)
	}
	return nil
}
