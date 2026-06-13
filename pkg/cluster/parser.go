package cluster

import (
	"errors"
	"regexp"
	"strings"
)

// Compilamos la regex una sola vez para que sea eficiente.
// Asumo: Region(4) + Zona(2) + Negocio(X) + Entorno(1) + AKS
var clusterRegex = regexp.MustCompile(`^[A-Z]{6}([A-Z]+)([A-Z])AKS$`)

// Algunos clusters usan un naming más descriptivo.
// Ejemplo: EMAZ-DEVOPS-TOOLS-P-PGA-AGENT-01 => business=PGA, environment=P
var clusterRegexAlt = regexp.MustCompile(`^EMAZ-DEVOPS-TOOLS-([A-Z])-([A-Z]+)-AGENT-[0-9]+$`)

// ParseClusterName extrae el negocio y el entorno del nombre del clúster.
// Devuelve dos strings y un error. En Go el manejo de errores es explícito.
func ParseClusterName(clusterName string) (string, string, error) {
	// Buscamos las coincidencias
	matches := clusterRegex.FindStringSubmatch(clusterName)
	if len(matches) == 3 {
		business := matches[1]
		environment := strings.ToLower(matches[2]) // Pasamos el entorno a minúscula

		return business, environment, nil
	}

	matches = clusterRegexAlt.FindStringSubmatch(clusterName)
	if len(matches) != 3 {
		return "", "", errors.New("formato de cluster invalido: " + clusterName)
	}

	business := matches[2]
	environment := strings.ToLower(matches[1]) // Pasamos el entorno a minúscula

	return business, environment, nil
}
