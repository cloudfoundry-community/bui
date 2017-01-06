package bosh

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/starkandwayne/goutils/log"
)

// GetStemcells from given BOSH
func (c *Client) GetStemcells(auth Auth) (stemcells []Stemcell, err error) {
	r := c.NewRequest("GET", "/stemcells")
	respBody, err := c.DoAuthRequest(r, auth)
	if err != nil {
		log.Errorf("GetStemcells - requesting stemcells  %v", err)
		return
	}
	err = json.Unmarshal(respBody, &stemcells)
	if err != nil {
		log.Errorf("GetStemcells - unmarshalling stemcells %v, payload %s", err, string(respBody))
		return
	}
	return
}

// GetReleases from given BOSH
func (c *Client) GetReleases(auth Auth) (releases []Release, err error) {
	r := c.NewRequest("GET", "/releases")
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("GetReleases - requesting releases  %v", err)
		return
	}
	err = json.Unmarshal(respBody, &releases)
	if err != nil {
		log.Errorf("GetReleases - unmarshalling releases %v, payload %s", err, string(respBody))
		return
	}
	return
}

// GetDeployments from given BOSH
func (c *Client) GetDeployments(auth Auth) (deployments []Deployment, err error) {
	r := c.NewRequest("GET", "/deployments")
	respBody, err := c.DoAuthRequest(r, auth)
	if err != nil {
		log.Errorf("GetDeployments - requesting deployments  %v", err)
		return
	}
	err = json.Unmarshal(respBody, &deployments)
	if err != nil {
		log.Errorf("GetDeployments - unmarshalling deployments %v, payload %s", err, string(respBody))
		return
	}
	return
}

// GetDeployment from given BOSH
func (c *Client) GetDeployment(name string, auth Auth) (manifest Manifest, err error) {
	r := c.NewRequest("GET", "/deployments/"+name)
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("GetDeployment - requesting deployment manifest %v", err)
		return
	}
	err = json.Unmarshal(respBody, &manifest)
	if err != nil {
		log.Errorf("GetDeployment - unmarshalling deployment manifest %v, payload %s", err, string(respBody))
		return
	}
	return
}

// CreateDeployment from given BOSH
func (c *Client) CreateDeployment(manifest string, auth Auth) (task Task, err error) {
	r := c.NewRequest("POST", "/deployments")
	buffer := bytes.NewBufferString(manifest)
	r.Body = buffer
	r.Header["Content-Type"] = "text/yaml"

	respBody, err := c.DoAuthRequest(r, auth)
	if err != nil {
		log.Errorf("CreateDeployment - requesting create deployment %v", err)
		return
	}
	err = json.Unmarshal(respBody, &task)
	if err != nil {
		log.Errorf("CreateDeployment - unmarshalling task %v, payload %s", err, string(respBody))
		return
	}
	return
}

// GetDeploymentVMs from given BOSH
func (c *Client) GetDeploymentVMs(name string, auth Auth) (vms []VM, err error) {
	var task Task
	r := c.NewRequest("GET", "/deployments/"+name+"/vms?format=full")
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("GetDeploymentVMs - requesting deployment vms %v", err)
		return
	}
	err = json.Unmarshal(respBody, &task)
	if err != nil {
		log.Errorf("GetDeploymentVMs - unmarshalling tasks %v, payload %s", err, string(respBody))
		return
	}
	output := c.WaitForTaskResult(task.ID, auth)
	for _, value := range output {
		if len(value) > 0 {
			var vm VM
			err = json.Unmarshal([]byte(value), &vm)
			if err != nil {
				log.Errorf("GetDeploymentVMs - unmarshalling vms %v, payload %s", err, value)
				return
			}
			vms = append(vms, vm)
		}
	}
	return
}

// SSH from given BOSH
func (c *Client) SSH(sshRequest SSHRequest, auth Auth) (sshResponses []SSHResponse, err error) {
	var task Task
	r := c.NewRequest("POST", "/deployments/"+sshRequest.DeploymentName+"/ssh")
	jsonRequest, err := json.Marshal(sshRequest)
	if err != nil {
		log.Errorf("SSH - requesting marshal ssh %v", err)
		return
	}
	buffer := bytes.NewBufferString(string(jsonRequest))
	r.Body = buffer
	r.Header["Content-Type"] = "application/json"
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("SSH - requesting ssh %v", err)
		return
	}
	err = json.Unmarshal(respBody, &task)
	if err != nil {
		log.Errorf("SSH - unmarshalling tasks for ssh result %v", err)
		return
	}
	output := c.WaitForTaskResult(task.ID, auth)

	err = json.Unmarshal([]byte(output[0]), &sshResponses)
	if err != nil {
		log.Errorf("SSH - unmarshalling ssh response %v, payload %s", err, output[0])
		return
	}

	return
}

// GetTasks from given BOSH
func (c *Client) GetTasks(auth Auth) (tasks []Task, err error) {
	r := c.NewRequest("GET", "/tasks")
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("GetTasks - requesting tasks  %v", err)
		return
	}

	err = json.Unmarshal(respBody, &tasks)
	if err != nil {
		log.Errorf("GetTasks - unmarshalling tasks %v, payload %s", err, string(respBody))
		return
	}
	return
}

// GetRunningTasks from given BOSH
func (c *Client) GetRunningTasks(auth Auth) (tasks []Task, err error) {
	r := c.NewRequest("GET", "/tasks?state=queued,processing,cancelling&verbose=2")
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("GetRunningTasks - requesting tasks  %v", err)
		return
	}
	err = json.Unmarshal(respBody, &tasks)
	if err != nil {
		log.Errorf("GetRunningTasks - unmarshalling tasks %v, payload %s", err, string(respBody))
		return
	}
	return
}

// GetTask from given BOSH
func (c *Client) GetTask(id int, auth Auth) (task Task, err error) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID)
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("GetTask - requesting task %v", err)
		return
	}
	err = json.Unmarshal(respBody, &task)
	if err != nil {
		log.Errorf("GetTask - unmarshalling task %v, payload %s", err, string(respBody))
		return
	}
	return
}

// GetTaskResult from given BOSH
func (c *Client) GetTaskResult(id int, auth Auth) (output []string) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID+"/output?type=result")
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("GetTaskResult - requesting task %v", err)
	}
	output = strings.Split(string(respBody), "\n")
	return
}

// WaitForTaskResult from given BOSH
func (c *Client) WaitForTaskResult(id int, auth Auth) (output []string) {
	for {
		taskStatus, err := c.GetTask(id, auth)
		if err != nil {
			log.Errorf("WaitForTaskResult - getting task %v", err)
		}
		if taskStatus.State == "done" || taskStatus.State == "error" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID+"/output?type=result")
	respBody, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Errorf("WaitForTaskResult - requesting task %v", err)
	}
	output = strings.Split(string(respBody), "\n")
	return
}
