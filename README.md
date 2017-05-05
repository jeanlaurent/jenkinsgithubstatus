# Jenkins to Github Status

Send all queued builds in jenkins as pending status in github

* Call the jenkins api, and retrieve all queued jobs
* Filter out the parent jobs that jenkins insist in putting in that list
* Retrieve the relevant information from those jobs and set the corresponding status in github


* usage :
```
jenkinsgithubstatus
no jenkins nor file specified, aborting
Usage of jenkinsgithubstatus:
  -file string
    	an optional local file
  -github_project string
    	A Github org/repo
  -github_token string
    	A Github token with repo rights.
  -jenkins_server string
    	the address of the jenkins server
  -jenkins_token string
    	A valid token linked to the jenkins user
  -jenkins_user string
    	the name of the jenkins user
```

* Sample Output
```
Queued jobs  17
Downstream jobs  12
win-10586-others
win-10586-compose
mac-1012-others
mac-1011-others
mac-1011-network
mac-1012-network
mac-1012-compose
mac-1011-compose
mac-1012-compose
mac-1011-others
mac-1011-compose
mac-1011-network
```
