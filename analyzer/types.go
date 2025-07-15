// analyzer/types.go
package analyzer

// A MicroserviceInfo represents key information on a service
// for following processing.
type MicroserviceInfo struct {
	Name         string
	EnvVariables []string
	Ports        []int
	Dependencies []string
	SharedDB     []ContainerDbConnection
	SourceType   string
}

// A ContainerDbConnection serves databse connection values
type ContainerDbConnection struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
