package main

import (
	"fmt"
	"log"

	"diploma/analyzer"
	"diploma/helmgen"
	"diploma/recommendations"
)

func main() {
	projectDir := "./targets/target_act"
	outputDir := "./outputs/output_act"

	services, err := analyzer.AnalyzeDockerProject(projectDir)
	if err != nil {
		log.Fatalf("Ошибка анализа проекта: %v", err)
	}
	fmt.Printf("Микросервисов: %d\n", len(services))

	clusters, err := recommendations.BuildClusters(services)
	if err != nil {
		log.Fatalf("Ошибка построения кластеров: %v", err)
	}
	fmt.Printf("Кластеров: %d\n", len(clusters))

	images := make(map[string]string)
	ports := make(map[string][]int)
	for _, svc := range services {
		images[svc.Name] = svc.Name
		ports[svc.Name] = svc.Ports
	}

	err = helmgen.GenerateHelmCharts(clusters, images, ports, outputDir)
	if err != nil {
		log.Fatalf("Ошибка генерации Helm-чартов: %v", err)
	}
}
