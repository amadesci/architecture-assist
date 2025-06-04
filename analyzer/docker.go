// analyzer/docker.go
package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ParseDockerfile анализирует Dockerfile для выявления ключевых характеристик
func ParseDockerfile(path string) (MicroserviceInfo, error) {
	info := MicroserviceInfo{
		Name:       filepath.Base(filepath.Dir(path)),
		SourceType: "dockerfile",
	}

	file, err := os.Open(path)
	if err != nil {
		return info, fmt.Errorf("error opening Dockerfile: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		instruction := strings.ToUpper(parts[0])
		args := parts[1:]

		switch instruction {
		case "EXPOSE":
			info.processExpose(args)
		case "ENV":
			info.processEnv(args)
		}
	}

	// Определяем общие ресурсы по переменным окружения
	info.detectSharedResources()

	return info, scanner.Err()
}

func (mi *MicroserviceInfo) processExpose(args []string) {
	for _, portStr := range args {
		port, err := strconv.Atoi(portStr)
		if err == nil && port > 0 && port < 65535 {
			mi.Ports = append(mi.Ports, port)
		}
	}
}

func (mi *MicroserviceInfo) processEnv(args []string) {
	if len(args) >= 2 {
		// Объединяем все аргументы после первого как значение
		value := strings.Join(args[1:], " ")
		mi.EnvVariables = append(mi.EnvVariables, args[0]+"="+value)
	}
}

func (mi *MicroserviceInfo) detectSharedResources() {
	dbPrefixes := make(map[string]map[string]string)

	// Сбор сырых данных из переменных окружения
	for _, env := range mi.EnvVariables {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		keyParts := strings.Split(key, "_")
		if len(keyParts) < 2 {
			continue
		}

		// Обработка БД: TYPE_HOST, TYPE_PORT
		if strings.HasSuffix(key, "_HOST") {
			dbType := strings.ToLower(strings.TrimSuffix(key, "_HOST"))
			if dbPrefixes[dbType] == nil {
				dbPrefixes[dbType] = make(map[string]string)
			}
			dbPrefixes[dbType]["host"] = value
		}

		if strings.HasSuffix(key, "_PORT") {
			dbType := strings.ToLower(strings.TrimSuffix(key, "_PORT"))
			if dbPrefixes[dbType] == nil {
				dbPrefixes[dbType] = make(map[string]string)
			}
			dbPrefixes[dbType]["port"] = value
		}
	}

	// Формирование подключений к БД
	for dbType, data := range dbPrefixes {
		if data["host"] == "" {
			continue
		}

		port := 0
		if p, err := strconv.Atoi(data["port"]); err == nil {
			port = p
		}

		// Определение порта по умолчанию
		if port == 0 {
			switch dbType {
			case "postgres":
				port = 5432
			case "mysql":
				port = 3306
			}
		}

		mi.SharedDB = append(mi.SharedDB, ContainerDbConnection{
			Type: dbType,
			Host: data["host"],
			Port: port,
		})
	}

	// Формирование списка кэшей
	// for cacheType := range cacheTypes {
	// 	mi.SharedCache = append(mi.SharedCache, cacheType)
	// }
}
