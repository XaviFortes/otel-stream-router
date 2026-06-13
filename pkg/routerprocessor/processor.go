package routerprocessor

import (
	"context"
	"os"
	"strings"
	"regexp"

	"go.opentelemetry.io/collector/pdata/plog"

	// OJO: Cambia esto si tu modulo se llama distinto en go.mod
	"github.com/xavifortes/otel-stream-router/pkg/cluster"
)

type customRouter struct {
	defaultStream string
	biz           string
	env           string
}

func newCustomRouter(cfg *Config) *customRouter {
	// Leemos el nombre del clúster de las variables de entorno del Pod de K8s.
	// Si no existe, no rompemos nada, nos quedamos vacíos.
	clusterName := os.Getenv("CLUSTER_NAME")

	biz := "unknown"
	env := "unknown"

	// Llamamos a tu código de la Fase 1. Lo hacemos UNA sola vez al arrancar,
	// no por cada puto log que llega, para no freír la CPU.
	if b, e, err := cluster.ParseClusterName(clusterName); err == nil {
		biz = b
		env = e
	}

	return &customRouter{
		defaultStream: cfg.DefaultStream,
		biz:           biz,
		env:           env,
	}
}

// processLogs intercepta el tráfico de logs.
func (r *customRouter) processLogs(ctx context.Context, logs plog.Logs) (plog.Logs, error) {
	// Los datos en OTel vienen agrupados por "Resources" (ej: un Pod específico)
	for i := 0; i < logs.ResourceLogs().Len(); i++ {
		rs := logs.ResourceLogs().At(i)
		resourceAttrs := rs.Resource().Attributes()

		// Extract nodepool from resource attribute k8s.node.name (injected by Helm as OTEL_K8S_NODE_NAME)
		// AKS node names typically contain the agentpool: aks-agentpool-12345678-vmss
		nodepool := ""
		if v, ok := resourceAttrs.Get("k8s.node.name"); ok {
			nodeName := v.Str()
			nodepool = extractNodepoolFromNodeName(nodeName)
		}

		// Iterate all scope logs / log records to prefer record-level namespace if present
		// and build the stream name per-record. We write the computed value into the
		// Resource attributes under _stream_name so downstream exporters that rely on
		// a resource-scoped attribute (like OpenObserve routing) can consume it.
		for si := 0; si < rs.ScopeLogs().Len(); si++ {
			sl := rs.ScopeLogs().At(si)
			for ri := 0; ri < sl.LogRecords().Len(); ri++ {
				lr := sl.LogRecords().At(ri)

				// Prefer namespace on the log record (k8s receivers may attach at record level)
				namespace := r.defaultStream
				if val, ok := lr.Attributes().Get("k8s.namespace.name"); ok {
					namespace = val.Str()
				} else if val, ok := resourceAttrs.Get("k8s.namespace.name"); ok {
					namespace = val.Str()
				}

				// Construct stream name with nodepool included
				streamName := buildStreamName(r.defaultStream, namespace, r.biz, r.env, nodepool)

				// Inject as a LogRecord attribute for per-record routing
				lr.Attributes().PutStr("_stream_name", streamName)
			}
		}
	}

	return logs, nil
}

func buildStreamName(defaultStream, namespace, business, environment, nodepool string) string {
	parts := make([]string, 0, 4)

	if namespace != "" {
		parts = append(parts, namespace)
	} else if defaultStream != "" {
		parts = append(parts, defaultStream)
	}

	if business != "" && business != "unknown" {
		parts = append(parts, business)
	}

	if environment != "" && environment != "unknown" {
		parts = append(parts, strings.ToLower(environment))
	}

	if nodepool != "" && nodepool != "unknown" {
		parts = append(parts, nodepool)
	}

	// If nothing meaningful was appended, fall back to defaultStream explicitly
	if len(parts) == 0 && defaultStream != "" {
		return defaultStream
	}

	return strings.Join(parts, "-")
}

// extractNodepoolFromNodeName derives a nodepool/agentpool identifier from an AKS node name.
// Example: aks-agentpool-12345678-vmss -> agentpool
func extractNodepoolFromNodeName(nodeName string) string {
	if nodeName == "" {
		return ""
	}
	// Pattern: aks-<agentpool>-<rest>
	// Capture the agentpool token between the first and second dash when prefixed with aks-
	// Case-insensitive to be safe.
	var nodeRegex = regexp.MustCompile(`(?i)^aks-([a-z0-9-]+)-`)
	if matches := nodeRegex.FindStringSubmatch(nodeName); len(matches) == 2 {
		return matches[1]
	}
	return nodeName
}
