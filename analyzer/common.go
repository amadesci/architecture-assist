// analyzer/common.go
package analyzer

import (
	"strconv"
	"strings"
)

func detectDBConnections(env []string) []ContainerDbConnection {
	connections := make(map[string]ContainerDbConnection)

	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := parts[0], parts[1]
		if strings.HasSuffix(key, "_DB_HOST") {
			dbType := strings.TrimSuffix(key, "_DB_HOST")
			conn := connections[dbType]
			conn.Type = dbType
			conn.Host = value
			connections[dbType] = conn
		}

		if strings.HasSuffix(key, "_DB_PORT") {
			dbType := strings.TrimSuffix(key, "_DB_PORT")
			conn := connections[dbType]
			port, _ := strconv.Atoi(value)
			conn.Port = port
			connections[dbType] = conn
		}
	}

	var result []ContainerDbConnection
	for _, c := range connections {
		result = append(result, c)
	}
	return result
}
