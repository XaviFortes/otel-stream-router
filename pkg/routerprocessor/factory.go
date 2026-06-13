package routerprocessor

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
)

const (
	// typeStr es el nombre con el que llamarás a tu procesador en el config.yaml (ej: router_custom)
	typeStr = "router_custom"
)

// NewFactory crea la factoría para nuestro procesador.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		component.MustNewType(typeStr),
		createDefaultConfig,
		// Aquí luego le diremos si procesa Logs, Metrics o Traces.
		// De momento lo dejamos preparado.
		processor.WithLogs(createLogsProcessor, component.StabilityLevelAlpha),
		// Si queremos métricas, descomentamos esta línea y creáis la función
		// processor.WithMetrics(createMetricsProcessor, component.StabilityLevelAlpha),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		DefaultStream: "default",
	}
}

func createLogsProcessor(
	ctx context.Context,
	set processor.Settings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (processor.Logs, error) {

	oCfg := cfg.(*Config)

	// Creamos la instancia de nuestro procesador
	router := newCustomRouter(oCfg)

	// processorhelper nos ahorra escribir todo el código de arranque y parada del componente.
	// Solo le pasamos nuestra función processLogs que mutará los datos.
	return processorhelper.NewLogs(
		ctx,
		set,
		cfg,
		nextConsumer,
		router.processLogs,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
	)
}
