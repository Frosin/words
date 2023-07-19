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

	WordsGetPhrasesRequest = service.NewGauge(
		GaugeOpts{
			Name: "word_get_phrase_req",
			Help: "get phrases request elapsed time",
		},
	)
)
