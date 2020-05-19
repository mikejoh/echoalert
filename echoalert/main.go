package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Alert struct {
	Annotations  map[string]interface{} `json:"annotations"`
	EndsAt       string                 `json:"endsAt"`
	GeneratorURL string                 `json:"generatorURL"`
	Labels       map[string]interface{} `json:"labels"`
	StartsAt     string                 `json:"startsAt"`
}

type Alerts struct {
	Alerts            []Alert                `json:"alerts"`
	CommonAnnotations map[string]interface{} `json:"commonAnnotations"`
	CommonLabels      map[string]interface{} `json:"commonLabels"`
	ExternalURL       string                 `json:"externalURL"`
	GroupKey          int                    `json:"groupKey"`
	GroupLabels       map[string]interface{} `json:"groupLabels"`
	Receiver          string                 `json:"receiver"`
	Status            string                 `json:"status"`
	Version           int                    `json:"version"`
}

type ServeAlerts struct {
	ServeTemplate *template.Template
	Alerts []Alerts
	Counters map[string]int
}

type ShapeShifterResponse struct {
	Code int
	Body string
}

func NewServeAlerts() (*ServeAlerts, error) {
	sa := new(ServeAlerts)

	c := map[string]int{
		"/default": 0,
		"/super_critical": 0,
		"/team": 0,
	}

	sa.Counters = c

	t, err := template.ParseFiles("/index.tmpl")
	if err != nil {
		return sa, fmt.Errorf("error parsing template file: %s", err)
	}

	sa.ServeTemplate = t

	return sa, nil
}

func (sa *ServeAlerts) serveIndex(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := sa.ServeTemplate.Execute(w, &sa); err != nil {
		log.Error(err)
	}
}

func (sa *ServeAlerts) webhookDefault(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		sa.Counters[req.RequestURI] += 1

		var alerts Alerts

		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error(err)
		}

		json.Unmarshal(b, &alerts)

		sa.Alerts = append(sa.Alerts, alerts)
	}
}

func (sa *ServeAlerts) webhookSuperCritical(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		sa.Counters[req.RequestURI] += 1
		var alerts Alerts

		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error(err)
		}

		json.Unmarshal(b, &alerts)

		sa.Alerts = append(sa.Alerts, alerts)
	}
}

func (sa *ServeAlerts) webhookTeam(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		sa.Counters[req.RequestURI] += 1
		var alerts Alerts

		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error(err)
		}

		json.Unmarshal(b, &alerts)

		sa.Alerts = append(sa.Alerts, alerts)
	}
}

// TODO(mikejoh): Handle cases when any of the query parameters are empty
func shapeShifter(w http.ResponseWriter, req *http.Request) {
	values := req.URL.Query()

	code, err := strconv.Atoi(values["code"][0])
	if err != nil {
		log.Fatal(err)
	}
	body := values["body"][0]

	r := &http.Response{
		Status:        http.StatusText(code),
		StatusCode:    code,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
	}

	w.WriteHeader(code)

	if err := r.Write(w); err != nil {
		log.Error(err)
	}
}

//TODO(mikejoh): Implement a randomly failing (HTTP response code != 2xx) endpoint
func main() {
	sa, err := NewServeAlerts()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", sa.serveIndex)
	http.HandleFunc("/default", sa.webhookDefault)
	http.HandleFunc("/super_critical", sa.webhookSuperCritical)
	http.HandleFunc("/team", sa.webhookTeam)
	http.HandleFunc("/shapeshifter", shapeShifter)

	log.Info("Starting echoalert..")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
