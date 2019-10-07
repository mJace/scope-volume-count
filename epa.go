package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// Plugin groups the methods a plugin needs
type Plugin struct {
	lock       sync.Mutex
}

type pluginSpec struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Description string   `json:"description,omitempty"`
	Interfaces  []string `json:"interfaces"`
	APIVersion  string   `json:"api_version,omitempty"`
}

type node struct {
	Node        string      `json:"nodes"`
}

type metadataTemplate struct {
	ID       string  `json:"id"`
	Label    string  `json:"label,omitempty"`
	Format   string  `json:"format,omitempty"`
	Priority float64 `json:"priority,omitempty"`
}

type topology struct {
	Nodes           	string           		`json:"nodes"`
	MetadataTemplate 	metadataTemplate 	`json:"metadata_templates"`
}

type report struct {
	Container 	topology
	Plugins 	[]pluginSpec
}

func (p *Plugin) metadataTemplates() metadataTemplate {
	return metadataTemplate{
			ID:       "affinity",
			Label:    "CPU Affinity",
			Format:   "latest",
			Priority: 0.1,
	}
}

func (p *Plugin) nodes() string {
	return "hi"
}

func (p *Plugin) makeReport() (*report, error) {
	rpt := &report{
		Container: topology{
			Nodes: p.nodes(),
			MetadataTemplate: p.metadataTemplates(),
		},
		Plugins: []pluginSpec{
			{
				ID:          "epa",
				Label:       "epa",
				Description: "Show EPA information for container",
				Interfaces:  []string{"reporter"},
				APIVersion:  "1",
			},
		},
	}
	return rpt, nil
}

func (p *Plugin) Report(w http.ResponseWriter, r *http.Request) {
	p.lock.Lock()
	defer p.lock.Unlock()
	log.Println(r.URL.String())
	rpt, err := p.makeReport()
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	raw, err := json.Marshal(*rpt)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(raw); err != nil {
		log.Fatalln("Failed to write back report")
		return
	}
}