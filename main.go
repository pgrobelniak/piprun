package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Execution struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

func request(method, url, bearer, data string, execution *Execution) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(execution)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		workspace := r.URL.Query().Get("workspace")
		project := r.URL.Query().Get("project")
		pipeline := r.URL.Query().Get("pipeline")
		bearer := "Bearer " + r.URL.Query().Get("token")
		execution := &Execution{}
		url := fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/pipelines/%s/executions", workspace, project, pipeline)
		err := request("POST", url, bearer, "{\"to_revision\":{\"revision\":\"HEAD\"}}", execution)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url = fmt.Sprintf("https://api.buddy.works/workspaces/%s/projects/%s/pipelines/%s/executions/%d", workspace, project, pipeline, execution.Id)
		for execution.Status == "INPROGRESS" || execution.Status == "ENQUEUED" {
			time.Sleep(time.Second)
			err = request("GET", url, bearer, "", execution)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		w.Write([]byte(execution.Status))
	})
	http.ListenAndServe("localhost:80", nil)
}
