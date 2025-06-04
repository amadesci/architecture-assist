// analyzer/compose.go
package analyzer

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type DockerCompose struct {
	Services map[string]struct {
		Environment map[string]string `yaml:"environment"`
		Ports       []string          `yaml:"ports"`
		DependsOn   []string          `yaml:"depends_on"`
		Build       struct {
			Context    string `yaml:"context"`
			Dockerfile string `yaml:"dockerfile"`
		} `yaml:"build"`
	} `yaml:"services"`
}

func AnalyzeDockerProject(projectRoot string) ([]MicroserviceInfo, error) {
	var services []MicroserviceInfo

	// Анализ docker-compose
	composePath := filepath.Join(projectRoot, "docker-compose.yaml")
	if _, err := os.Stat(composePath); err == nil {
		composeServices := parseComposeFile(composePath)
		services = append(services, composeServices...)
	}

	// Анализ отдельных Dockerfile
	if len(services) == 0 {
		dockerfilePath := filepath.Join(projectRoot, "Dockerfile")
		if _, err := os.Stat(dockerfilePath); err == nil {
			services = append(services, analyzeDockerfile(dockerfilePath))
		}
	}

	return services, nil
}

func parseComposeFile(path string) []MicroserviceInfo {
	data, _ := os.ReadFile(path)

	var compose DockerCompose
	yaml.Unmarshal(data, &compose)

	var services []MicroserviceInfo

	for name, svc := range compose.Services {
		info := MicroserviceInfo{
			Name:         name,
			SourceType:   "docker-compose",
			Dependencies: svc.DependsOn,
		}

		// Анализ переменных окружения
		for k, v := range svc.Environment {
			info.EnvVariables = append(info.EnvVariables, k+"="+v)
		}

		// Анализ портов
		for _, p := range svc.Ports {
			parts := strings.Split(p, ":")
			if port, err := strconv.Atoi(parts[0]); err == nil {
				info.Ports = append(info.Ports, port)
			}
		}

		// Анализ Dockerfile
		if svc.Build.Context != "" {
			dockerfilePath := filepath.Join(svc.Build.Context, svc.Build.Dockerfile)
			dfInfo := analyzeDockerfile(dockerfilePath)
			info.EnvVariables = append(info.EnvVariables, dfInfo.EnvVariables...)
			info.Ports = append(info.Ports, dfInfo.Ports...)

			// if dfInfo := analyzeDockerfile(dockerfilePath); dfInfo != (MicroserviceInfo{}) {
			//     info.EnvVariables = append(info.EnvVariables, dfInfo.EnvVariables...)
			//     info.Ports = append(info.Ports, dfInfo.Ports...)
			// }
		}

		// Обнаружение БД
		info.SharedDB = detectDBConnections(info.EnvVariables)

		services = append(services, info)
	}

	return services
}

func analyzeDockerfile(path string) MicroserviceInfo {
	data, _ := os.ReadFile(path)
	info := MicroserviceInfo{
		SourceType: "dockerfile",
	}

	// Парсинг Dockerfile
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "EXPOSE ") {
			parts := strings.Split(line, " ")
			if port, err := strconv.Atoi(parts[1]); err == nil {
				info.Ports = append(info.Ports, port)
			}
		}
		if strings.HasPrefix(line, "ENV ") {
			info.EnvVariables = append(info.EnvVariables, strings.TrimPrefix(line, "ENV "))
		}
	}

	return info
}
