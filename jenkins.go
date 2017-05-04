package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type cause struct {
	UpstreamBuild int `json:"upstreamBuild"`
}

type parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type action struct {
	Causes     []cause     `json:"causes"`
	Parameters []parameter `json:"parameters"`
}

type task struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type item struct {
	Actions []action `json:"actions"`
	//Params  string   `json:"params"`
	Task task `json:"task"`
}

type queuedJobs struct {
	Items []item `json:"items"`
}

func readJobsFromDisk(file string) (queuedJobs, error) {
	var emptyJobs queuedJobs
	queueAsJSON, err := ioutil.ReadFile(file)
	if err != nil {
		return emptyJobs, err
	}
	return unmarshal(queueAsJSON)
}

func fetchQueuedJobs(user string, token string, server string) (queuedJobs, error) {
	var emptyJobs queuedJobs
	crumb, err := getJenkinsCrumb(user, token, server)
	if err != nil {
		return emptyJobs, err
	}
	queueAsJSON, err := retrieveQueuedJob(user, token, server, crumb)
	if err != nil {
		return emptyJobs, err
	}
	return unmarshal(queueAsJSON)
}

func unmarshal(jsonPayload []byte) (queuedJobs, error) {
	var jobs queuedJobs
	err := json.Unmarshal(jsonPayload, &jobs)
	if err != nil {
		return jobs, err
	}
	return jobs, nil
}

func retrieveQueuedJob(user string, token string, server string, jenkinsCrumb string) ([]byte, error) {
	var url = fmt.Sprintf("https://%v:%v@%v/queue/api/json", user, token, server)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	crumbKeyValue := strings.Split(jenkinsCrumb, ":")
	request.Header.Set(crumbKeyValue[0], crumbKeyValue[1])

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	//GET -H "$CRUMB" "https://$JENKINS_USER:$JENKINS_TOKEN@$JENKINS_SERVER/queue/api/json")
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Non http response code (%v)", response.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func getJenkinsCrumb(user string, token string, server string) (string, error) {
	var url = fmt.Sprintf("https://%v:%v@%v/crumbIssuer/api/xml?xpath=concat(//crumbRequestField,\":\",//crumb)", user, token, server)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", fmt.Errorf("Non http response code (%v)", response.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(bodyBytes), nil
}
