package metrics

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricsServer(t *testing.T) {
	RunMetrics()

	time.Sleep(time.Millisecond * 100)

	WordsOperationResults.WithLabelValues("", "").Add(777)

	client := http.Client{
		Timeout: 1 * time.Second,
	}
	path := fmt.Sprintf("http://:%s%s", defaultMetricsPort, defaultMetricsPath)

	res, err := client.Get(path)
	assert.NoError(t, err)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	response := string(body)

	assert.Contains(t, response, "words_operation_results")
}
