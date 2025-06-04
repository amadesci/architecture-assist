// analyzer/types.go
package analyzer

type ContainerDbConnection struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MicroserviceInfo struct {
	Name         string
	EnvVariables []string
	Ports        []int
	Dependencies []string // Из docker-compose depends_on
	SharedDB     []ContainerDbConnection
	SourceType   string
}
