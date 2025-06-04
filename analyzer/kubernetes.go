// // analyzer/kubernetes.go
package analyzer

// import (
// 	"regexp"
// 	"strconv"
// 	"strings"

// 	"gopkg.in/yaml.v3"
// )

// // Базовые структуры ========================================================

// type metadata struct {
// 	Name   string            `yaml:"name"`
// 	Labels map[string]string `yaml:"labels"`
// }

// // Deployment обработчик ====================================================
// //
// //	type Deployment struct {
// //		Metadata metadata `yaml:"metadata"`
// //		Spec     struct {
// //			Template struct {
// //				Spec struct {
// //					Containers []struct {
// //						Env []struct {
// //							Name  string `yaml:"name"`
// //							Value string `yaml:"value"`
// //						} `yaml:"env"`
// //					} `yaml:"containers"`
// //				} `yaml:"spec"`
// //			} `yaml:"template"`
// //		} `yaml:"spec"`
// //	}
// //
// // Deployment analysis
// type Deployment struct {
// 	Metadata struct {
// 		Name string `yaml:"name"`
// 	} `yaml:"metadata"`
// 	Spec struct {
// 		Template struct {
// 			Spec struct {
// 				Containers []struct {
// 					Env []struct {
// 						Name  string `yaml:"name"`
// 						Value string `yaml:"value"`
// 					} `yaml:"env"`
// 					Ports []struct {
// 						ContainerPort int `yaml:"containerPort"`
// 					} `yaml:"ports"`
// 				} `yaml:"containers"`
// 			} `yaml:"spec"`
// 		} `yaml:"template"`
// 	} `yaml:"spec"`
// }

// func parseDeployment(content []byte) (K8sManifest, error) {
// 	var d Deployment
// 	err := yaml.Unmarshal(content, &d)
// 	return &d, err
// }

// func (d *Deployment) ResourceType() K8sResourceType { return DeploymentRes }
// func (d *Deployment) GetName() string               { return d.Metadata.Name }

// func (d *Deployment) ServiceDependencies() []string {
// 	serviceRegex := regexp.MustCompile(`([a-zA-Z0-9-]+?)(?:-svc|\.svc|\.service)`)
// 	unique := make(map[string]struct{})

// 	for _, container := range d.Spec.Template.Spec.Containers {
// 		for _, env := range container.Env {
// 			matches := serviceRegex.FindAllStringSubmatch(env.Value, -1)
// 			for _, m := range matches {
// 				if len(m) > 1 {
// 					unique[m[1]] = struct{}{}
// 				}
// 			}
// 		}
// 	}

// 	return mapKeys(unique)
// }

// func (d *Deployment) DatabaseConnections() []ContainerDbConnection {
// 	var connections []ContainerDbConnection
// 	dbVars := make(map[string]map[string]string) // key: dbType

// 	for _, container := range d.Spec.Template.Spec.Containers {
// 		for _, env := range container.Env {
// 			// Обработка переменных БД
// 			if strings.HasSuffix(env.Name, "_DB_HOST") {
// 				dbType := strings.TrimSuffix(env.Name, "_DB_HOST")
// 				if dbVars[dbType] == nil {
// 					dbVars[dbType] = make(map[string]string)
// 				}
// 				dbVars[dbType]["host"] = env.Value
// 			}

// 			if strings.HasSuffix(env.Name, "_DB_PORT") {
// 				dbType := strings.TrimSuffix(env.Name, "_DB_PORT")
// 				if dbVars[dbType] == nil {
// 					dbVars[dbType] = make(map[string]string)
// 				}
// 				dbVars[dbType]["port"] = env.Value
// 			}
// 		}
// 	}

// 	for dbType, data := range dbVars {
// 		conn := ContainerDbConnection{
// 			Type: strings.ToLower(dbType),
// 			Host: data["host"],
// 		}

// 		if port, err := strconv.Atoi(data["port"]); err == nil {
// 			conn.Port = port
// 		} else {
// 			conn.Port = getDefaultPort(conn.Type)
// 		}

// 		if conn.Host != "" {
// 			connections = append(connections, conn)
// 		}
// 	}

// 	return connections
// }

// func (d *Deployment) CacheConnections() []string {
// 	caches := make(map[string]struct{})

// 	for _, container := range d.Spec.Template.Spec.Containers {
// 		for _, env := range container.Env {
// 			if strings.HasPrefix(env.Name, "CACHE_") {
// 				cacheType := strings.TrimPrefix(env.Name, "CACHE_")
// 				caches[strings.ToLower(cacheType)] = struct{}{}
// 			}
// 		}
// 	}

// 	return mapKeys(caches)
// }

// // Service обработчик =======================================================
// type Service struct {
// 	Metadata metadata `yaml:"metadata"`
// 	Spec     struct {
// 		Selector map[string]string `yaml:"selector"`
// 	} `yaml:"spec"`
// }

