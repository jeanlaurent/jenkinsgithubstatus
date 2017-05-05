package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type effectiveJob struct {
	URL     string
	Context string
	Commit  string
}

func main() {
	jenkinsUser := flag.String("jenkins_user", "", "the name of the jenkins user")
	jenkinsToken := flag.String("jenkins_token", "", "A valid token linked to the jenkins user")
	jenkinsServer := flag.String("jenkins_server", "", "the address of the jenkins server")
	githubToken := flag.String("github_token", "", "A Github token with repo rights.")
	githubProject := flag.String("github_project", "", "A Github org/repo")
	file := flag.String("file", "", "an optional local file")
	readFile := false
	flag.Parse()

	if *file != "" {
		readFile = true
	}

	// fetch jenkins queued jobs
	var allJobs queuedJobs
	var err error
	if readFile {
		fmt.Println("Reading file", *file)
		allJobs, err = readJobsFromDisk(*file)
	} else {
		if *jenkinsServer == "" {
			fmt.Println("no jenkins nor file specified, aborting")
			flag.Usage()
			os.Exit(-1)
		}
		fmt.Println("Fetching queued jobs from", *jenkinsServer)
		allJobs, err = fetchQueuedJobs(*jenkinsUser, *jenkinsToken, *jenkinsServer)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if len(allJobs.Items) == 0 {
		fmt.Println("No queued jobs, life is good.")
		os.Exit(0)
	}
	fmt.Println("Queued jobs ", len(allJobs.Items))
	// filter out parent master jobs
	downStreamJobs := filterParentJob(allJobs)
	if len(downStreamJobs) == 0 {
		fmt.Println("No queued downstream jobs, life is good.")
		os.Exit(0)
	}
	fmt.Println("Downstream jobs ", len(downStreamJobs))

	// Update github status for each downstream job
	jobs := transformJenkinsJobsIntoGithubJobs(downStreamJobs)
	for _, job := range jobs {
		fmt.Println(job.Context)
		if *githubToken == "" || *githubProject == "" {
			continue
		}
		err = sendGithubStatus(*githubProject, *githubToken, job)
		if err != nil {
			fmt.Println(err)
		}
	}
	if *githubToken == "" || *githubProject == "" {
		fmt.Println("No githubtoken present skipping sending status to github")
		os.Exit(-1)
	}
}

func transformJenkinsJobsIntoGithubJobs(jenkinsJobs []item) []effectiveJob {
	var jobs []effectiveJob
	for _, job := range jenkinsJobs {
		commit, err := find("commit", job)
		if err != nil {
			fmt.Println(err)
			continue
		}
		url := job.Task.URL
		name := extractContext(job.Task.Name)
		jobs = append(jobs, effectiveJob{URL: url, Context: name, Commit: commit})
	}
	return jobs
}

func filterParentJob(allJobs queuedJobs) []item {
	var downStreamJob []item
	for _, item := range allJobs.Items {
		lastAction := item.Actions[len(item.Actions)-1]
		for _, cause := range lastAction.Causes {
			if cause.UpstreamBuild > 0 {
				downStreamJob = append(downStreamJob, item)
			}
		}
	}
	return downStreamJob
}

func find(key string, job item) (string, error) {
	for _, action := range job.Actions {
		for _, parameter := range action.Parameters {
			if parameter.Name == key {
				return parameter.Value, nil
			}
		}
	}
	return "", errors.New("can't find key")
}

// split something like label=win-srv2016,testsToRun=compose
// into win-srv2016-compose
func extractContext(context string) string {
	var values []string
	words := strings.Split(context, ",")
	for _, word := range words {
		keyValue := strings.Split(word, "=")
		values = append(values, keyValue[1])
	}
	return strings.Join(values, "-")
}
