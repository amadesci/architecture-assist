// recommendations/types.go
package recommendations

type ServiceCluster struct {
	Name            string
	Services        []string
	SharedDBs       []SharedDB
	Env             string
	Rationale       []string
	Recommendations []Recommendation
}

type SharedDB struct {
	Type string
	Host string
	Port int
}

type DependencyEdge struct {
	Source   string
	Target   string
	Weight   int
	Env      string
	EdgeType string // "db" | "explicit"
}

type Recommendation struct {
	Type        string
	Description string
	Env         string
	Priority    int
}
