// analyzer/environment.go
package analyzer

import (
	"os"
	"path/filepath"
)

type EnvironmentType int

const (
	DockerComposeEnv EnvironmentType = iota
	DockerOnly
	UnknownEnv
)

func DetectEnvironment(projectRoot string) EnvironmentType {
	// Проверяем docker-compose
	if _, err := os.Stat(filepath.Join(projectRoot, "docker-compose.yml")); err == nil {
		return DockerComposeEnv
	}

	// Проверяем Dockerfile
	if _, err := os.Stat(filepath.Join(projectRoot, "Dockerfile")); err == nil {
		return DockerOnly
	}

	return UnknownEnv
}
