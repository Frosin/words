package metrics

const metricsPrefix = ""

var service = NewMetrics().WithNamespace(metricsPrefix)

var (
	WordsOperationResults = service.NewGauge(
		GaugeOpts{
			Name: "words_operation_results",
			Help: "Handler and sending operations results",
		},
	)
)
