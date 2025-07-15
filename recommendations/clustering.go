// recommendations/clustering.go
package recommendations

import (
	"fmt"
	"strings"

	"diploma/analyzer"
)

const (
	DBConnectionWeight = 5
	ExplicitDepWeight  = 3
)

// BuildClusters puts services into clusters and
// returns a slice of ServiceCluster and err=nil if success.
func BuildClusters(services []analyzer.MicroserviceInfo) ([]ServiceCluster, error) {
	envGroups := make(map[string][]analyzer.MicroserviceInfo)

	for _, svc := range services {
		env := detectEnvironment(svc)
		envGroups[env] = append(envGroups[env], svc)
	}

	var clusters []ServiceCluster

	for env, envServices := range envGroups {
		edges := buildEnvGraph(envServices, env)
		dbClusters := clusterByDB(edges, env)
		finalClusters := mergeClusters(dbClusters, edges)

		for i := range finalClusters {
			recs := generateRecommendations([]ServiceCluster{finalClusters[i]}, env)
			finalClusters[i].Recommendations = recs
		}

		clusters = append(clusters, finalClusters...)
	}

	return clusters, nil
}

func detectEnvironment(svc analyzer.MicroserviceInfo) string {
	for _, db := range svc.SharedDB {
		host := strings.ToLower(db.Host)
		switch {
		case strings.Contains(host, "-prod"):
			return "prod"
		case strings.Contains(host, "-test") || strings.Contains(host, "-stag"):
			return "test"
		}
	}
	return "default"
}

func buildEnvGraph(services []analyzer.MicroserviceInfo, env string) []DependencyEdge {
	var edges []DependencyEdge
	dbIndex := make(map[string][]string)

	for _, svc := range services {
		for _, db := range svc.SharedDB {
			key := fmt.Sprintf("%s:%d", db.Host, db.Port)
			dbIndex[key] = append(dbIndex[key], svc.Name)
		}

		for _, dep := range svc.Dependencies {
			edges = append(edges, DependencyEdge{
				Source:   svc.Name,
				Target:   dep,
				Weight:   ExplicitDepWeight,
				Env:      env,
				EdgeType: "explicit",
			})
		}
	}

	for _, services := range dbIndex {
		for i := 0; i < len(services); i++ {
			for j := i + 1; j < len(services); j++ {
				edges = append(edges, DependencyEdge{
					Source:   services[i],
					Target:   services[j],
					Weight:   DBConnectionWeight,
					Env:      env,
					EdgeType: "db",
				})
			}
		}
	}

	return edges
}

func clusterByDB(edges []DependencyEdge, env string) []ServiceCluster {
	clusters := make(map[string]*ServiceCluster)
	serviceMap := make(map[string]string)

	for _, edge := range edges {
		if edge.EdgeType == "db" {
			clusterKey := fmt.Sprintf("%s-%s", edge.Source, edge.Target)
			if _, exists := clusters[clusterKey]; !exists {
				clusters[clusterKey] = &ServiceCluster{
					Name:     fmt.Sprintf("cluster-%d", len(clusters)+1),
					Env:      env,
					Services: []string{edge.Source, edge.Target},
					SharedDBs: []SharedDB{{
						Host: strings.Split(edge.Source, "-")[0] + "-db",
						Port: 5432,
					}},
					Rationale: []string{"Shared database"},
				}
			}
			serviceMap[edge.Source] = clusterKey
			serviceMap[edge.Target] = clusterKey
		}
	}
	result := make([]ServiceCluster, 0, len(clusters))
	for _, c := range clusters {
		result = append(result, *c)
	}
	return result
}

func mergeClusters(clusters []ServiceCluster, edges []DependencyEdge) []ServiceCluster {
	serviceToCluster := make(map[string]int)
	for i, cluster := range clusters {
		for _, svc := range cluster.Services {
			serviceToCluster[svc] = i
		}
	}

	for _, edge := range edges {
		if edge.EdgeType == "explicit" {
			srcIdx, srcExists := serviceToCluster[edge.Source]
			tgtIdx, tgtExists := serviceToCluster[edge.Target]

			if srcExists && tgtExists && srcIdx != tgtIdx {
				clusters[srcIdx].Services = append(
					clusters[srcIdx].Services,
					clusters[tgtIdx].Services...,
				)
				clusters[srcIdx].Rationale = append(
					clusters[srcIdx].Rationale,
					fmt.Sprintf("Explicit dependency: %s -> %s", edge.Source, edge.Target),
				)

				clusters = append(clusters[:tgtIdx], clusters[tgtIdx+1:]...)

				for _, svc := range clusters[srcIdx].Services {
					serviceToCluster[svc] = srcIdx
				}
			}
		}
	}
	return clusters
}
