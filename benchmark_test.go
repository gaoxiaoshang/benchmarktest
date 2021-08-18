package benchmarktest
import (
	"sync"
	"testing"
	"time"
	"github.com/panjf2000/ants/v2"
)

const (
	RunTimes           = 100000000
	BenchParam         = 10
	BenchAntsSize      = 2000000
	DefaultExpiredTime = 10 * time.Second

)

func demoFunc() {
	time.Sleep(time.Duration(BenchParam) * time.Millisecond)
}

func BenchmarkGoroutines(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			go func() {
				demoFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkChannel(b *testing.B) {
	var sema = make(chan struct{}, BenchAntsSize)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			sema <- struct{}{}
			go func() {
				demoFunc()
				<-sema
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkPool(b *testing.B) {
	var wg sync.WaitGroup
	p, _ := ants.NewPool(BenchAntsSize, ants.WithExpiryDuration(DefaultExpiredTime))
	defer p.Release()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes)
		for j := 0; j < RunTimes; j++ {
			_ = p.Submit(func() {
				demoFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
}