package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_isWorkerSleepTime(t *testing.T) {
	t.Parallel()
	t.Run("before, sleep time", func(t *testing.T) {
		startTime, err := time.Parse("15:04", "20:20")
		require.NoError(t, err)
		now, err := time.Parse(time.RFC3339, "2006-01-02T20:19:05+03:00")
		require.NoError(t, err)
		duration := time.Hour * 3

		res := isWorkerSleepTime(startTime, now, duration)
		require.True(t, res)
	})

	t.Run("in, NOT sleep time", func(t *testing.T) {
		startTime, err := time.Parse("15:04", "20:20")
		require.NoError(t, err)
		now, err := time.Parse(time.RFC3339, "2006-01-02T20:21:05+03:00")
		require.NoError(t, err)
		duration := time.Hour * 3

		res := isWorkerSleepTime(startTime, now, duration)
		require.False(t, res)
	})

	t.Run("after, sleep time", func(t *testing.T) {
		startTime, err := time.Parse("15:04", "20:20")
		require.NoError(t, err)
		now, err := time.Parse(time.RFC3339, "2006-01-02T23:21:05+03:00")
		require.NoError(t, err)
		duration := time.Hour * 3

		res := isWorkerSleepTime(startTime, now, duration)
		require.True(t, res)
	})

}