// func parseService(content []byte) (K8sManifest, error) {
// 	var s Service
// 	err := yaml.Unmarshal(content, &s)
// 	return &s, err
// }

// func (s *Service) ResourceType() K8sResourceType { return ServiceRes }
// func (s *Service) GetName() string               { return s.Metadata.Name }
// func (s *Service) ServiceDependencies() []string {
// 	if app, exists := s.Spec.Selector["app"]; exists {
// 		return []string{app}
// 	}
// 	return nil
// }
// func (s *Service) DatabaseConnections() []ContainerDbConnection { return nil }
// func (s *Service) CacheConnections() []string                   { return nil }

// // PersistentVolumeClaim обработчик =========================================
// type PersistentVolumeClaim struct {
// 	Metadata metadata `yaml:"metadata"`
// }

// func parsePVC(content []byte) (K8sManifest, error) {
// 	var p PersistentVolumeClaim
// 	err := yaml.Unmarshal(content, &p)
// 	return &p, err
// }

// func (p *PersistentVolumeClaim) ResourceType() K8sResourceType { return PersistentVolumeClaimRes }
// func (p *PersistentVolumeClaim) GetName() string               { return p.Metadata.Name }
// func (p *PersistentVolumeClaim) ServiceDependencies() []string { return nil }
// func (p *PersistentVolumeClaim) DatabaseConnections() []ContainerDbConnection {
// 	if dbType, ok := p.Metadata.Labels["db-type"]; ok {
// 		return []ContainerDbConnection{{Type: dbType}}
// 	}
// 	return nil
// }
// func (p *PersistentVolumeClaim) CacheConnections() []string { return nil }

// // ConfigMap обработчик =====================================================
// type ConfigMap struct {
// 	Metadata   metadata          `yaml:"metadata"`
// 	Data       map[string]string `yaml:"data"`
// 	BinaryData map[string][]byte `yaml:"binaryData"`
// }

// func parseConfigMap(content []byte) (K8sManifest, error) {
// 	var c ConfigMap
// 	err := yaml.Unmarshal(content, &c)
// 	return &c, err
// }

// func (c *ConfigMap) ResourceType() K8sResourceType { return ConfigMapRes }
// func (c *ConfigMap) GetName() string               { return c.Metadata.Name }
// func (c *ConfigMap) ServiceDependencies() []string { return nil }
// func (c *ConfigMap) DatabaseConnections() []ContainerDbConnection {
// 	var conns []ContainerDbConnection
// 	pattern := regexp.MustCompile(`(?i)(postgresql|mysql|mongodb)://([^:]+):(\d+)`)

// 	for _, v := range c.Data {
// 		matches := pattern.FindAllStringSubmatch(v, -1)
// 		for _, m := range matches {
// 			port, _ := strconv.Atoi(m[3])
// 			conns = append(conns, ContainerDbConnection{
// 				Type: strings.ToLower(m[1]),
// 				Host: m[2],
// 				Port: port,
// 			})
// 		}
// 	}
// 	return conns
// }

// func (c *ConfigMap) CacheConnections() []string {
// 	caches := make(map[string]struct{})
// 	patterns := map[string]*regexp.Regexp{
// 		"redis":     regexp.MustCompile(`redis://`),
// 		"memcached": regexp.MustCompile(`memcached://`),
// 	}

// 	for _, v := range c.Data {
// 		for cacheType, re := range patterns {
// 			if re.MatchString(v) {
// 				caches[cacheType] = struct{}{}
// 			}
// 		}
// 	}
// 	return mapKeys(caches)
// }

// // Вспомогательные функции ==================================================
// func getDefaultPort(dbType string) int {
// 	switch strings.ToLower(dbType) {
// 	case "postgres", "postgresql":
// 		return 5432
// 	case "mysql":
// 		return 3306
// 	case "redis":
// 		return 6379
// 	default:
// 		return 0
// 	}
// }

// func mapKeys(m map[string]struct{}) []string {
// 	keys := make([]string, 0, len(m))
// 	for k := range m {
// 		keys = append(keys, k)
// 	}
// 	return keys
// }

// // ParseK8sManifest главная функция парсинга
// func ParseK8sManifest(content []byte) (K8sManifest, error) {
// 	var base struct {
// 		Kind string `yaml:"kind"`
// 	}
// 	if err := yaml.Unmarshal(content, &base); err != nil {
// 		return nil, err
// 	}

// 	switch strings.ToLower(base.Kind) {
// 	case "deployment":
// 		return parseDeployment(content)
// 	case "service":
// 		return parseService(content)
// 	case "persistentvolumeclaim":
// 		return parsePVC(content)
// 	case "configmap":
// 		return parseConfigMap(content)
// 	default:
// 		return nil, nil
// 	}
// }
