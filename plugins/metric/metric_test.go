package metric

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func Test(t *testing.T) {
	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			CounterVec.Count("test_count", []string{"k1", "v1", "k2", "v2"})
		}()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			GaugeVec.Set("test_gauge", []string{"k1", "v1", "k2", "v2"}, 2)
		}()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startAt := time.Now()
			<-time.After(time.Millisecond * time.Duration(rand.Intn(100)))
			HistogramVec.Timing("test_his_seconds", []string{"k1", "v1", "k2", "v2"}, startAt)
		}()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			startAt := time.Now()
			<-time.After(time.Millisecond * time.Duration(rand.Intn(100)))
			SummaryVec.Timing("test_sum_seconds", []string{"k1", "v1", "k2", "v2"}, startAt)
		}()
	}

	wg.Wait()

	CounterVec.Count("test_count", []string{"k1", "v1", "k2", "v2"}, 5)
	GaugeVec.Set("test_gauge", []string{"k1", "v1", "k2", "v2"}, 10)
	HistogramVec.Timing("test_his_seconds", []string{"k1", "v1", "k2", "v2"}, time.Now())
	SummaryVec.Timing("test_sum_seconds", []string{"k1", "v1", "k2", "v2"}, time.Now())

	mc := make(chan prometheus.Metric, 1)
	mt := &dto.Metric{}

	CounterVec.m.bag["test_count"].(*prometheus.CounterVec).Collect(mc)
	m := <-mc
	m.Write(mt)
	fmt.Printf("%#v\n", mt.String())
	if !strings.Contains(mt.String(), "105") {
		t.Fatal("counter should be 105")
	}

	GaugeVec.m.bag["test_gauge"].(*prometheus.GaugeVec).Collect(mc)
	m = <-mc
	mt.Reset()
	m.Write(mt)
	fmt.Printf("%#v\n", mt.String())
	if !strings.Contains(mt.String(), "10") {
		t.Fatal("counter should be 10")
	}

	HistogramVec.m.bag["test_his_seconds"].(*prometheus.HistogramVec).Collect(mc)
	m = <-mc
	mt.Reset()
	m.Write(mt)
	fmt.Printf("%#v\n", mt.String())
	if !strings.Contains(mt.String(), "sample_count:101") {
		t.Fatal("his sample_count should be 101")
	}

	SummaryVec.m.bag["test_sum_seconds"].(*prometheus.SummaryVec).Collect(mc)
	m = <-mc
	mt.Reset()
	m.Write(mt)
	fmt.Printf("%#v\n", mt.String())
	if !strings.Contains(mt.String(), "sample_count:101") {
		t.Fatal("sum sample_count should be 101")
	}

	for i := 0; i < 100; i++ {
		fmt.Println(genLabels(map[string]string{"a": "b", "x": "y"}))
	}

}
