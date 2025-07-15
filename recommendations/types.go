// recommendations/types.go
package recommendations

// A ServiceCluster represents the contents and charachteristics.
// of the cluster.
type ServiceCluster struct {
	Name            string
	Services        []string
	SharedDBs       []SharedDB
	Env             string
	Rationale       []string
	Recommendations []Recommendation
}

// A SharedDB serves database description.
type SharedDB struct {
	Type string
	Host string
	Port int
}

// A DependencyEdge serves values of the service graph edge.
type DependencyEdge struct {
	Source   string
	Target   string
	Weight   int
	Env      string
	EdgeType string // "db" | "explicit"
}

// A Recommendation serves and additional recommendation for the cluster
// (e.g. chache, gateway integration).
type Recommendation struct {
	Type        string
	Description string
	Env         string
	Priority    int
}
