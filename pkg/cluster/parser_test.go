package cluster

import (
	"testing"
)

func TestParseClusterName(t *testing.T) {
	// Definimos nuestra tabla de pruebas
	tests := []struct {
		name         string // Nombre del test
		clusterInput string // Lo que le pasamos a la función
		expectedBiz  string // Negocio esperado
		expectedEnv  string // Entorno esperado
		expectError  bool   // ¿Debería dar error?
	}{
		{
			name:         "Cluster de Produccion Labor",
			clusterInput: "EMEAWELABORPAKS",
			expectedBiz:  "LABOR",
			expectedEnv:  "p", // Lo normalizamos a minúscula
			expectError:  false,
		},
		{
			name:         "Cluster de Desarrollo Logistica",
			clusterInput: "EMEAWELOGISTICSDAKS",
			expectedBiz:  "LOGISTICS",
			expectedEnv:  "d",
			expectError:  false,
		},
		{
			name:         "Cluster con nombre entero",
			clusterInput: "EMEAWEGPAHAKS-01-admin",
			expectedBiz:  "GPA",
			expectedEnv:  "h",
			expectError:  true,
		},
		{
			name:         "Cluster con formato basura",
			clusterInput: "CLUSTER_MIO_DE_PRUEBAS",
			expectedBiz:  "",
			expectedEnv:  "",
			expectError:  true,
		},
		{
			name:         "Cluster con naming descriptivo",
			clusterInput: "EMAZ-DEVOPS-TOOLS-P-GPA-AGENT-01",
			expectedBiz:  "GPA",
			expectedEnv:  "p",
			expectError:  false,
		},
	}

	// Ejecutamos cada caso
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			biz, env, err := ParseClusterName(tc.clusterInput)

			// Comprobamos si esperábamos un error y no ha ocurrido (o viceversa)
			if (err != nil) != tc.expectError {
				t.Fatalf("Esperaba error: %v, pero obtuve: %v", tc.expectError, err)
			}

			// Si no esperamos error, comprobamos que los valores coinciden
			if !tc.expectError {
				if biz != tc.expectedBiz {
					t.Errorf("Esperaba negocio %q, obtuve %q", tc.expectedBiz, biz)
				}
				if env != tc.expectedEnv {
					t.Errorf("Esperaba entorno %q, obtuve %q", tc.expectedEnv, env)
				}
			}
		})
	}
}
