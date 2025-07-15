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

// DetectEnvironment searches if the target environment contains
// docker-compose or only Dockerfiles.
// Returns 0 for docker-compose one, 1 for Dockerfile-only, 2 for other (including occured errors)
func DetectEnvironment(projectRoot string) EnvironmentType {
	if _, err := os.Stat(filepath.Join(projectRoot, "docker-compose.yml")); err == nil {
		return DockerComposeEnv
	}

	if _, err := os.Stat(filepath.Join(projectRoot, "Dockerfile")); err == nil {
		return DockerOnly
	}

	return UnknownEnv
}
