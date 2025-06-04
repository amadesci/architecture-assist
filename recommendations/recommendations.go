// recommendations/recommendations.go
package recommendations

import "fmt"

func generateRecommendations(clusters []ServiceCluster, env string) []Recommendation {
	var recs []Recommendation

	for _, cluster := range clusters {
		// Рекомендация API Gateway
		if len(cluster.Services) >= 3 {
			recs = append(recs, Recommendation{
				Type:        "api-gateway",
				Description: fmt.Sprintf("API Gateway для %s (%d сервисов)", cluster.Name, len(cluster.Services)),
				Env:         env,
				Priority:    2,
			})
		}

		// Рекомендация кэша
		if len(cluster.SharedDBs) > 0 && len(cluster.Services) > 1 {
			recs = append(recs, Recommendation{
				Type:        "cache",
				Description: fmt.Sprintf("Кэш Redis для кластера %s", cluster.Name),
				Env:         env,
				Priority:    3,
			})
		}
	}

	return recs
}
