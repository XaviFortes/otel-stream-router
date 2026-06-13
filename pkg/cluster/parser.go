package cluster

import (
	"errors"
	"regexp"
	"strings"
)

// Compilamos la regex una sola vez para que sea eficiente.
// Asumo: Region(4) + Zona(2) + Negocio(X) + Entorno(1) + AKS
var clusterRegex = regexp.MustCompile(`^[A-Z]{6}([A-Z]+)([A-Z])AKS$`)

// ParseClusterName extrae el negocio y el entorno del nombre del clúster.
// Devuelve dos strings y un error. En Go el manejo de errores es explícito.
func ParseClusterName(clusterName string) (string, string, error) {
	// Buscamos las coincidencias
	matches := clusterRegex.FindStringSubmatch(clusterName)

	// Si matches no tiene 3 elementos (el string entero, el grupo 1 y el grupo 2)
	// es que el nombre no cumple la convención.
	if len(matches) != 3 {
		return "", "", errors.New("formato de cluster invalido: " + clusterName)
	}

	business := matches[1]
	environment := strings.ToLower(matches[2]) // Pasamos el entorno a minúscula

	return business, environment, nil
}
