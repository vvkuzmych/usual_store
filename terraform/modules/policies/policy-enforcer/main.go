package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type PolicyEnforcer struct {
	dockerClient *client.Client
	opaServer    string
	httpClient   *http.Client
}

type OPARequest struct {
	Input map[string]interface{} `json:"input"`
}

type OPAResponse struct {
	Result map[string]interface{} `json:"result"`
}

func main() {
	opaServer := os.Getenv("OPA_SERVER")
	if opaServer == "" {
		opaServer = "http://opa-server:8181"
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal("Failed to create Docker client:", err)
	}
	defer cli.Close()

	enforcer := &PolicyEnforcer{
		dockerClient: cli,
		opaServer:    opaServer,
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}

	log.Println("Policy Enforcer started")
	log.Printf("OPA Server: %s\n", opaServer)

	// Start monitoring Docker events
	go enforcer.MonitorDockerEvents()

	// Start periodic compliance checks
	go enforcer.PeriodicComplianceCheck()

	// Start HTTP server for health checks and API
	http.HandleFunc("/health", enforcer.HealthHandler)
	http.HandleFunc("/enforce", enforcer.EnforceHandler)
	http.HandleFunc("/audit", enforcer.AuditHandler)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("HTTP server failed:", err)
	}
}

func (e *PolicyEnforcer) MonitorDockerEvents() {
	ctx := context.Background()
	events, errs := e.dockerClient.Events(ctx, types.EventsOptions{})

	for {
		select {
		case event := <-events:
			e.handleDockerEvent(event)
		case err := <-errs:
			log.Printf("Error monitoring Docker events: %v\n", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (e *PolicyEnforcer) handleDockerEvent(event types.EventMessage) {
	switch event.Action {
	case "start":
		log.Printf("Container started: %s\n", event.Actor.Attributes["name"])
		e.validateContainer(event.Actor.ID)
	case "create":
		log.Printf("Container created: %s\n", event.Actor.Attributes["name"])
	case "stop":
		log.Printf("Container stopped: %s\n", event.Actor.Attributes["name"])
	}
}

func (e *PolicyEnforcer) validateContainer(containerID string) {
	ctx := context.Background()

	// Get container details
	container, err := e.dockerClient.ContainerInspect(ctx, containerID)
	if err != nil {
		log.Printf("Failed to inspect container %s: %v\n", containerID, err)
		return
	}

	// Build policy input
	input := map[string]interface{}{
		"container_name": container.Name,
		"image":          container.Config.Image,
		"privileged":     container.HostConfig.Privileged,
		"user":           container.Config.User,
		"networks":       getNetworkNames(container.NetworkSettings),
		"env_vars":       parseEnvVars(container.Config.Env),
		"volumes":        getVolumeMounts(container.HostConfig.Binds),
		"exposed_ports":  getExposedPorts(container.Config.ExposedPorts),
	}

	// Check security policy
	if violations := e.checkPolicy("usualstore/security/security_violations", input); len(violations) > 0 {
		log.Printf("⚠️  Security violations for container %s:\n", container.Name)
		for _, v := range violations {
			log.Printf("  - %v\n", v)
		}
	}

	// Check network policy
	if allowed := e.checkPolicy("usualstore/network/allow", input); !allowed {
		log.Printf("❌ Network policy violation for container %s\n", container.Name)
	}

	// Check resource limits
	if valid := e.checkPolicy("usualstore/resources/valid_resources", input); !valid {
		log.Printf("⚠️  Resource limit recommendations for container %s\n", container.Name)
	}
}

func (e *PolicyEnforcer) checkPolicy(policyPath string, input map[string]interface{}) interface{} {
	url := fmt.Sprintf("%s/v1/data/%s", e.opaServer, policyPath)

	reqBody := OPARequest{Input: input}
	reqJSON, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("Failed to create OPA request: %v\n", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.httpClient.Do(req)
	if err != nil {
		log.Printf("Failed to query OPA: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var opaResp OPAResponse
	if err := json.NewDecoder(resp.Body).Decode(&opaResp); err != nil {
		log.Printf("Failed to decode OPA response: %v\n", err)
		return nil
	}

	return opaResp.Result["result"]
}

func (e *PolicyEnforcer) PeriodicComplianceCheck() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Running periodic compliance check...")
		e.runComplianceCheck()
	}
}

func (e *PolicyEnforcer) runComplianceCheck() {
	ctx := context.Background()
	containers, err := e.dockerClient.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Printf("Failed to list containers: %v\n", err)
		return
	}

	log.Printf("Checking %d containers for compliance\n", len(containers))

	for _, container := range containers {
		e.validateContainer(container.ID)
	}
}

func (e *PolicyEnforcer) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (e *PolicyEnforcer) EnforceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// TODO: Implement enforcement logic
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "enforced"})
}

func (e *PolicyEnforcer) AuditHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Return audit logs
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"logs": []string{}})
}

// Helper functions
func getNetworkNames(settings *types.NetworkSettings) []string {
	var networks []string
	for name := range settings.Networks {
		networks = append(networks, name)
	}
	return networks
}

func parseEnvVars(envs []string) []map[string]string {
	var result []map[string]string
	for _, env := range envs {
		// Parse KEY=VALUE format
		result = append(result, map[string]string{"raw": env})
	}
	return result
}

func getVolumeMounts(binds []string) []map[string]string {
	var result []map[string]string
	for _, bind := range binds {
		result = append(result, map[string]string{"bind": bind})
	}
	return result
}

func getExposedPorts(ports map[string]struct{}) []string {
	var result []string
	for port := range ports {
		result = append(result, port)
	}
	return result
}
