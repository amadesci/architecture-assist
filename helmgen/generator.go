// helm/generator.go
package helmgen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"diploma/recommendations"
)

// A ChartConfig represents values for chart generation.
type ChartConfig struct {
	Cluster         recommendations.ServiceCluster
	ImageMap        map[string]string // service -> Docker image
	PortMap         map[string][]int  // service -> ports
	ReplicaCount    int               // amount of replicas
	Recommendations []recommendations.Recommendation
	OutputDir       string // chart generation directory
}

// A template for Deployment.yaml with multiple containers surropt.
const deploymentTemplate = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Cluster.Name }}-deployment
  labels:
    app: {{ .Cluster.Name }}
spec:
  replicas: {{ .ReplicaCount }}
  selector:
    matchLabels:
      app: {{ .Cluster.Name }}
  template:
    metadata:
      labels:
        app: {{ .Cluster.Name }}
    spec:
      containers:
{{- range $svc := .Cluster.Services }}
        - name: {{ $svc }}
          image: {{ index $.ImageMap $svc }}
          imagePullPolicy: IfNotPresent
          ports:
{{- range $port := index $.PortMap $svc }}
            - containerPort: {{ $port }}
{{- end }}
{{- end }}
{{- if $.HasApiGateway }}
        - name: api-gateway
          image: "your-api-gateway-image:latest"
          ports:
            - containerPort: 8080
{{- end }}
{{- if $.HasCache }}
        - name: redis-cache
          image: "redis:6-alpine"
          ports:
            - containerPort: 6379
{{- end }}
`

// A template for Service.yaml to access k8s' Deployment
const serviceTemplate = `
apiVersion: v1
kind: Service
metadata:
  name: {{ .Cluster.Name }}-service
  labels:
    app: {{ .Cluster.Name }}
spec:
  selector:
    app: {{ .Cluster.Name }}
  ports:
{{- range $svc := .Cluster.Services }}
{{- range $port := index $.PortMap $svc }}
    - protocol: TCP
      port: {{ $port }}
      targetPort: {{ $port }}
{{- end }}
{{- end }}
  type: ClusterIP
`

func createChartConfig(cluster recommendations.ServiceCluster, images map[string]string, ports map[string][]int, outputDir string) ChartConfig {
	replicaCount := 1
	if cluster.Env == "prod" {
		replicaCount = 3
	} else if cluster.Env == "test" {
		replicaCount = 1
	}

	return ChartConfig{
		Cluster:         cluster,
		ImageMap:        images,
		PortMap:         ports,
		ReplicaCount:    replicaCount,
		Recommendations: cluster.Recommendations,
		OutputDir:       outputDir,
	}
}

// HasApiGateway checks for an API Gateway integration recommendation.
func (c ChartConfig) HasApiGateway() bool {
	for _, rec := range c.Recommendations {
		if rec.Type == "api-gateway" {
			return true
		}
	}
	return false
}

// HasCache checks for a Redis cache integration recommendation.
func (c ChartConfig) HasCache() bool {
	for _, rec := range c.Recommendations {
		if rec.Type == "cache" {
			return true
		}
	}
	return false
}

// GenerateHelmCharts generates charts for each cluster.
func GenerateHelmCharts(
	clusters []recommendations.ServiceCluster,
	images map[string]string,
	ports map[string][]int,
	outputBaseDir string) error {
	for _, cluster := range clusters {
		chartDir := filepath.Join(outputBaseDir, cluster.Name)
		err := os.MkdirAll(chartDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create chart dir %s: %v", chartDir, err)
		}

		config := createChartConfig(cluster, images, ports, chartDir)

		err = generateDeploymentYaml(config)
		if err != nil {
			return fmt.Errorf("failed to generate deployment for %s: %v", cluster.Name, err)
		}

		err = generateServiceYaml(config)
		if err != nil {
			return fmt.Errorf("failed to generate service for %s: %v", cluster.Name, err)
		}
	}
	return nil
}

func generateDeploymentYaml(config ChartConfig) error {
	tpl, err := template.New("deployment").Parse(deploymentTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, config)
	if err != nil {
		return err
	}
	filePath := filepath.Join(config.OutputDir, "deployment.yaml")
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func generateServiceYaml(config ChartConfig) error {
	tpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, config)
	if err != nil {
		return err
	}
	filePath := filepath.Join(config.OutputDir, "service.yaml")
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
