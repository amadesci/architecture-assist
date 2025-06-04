package main

import (
	"fmt"
	"log"

	"diploma/analyzer"
	"diploma/helmgen"
	"diploma/recommendations"
)

func main() {
	// Пути к проекту и выводу
	projectDir := "./targets/target_act" // Путь к папке с docker-compose и Dockerfile
	outputDir := "./outputs/output_act"  // Путь для генерации Helm-чартов

	// Шаг 1: Анализ микросервисов
	fmt.Println("Анализ микросервисов из проекта:", projectDir)
	services, err := analyzer.AnalyzeDockerProject(projectDir)
	if err != nil {
		log.Fatalf("Ошибка анализа проекта: %v", err)
	}
	fmt.Printf("Найдено микросервисов: %d\n", len(services))

	// Шаг 2: Построение кластеров и генерация рекомендаций
	fmt.Println("Формирование кластеров и генерация рекомендаций...")
	clusters, err := recommendations.BuildClusters(services)
	if err != nil {
		log.Fatalf("Ошибка построения кластеров: %v", err)
	}
	fmt.Printf("Сформировано кластеров: %d\n", len(clusters))

	// Получение всех образов и портов сервисов для генерации Helm-чартов
	images := make(map[string]string)
	ports := make(map[string][]int)
	for _, svc := range services {
		images[svc.Name] = svc.Name
		ports[svc.Name] = svc.Ports
	}

	// Шаг 3: Генерация Helm-чартов
	fmt.Println("Генерация Helm-чартов в:", outputDir)
	err = helmgen.GenerateHelmCharts(clusters, images, ports, outputDir)
	if err != nil {
		log.Fatalf("Ошибка генерации Helm-чартов: %v", err)
	}

	fmt.Println("Генерация Helm-чартов завершена успешно.")
}

// Дополнительно в analyzer реализуйте функцию AnalyzeProject для анализа docker-compose/dockerfile из папки:
// func AnalyzeProject(path string) ([]MicroserviceInfo, error)
