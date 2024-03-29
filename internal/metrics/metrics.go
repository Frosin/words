package metrics

const metricsPrefix = ""

var service = NewMetrics().WithNamespace(metricsPrefix)

var (
	WordsOperationResults = service.NewGaugeVec(
		GaugeOpts{
			Name: "words_operation_results",
			Help: "Handler and sending operations results",
		},
		[]string{"data", "error"},
	)

	WordsPhraseAdded = service.NewCounter(
		CounterOpts{
			Name: "word_phrase_add",
			Help: "user add some phrase",
		},
	)

	WordsPhraseEpoch1 = service.NewCounter(
		CounterOpts{
			Name: "word_phrase_epoch1",
			Help: "user add some text with phrase 1 epoch",
		},
	)

	WordsPhraseEpoch2 = service.NewCounter(
		CounterOpts{
			Name: "word_phrase_epoch2",
			Help: "user add some text with phrase 2 epoch",
		},
	)

	WordsPhraseEpoch3 = service.NewCounter(
		CounterOpts{
			Name: "word_phrase_epoch3",
			Help: "user add some text with phrase 3 epoch",
		},
	)

	WordsRequestDuration = service.NewHistogramVec(
		HistogramOpts{
			Name: "word_reqests_duration",
			Help: "repository requests elapsed time",
		},
		[]string{"request_name"},
	)
)
