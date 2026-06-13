package routerprocessor

import (
	"context"
	"os"

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
		attrs := rs.Resource().Attributes()

		// 1. Intentamos sacar el namespace. El procesador k8sattributes estándar
		// de OTel inyecta este dato automáticamente antes de que llegue a nosotros.
		namespace := r.defaultStream
		if val, ok := attrs.Get("k8s.namespace.name"); ok {
			namespace = val.Str()
		}

		// (Opcional) Si tuvieras un label en el namespace llamado "proyecto",
		// k8sattributes suele inyectarlo como "k8s.namespace.labels.proyecto".
		// Podrías hacer un attrs.Get() de eso para agrupar namespaces distintos.

		// 2. Construimos el nombre del stream: ej. "pagos-p" o "namespace-unknown"
		streamName := namespace + "-" + r.env

		// 3. Inyectamos el resultado como un nuevo atributo en el log.
		// attrs.PutStr("x_custom_stream", streamName)
		attrs.PutStr("_stream_name", streamName)
	}

	return logs, nil
}
