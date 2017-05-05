# Jenkins to Github Status

Send all queued builds in jenkins as pending status in github

* Call the jenkins api, and retrieve all queued jobs
* Filter out the parent jobs that jenkins insist in putting in that list
* retrieve the relevant information and set the corresponding status in github
