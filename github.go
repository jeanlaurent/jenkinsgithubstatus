package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type githubStatus struct {
	State       string `json:"state"`
	TargetURL   string `json:"target_url"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

func sendGithubStatus(githubOrgProject string, accessTokenWin string, accessTokenMac string, job effectiveJob) error {
	status := githubStatus{State: "pending", TargetURL: job.URL, Description: "Pending", Context: job.Context}

	token := accessTokenMac
	if strings.HasPrefix(job.Context, "win") {
		token = accessTokenWin
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/statuses/%s?access_token=%s", githubOrgProject, job.Commit, token)

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
