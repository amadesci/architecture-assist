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

func BuildClusters(services []analyzer.MicroserviceInfo) ([]ServiceCluster, error) {
	fmt.Printf("BuildClusters: \n")
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

		// Генерация рекомендаций для каждого кластера
		for i := range finalClusters {
			recs := generateRecommendations([]ServiceCluster{finalClusters[i]}, env)
			finalClusters[i].Recommendations = recs // Добавление рекомендаций в кластер
		}

		clusters = append(clusters, finalClusters...)
		//recs := generateRecommendations(finalClusters, env)
		//recommendations = append(recommendations, recs...)
	}

	fmt.Printf("/close BuildClusters: \n")
	return clusters, nil
}

func detectEnvironment(svc analyzer.MicroserviceInfo) string {
	fmt.Printf("detectEnvironment: \n")
	for _, db := range svc.SharedDB {
		host := strings.ToLower(db.Host)
		switch {
		case strings.Contains(host, "-prod"):
			return "prod"
		case strings.Contains(host, "-test") || strings.Contains(host, "-stag"):
			return "test"
		}
	}
	fmt.Printf("/close detectEnvironment: \n")
	return "default"
}

func buildEnvGraph(services []analyzer.MicroserviceInfo, env string) []DependencyEdge {
	fmt.Printf("buildEnvGraph: \n")
	var edges []DependencyEdge
	dbIndex := make(map[string][]string)

	// Обработка общих БД
	for _, svc := range services {
		fmt.Printf("Сервис: " + svc.Name + "\n")
		for _, db := range svc.SharedDB {
			fmt.Printf("   БД: " + db.Host + "\n")
			key := fmt.Sprintf("%s:%d", db.Host, db.Port)
			dbIndex[key] = append(dbIndex[key], svc.Name)
		}

		for _, dep := range svc.Dependencies {
			fmt.Printf("   Связь: " + dep + "\n")
			edges = append(edges, DependencyEdge{
				Source:   svc.Name,
				Target:   dep,
				Weight:   ExplicitDepWeight,
				Env:      env,
				EdgeType: "explicit",
			})
		}
	}

	// Добавление связей через БД
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

	// Явные зависимости из docker-compose
	// for _, svc := range services {
	// 	for _, dep := range svc.Dependencies {
	// 		edges = append(edges, DependencyEdge{
	// 			Source:   svc.Name,
	// 			Target:   dep,
	// 			Weight:   ExplicitDepWeight,
	// 			Env:      env,
	// 			EdgeType: "explicit",
	// 		})
	// 	}
	// }

	fmt.Printf("/close buildEnvGraph: \n")
	return edges
}

func clusterByDB(edges []DependencyEdge, env string) []ServiceCluster {
	fmt.Printf("clusterByDB: \n")
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
	fmt.Printf("/close clusterByDB: \n")
	return result
}

func mergeClusters(clusters []ServiceCluster, edges []DependencyEdge) []ServiceCluster {
	fmt.Printf("mergeClusters: \n")
	serviceToCluster := make(map[string]int)
	for i, cluster := range clusters {
		for _, svc := range cluster.Services {
			fmt.Printf("    for clusters\n")
			serviceToCluster[svc] = i
		}
	}

	for _, edge := range edges {
		if edge.EdgeType == "explicit" {
			fmt.Printf("    for edges\n")
			srcIdx, srcExists := serviceToCluster[edge.Source]
			tgtIdx, tgtExists := serviceToCluster[edge.Target]

			if srcExists && tgtExists && srcIdx != tgtIdx {
				// Объединение кластеров
				clusters[srcIdx].Services = append(
					clusters[srcIdx].Services,
					clusters[tgtIdx].Services...,
				)
				clusters[srcIdx].Rationale = append(
					clusters[srcIdx].Rationale,
					fmt.Sprintf("Explicit dependency: %s -> %s", edge.Source, edge.Target),
				)

				// Удаление объединённого кластера
				clusters = append(clusters[:tgtIdx], clusters[tgtIdx+1:]...)

				// Обновление индекса
				for _, svc := range clusters[srcIdx].Services {
					serviceToCluster[svc] = srcIdx
				}
			}
		}
	}

	fmt.Printf("/close mergeClusters: \n")
	return clusters
}
