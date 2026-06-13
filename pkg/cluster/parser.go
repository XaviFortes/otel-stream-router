package cluster

import (
	"errors"
	"regexp"
	"strings"
)

// Regex resistente y case-insensitive para nombres descriptivos tipo:
// emaz-devops-tools-p-pga-agent-01-aks-admin
// Captura business = devops-tools, environment = p
var emazRegex = regexp.MustCompile(`(?i)^emaz-([a-z0-9-]+)-([a-z])-(?:.*-)?agent-[0-9]+(?:-.*)?$`)

// Legacy compact naming format. Example: EMEAWELABORPAKS => business=LABOR, env=P
var legacyRegex = regexp.MustCompile(`(?i)^[A-Z]{6}([A-Z]+)([A-Z])aks$`)

// ParseClusterName extrae el negocio (business) y el entorno (environment)
// del nombre del clúster. Devuelve valores en minúscula cuando procede.
// Si no se reconoce el patrón, intenta una heurística basada en tokens
// buscando un token de entorno único ('p' o 'd').
func ParseClusterName(clusterName string) (string, string, error) {
	if clusterName == "" {
		return "unknown", "unknown", errors.New("cluster name vacío")
	}

	// 1) Intentamos el patrón específico EMAZ (descriptivo)
	if matches := emazRegex.FindStringSubmatch(clusterName); len(matches) == 3 {
		business := strings.ToLower(matches[1])
		environment := strings.ToLower(matches[2])
		if business == "" {
			business = "unknown"
		}
		if environment == "" {
			environment = "unknown"
		}
		return business, environment, nil
	}

	// 1b) Intentamos el patrón legacy compacto (EMEA...)
	if matches := legacyRegex.FindStringSubmatch(clusterName); len(matches) == 3 {
		business := strings.ToLower(matches[1])
		environment := strings.ToLower(matches[2])
		if business == "" {
			business = "unknown"
		}
		if environment == "" {
			environment = "unknown"
		}
		return business, environment, nil
	}

	// 2) Heurística: dividir en tokens y buscar un token de entorno ('p' o 'd')
	toks := strings.Split(strings.ToLower(clusterName), "-")
	for i, t := range toks {
		if t == "p" || t == "d" {
			// Si el nombre empieza con 'emaz' descartamos ese prefijo
			start := 0
			if len(toks) > 0 && toks[0] == "emaz" {
				start = 1
			}
			if i <= start {
				return "unknown", t, nil
			}
			business := strings.Join(toks[start:i], "-")
			if business == "" {
				business = "unknown"
			}
			return business, t, nil
		}
	}

	return "unknown", "unknown", errors.New("formato de cluster no reconocido: " + clusterName)
}
