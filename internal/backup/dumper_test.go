package backup_test

import (
	"sync/atomic"
	"test/internal/backup"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"go.uber.org/goleak"
)

var (
	testPeriod = time.Second
)

func TestDumper(t *testing.T) {
	var counter int32

	defer goleak.VerifyNone(t)

	dumpFn := func() error {
		atomic.AddInt32(&counter, 1)
		return nil
	}
	dumper := backup.NewDumper(dumpFn, &testPeriod)
	dumper.Start()

	for i := 1; i < 31; i++ {
		dumper.ScheduleUpdate()
		time.Sleep(time.Millisecond * 100)
	}
	dumper.ScheduleStop()
	time.Sleep(testPeriod)

	assert.Equal(t, counter, int32(3))
}
