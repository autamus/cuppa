package oras 

// Types and functions for manifests, configs, etc.

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ImageConfig struct {
	Architecture string `json:"architecture"`
	Config       struct {
		Hostname     string            `json:"Hostname"`
		Domainname   string            `json:"Domainname"`
		User         string            `json:"User"`
		AttachStdin  bool              `json:"AttachStdin"`
		AttachStdout bool              `json:"AttachStdout"`
		AttachStderr bool              `json:"AttachStderr"`
		Tty          bool              `json:"Tty"`
		OpenStdin    bool              `json:"OpenStdin"`
		StdinOnce    bool              `json:"StdinOnce"`
		Env          []string          `json:"Env"`
		Cmd          []string          `json:"Cmd"`
		Image        string            `json:"Image"`
		Volumes      interface{}       `json:"Volumes"`
		WorkingDir   string            `json:"WorkingDir"`
		Entrypoint   interface{}       `json:"Entrypoint"`
		OnBuild      interface{}       `json:"OnBuild"`
		Labels       map[string]string `json:"Labels"`
	} `json:"config"`
	Container       string `json:"container"`
	ContainerConfig struct {
		Hostname     string            `json:"Hostname"`
		Domainname   string            `json:"Domainname"`
		User         string            `json:"User"`
		AttachStdin  bool              `json:"AttachStdin"`
		AttachStdout bool              `json:"AttachStdout"`
		AttachStderr bool              `json:"AttachStderr"`
		Tty          bool              `json:"Tty"`
		OpenStdin    bool              `json:"OpenStdin"`
		StdinOnce    bool              `json:"StdinOnce"`
		Env          []string          `json:"Env"`
		Cmd          []string          `json:"Cmd"`
		Image        string            `json:"Image"`
		Volumes      interface{}       `json:"Volumes"`
		WorkingDir   string            `json:"WorkingDir"`
		Entrypoint   interface{}       `json:"Entrypoint"`
		OnBuild      interface{}       `json:"OnBuild"`
		Labels       map[string]string `json:"Labels"`
	} `json:"container_config"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	History       []struct {
		Created    time.Time `json:"created"`
		CreatedBy  string    `json:"created_by"`
		EmptyLayer bool      `json:"empty_layer,omitempty"`
	} `json:"history"`
	Os     string `json:"os"`
	Rootfs struct {
		Type    string   `json:"type"`
		DiffIds []string `json:"diff_ids"`
	} `json:"rootfs"`
}


type ImageManifest struct {
	SchemaVersion int `json:"schemaVersion"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Digest    string `json:"digest"`
		Size      int    `json:"size"`
	} `json:"config"`
	Layers []struct {
		MediaType   string `json:"mediaType"`
		Digest      string `json:"digest"`
		Size        int    `json:"size"`
		Annotations map[string]string `json:"annotations"`
	} `json:"layers"`
}

func GetRequest(url string, headers map[string]string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Read the response from the body, and return as string
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}


// GetImageManifest of an existing oras image
func GetImageManifest(container string) (ImageManifest, error) {

	url := "https://crane.ggcr.dev/manifest/" + container
	manifest := ImageManifest{}
	response, err := GetRequest(url, map[string]string{})
	if err != nil {
		return manifest, err
	}
	json.Unmarshal([]byte(response), &manifest)
	return manifest, err
}


// GetImageDigest gets a digest based on a full URI with tag
func GetImageDigest(container string) (string, error) {

	url := "https://crane.ggcr.dev/digest/" + container
	response, err := GetRequest(url, map[string]string{})
	if err != nil {
		return "", err
	}
	return string(response), nil
}


// GetImageConfig of an existing container
func GetImageConfig(container string) (ImageConfig, error) {

	// Get tags for current container image
	configUrl := "https://crane.ggcr.dev/config/" + container
	imageConf := ImageConfig{}
	response, err := GetRequest(configUrl, map[string]string{})
	if err != nil {
		return imageConf, err
	}
	json.Unmarshal([]byte(response), &imageConf)
	return imageConf, err
}

// Get image tags for a container
func GetImageTags(container string) ([]string, error) {
	tags := []string{}
	tagsUrl := "https://crane.ggcr.dev/ls/" + container
	response, err := GetRequest(tagsUrl, map[string]string{})
	if err != nil {
		return tags, err
	}
	tags = strings.Split(response, "\n")
	return tags, err
}
