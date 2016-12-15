package bosh

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GetStemcells from given BOSH
func (c *Client) GetStemcells(auth Auth) (stemcells []Stemcell, err error) {
	r := c.NewRequest("GET", "/stemcells")
	resp, err := c.DoAuthRequest(r, auth)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("Error requesting stemcells  %v", err)
		return
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading stemcells request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &stemcells)
	if err != nil {
		log.Printf("Error unmarshaling stemcells %v", err)
		return
	}
	return
}

// GetReleases from given BOSH
func (c *Client) GetReleases(auth Auth) (releases []Release, err error) {
	r := c.NewRequest("GET", "/releases")
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting releases  %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading releases request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &releases)
	if err != nil {
		log.Printf("Error unmarshaling releases %v", err)
		return
	}
	return
}

// GetDeployments from given BOSH
func (c *Client) GetDeployments(auth Auth) (deployments []Deployment, err error) {
	r := c.NewRequest("GET", "/deployments")
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting deployments  %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployments request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &deployments)
	if err != nil {
		log.Printf("Error unmarshaling deployments %v", err)
		return
	}
	return
}

// GetDeployment from given BOSH
func (c *Client) GetDeployment(name string, auth Auth) (manifest Manifest, err error) {
	r := c.NewRequest("GET", "/deployments/"+name)
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting deployment manifest %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployment manifest request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &manifest)
	if err != nil {
		log.Printf("Error unmarshaling deployment manifest %v", err)
		return
	}
	return
}

// DeleteDeployment from given BOSH
func (c *Client) DeleteDeployment(name string, auth Auth) (task Task, err error) {
	r := c.NewRequest("DELETE", "/deployments/"+name+"?force=true")
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting deleting deployment %v", err)
		return
	}
	defer resp.Body.Close()
	url, _ := resp.Location()
	re, _ := regexp.Compile(`(\d+)$`)
	stringId := re.FindStringSubmatch(url.Path)
	id, err := strconv.Atoi(stringId[0])
	if err != nil {
		return
	}
	task, err = c.GetTask(id, auth)
	return
}

// CreateDeployment from given BOSH
func (c *Client) CreateDeployment(manifest string, auth Auth) (task Task, err error) {
	r := c.NewRequest("POST", "/deployments")
	buffer := bytes.NewBufferString(manifest)
	r.Body = buffer
	r.Header["Content-Type"] = "text/yaml"

	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting create deployment %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
		return
	}
	return
}

// GetDeploymentVMs from given BOSH
func (c *Client) GetDeploymentVMs(name string, auth Auth) (vms []VM, err error) {
	var task Task
	r := c.NewRequest("GET", "/deployments/"+name+"/vms?format=full")
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting deployment vms %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading deployment vms request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling tasks %v", err)
		return
	}
	for {
		taskStatus, err := c.GetTask(task.ID, auth)
		if err != nil {
			log.Printf("Error getting task %v", err)
		}
		if taskStatus.State == "done" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	output := c.GetTaskResult(task.ID, auth)
	for _, value := range output {
		if len(value) > 0 {
			var vm VM
			err = json.Unmarshal([]byte(value), &vm)
			if err != nil {
				log.Printf("Error unmarshaling vms %v %v", value, err)
				return
			}
			vms = append(vms, vm)
		}
	}
	return
}

// GetTasks from given BOSH
func (c *Client) GetTasks(auth Auth) (tasks []Task, err error) {
	r := c.NewRequest("GET", "/tasks")
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting tasks  %v", err)
		return
	}

	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading tasks request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &tasks)
	if err != nil {
		log.Printf("Error unmarshaling tasks %v", err)
		return
	}
	return
}

// GetTask from given BOSH
func (c *Client) GetTask(id int, auth Auth) (task Task, err error) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID)
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting task %v", err)
		return
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
		return
	}
	err = json.Unmarshal(resBody, &task)
	if err != nil {
		log.Printf("Error unmarshaling task %v", err)
		return
	}
	return
}

// GetTaskResult from given BOSH
func (c *Client) GetTaskResult(id int, auth Auth) (output []string) {
	stringID := strconv.Itoa(id)
	r := c.NewRequest("GET", "/tasks/"+stringID+"/output?type=result")
	resp, err := c.DoAuthRequest(r, auth)

	if err != nil {
		log.Printf("Error requesting task %v", err)
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading task request %v", resBody)
	}
	output = strings.Split(string(resBody), "\n")
	return
}
