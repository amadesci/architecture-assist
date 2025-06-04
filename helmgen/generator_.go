// // helm/generator.go
package helmgen

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"text/template"

// 	"diploma/recommendations"
// )

// const (
// 	chartTemplate = `apiVersion: v2
// name: {{.ChartName}}
// description: Helm chart for {{.Cluster.Name}}
// version: 0.1.0
// appVersion: "1.0"
// dependencies:{{range .Dependencies}}
// - name: {{.Name}}
//   version: {{.Version}}
//   repository: {{.Repo}}{{end}}
// `

// 	valuesTemplate = `global:
//   env: {{.Cluster.Env}}
//   replicaCount: {{.ReplicaCount}}

// {{if .APIEnabled}}apiGateway:
//   enabled: true
//   image: nginx:latest{{end}}

// {{if .CacheEnabled}}redis:
//   enabled: true
//   host: redis-{{.Cluster.Env}}{{end}}

// services:{{range .Cluster.Services}}
//   {{.}}:
//     image: "{{index $.ImageMap .}}"
//     ports:{{range $port := index $.PortMap .}}
//       - {{$port}}{{end}}{{end}}
// `
// )

// type ChartConfig struct {
// 	Cluster      recommendations.ServiceCluster
// 	ImageMap     map[string]string
// 	PortMap      map[string][]int
// 	APIEnabled   bool
// 	CacheEnabled bool
// 	ReplicaCount int
// 	Dependencies []Dependency
// 	Recommendations []recommendations.Recommendation // Добавлено поле для рекомендаций
// }

// type Dependency struct {
// 	Name    string
// 	Version string
// 	Repo    string
// }

// func GenerateHelmCharts(clusters []recommendations.ServiceCluster, serviceDetails map[string]struct {
// 	Image string
// 	Ports []int
// }, outputDir string) error {
// 	for _, cluster := range clusters {
// 		chartDir := filepath.Join(outputDir, cluster.Name)
// 		if err := os.MkdirAll(chartDir, 0755); err != nil {
// 			return err
// 		}

// 		config := createChartConfig(cluster, serviceDetails)

// 		if err := generateFile(chartDir, "Chart.yaml", chartTemplate, config); err != nil {
// 			return err
// 		}

// 		if err := generateFile(chartDir, "values.yaml", valuesTemplate, config); err != nil {
// 			return err
// 		}

// 		if err := generateTemplates(chartDir, config); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func createChartConfig(cluster recommendations.ServiceCluster, details map[string]struct {
// 	Image string
// 	Ports []int
// }) ChartConfig {
// 	config := ChartConfig{
// 		Cluster:      cluster,
// 		ReplicaCount: 1,
// 		ImageMap:     make(map[string]string),
// 		PortMap:      make(map[string][]int),
// 		Recommendations: cluster.Recommendations, // Заполнение рекомендаций
// 	}

// 	// Заполнение данных образов и портов
// 	for _, svc := range cluster.Services {
// 		if detail, exists := details[svc]; exists {
// 			config.ImageMap[svc] = detail.Image
// 			config.PortMap[svc] = detail.Ports
// 		}
// 	}

// 	// Настройки окружения
// 	switch cluster.Env {
// 	case "prod":
// 		config.ReplicaCount = 3
// 	case "test":
// 		config.ReplicaCount = 1
// 	}

// 	// Обработка рекомендаций
// 	for _, rec := range cluster.Recommendations {
// 		switch rec.Type {
// 		case "api-gateway":
// 			config.APIEnabled = true
// 			config.Dependencies = append(config.Dependencies, Dependency{
// 				Name:    "nginx-ingress",
// 				Version: "1.0",
// 				Repo:    "https://charts.helm.sh/stable",
// 			})
// 		case "cache":
// 			config.CacheEnabled = true
// 			config.Dependencies = append(config.Dependencies, Dependency{
// 				Name:    "redis",
// 				Version: "12.0.0",
// 				Repo:    "https://charts.bitnami.com/bitnami",
// 			})
// 		}
// 	}

// 	return config
// }

// func generateFile(dir, filename, tpl string, data interface{}) error {
// 	path := filepath.Join(dir, filename)
// 	f, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	tmpl := template.Must(template.New("").Parse(tpl))
// 	return tmpl.Execute(f, data)
// }

// func generateTemplates(chartDir string, config ChartConfig) error {
// 	templatesDir := filepath.Join(chartDir, "templates")
// 	if err := os.MkdirAll(templatesDir, 0755); err != nil {
// 		return err
// 	}

// 	// Генерация Deployment
// 	for svc, image := range config.ImageMap {
// 		deployment := fmt.Sprintf(`apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: %s
// spec:
//   replicas: {{.Values.global.replicaCount}}
//   selector:
//     matchLabels:
//       app: %s
//   template:
//     metadata:
//       labels:
//         app: %s
//     spec:
//       containers:
//       - name: %s
//         image: "{{.Values.services.%s.image}}"
//         ports:{{range .Values.services.%s.ports}}
//         - containerPort: {{.}}{{end}}
// `, svc, svc, svc, svc, svc, svc)

// 		if err := os.WriteFile(filepath.Join(templatesDir, svc+"-deployment.yaml"), []byte(deployment), 0644); err != nil {
// 			return err
// 		}
// 	}

// 	// Генерация Service
// 	service := `apiVersion: v1
// kind: Service
// metadata:
//   name: {{.Chart.Name}}-svc
// spec:
//   selector:
//     app: {{.Chart.Name}}
//   ports:{{range $svc, $ports := .Values.services}}
//   - name: {{$svc}}
//     port: {{index $ports.ports 0}}
//     targetPort: {{index $ports.ports 0}}{{end}}
// `
// 	if err := os.WriteFile(filepath.Join(templatesDir, "service.yaml"), []byte(service), 0644); err != nil {
// 		return err
// 	}

// 	return nil
// }
