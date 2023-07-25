package misc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestMultiRun(t *testing.T) {
	err := MultiRun(func() error {
		time.Sleep(time.Second * 2)
		log.Println("func0")
		return nil
	}, func() error {
		log.Println("func1")
		return errors.New("func1 error")
	}, func() error {
		log.Println("func2")
		return nil
	})

	if err == nil {
		t.Fatal("should return err func1 error")
	}
}

func TestMultiRunWithCtx(t *testing.T) {
	err := MultiRunWithCtx(func(ctx context.Context) error {
		<-time.After(time.Second)
		select {
		case <-ctx.Done():
			if ctx.Err() != context.Canceled {
				t.Fatal("ctx err must be ctx canceled")
			}
			return nil
		default:
			return nil
		}
	}, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			t.Fatal("func1 ctx err should empty", ctx.Err())
			return nil
		default:
			return errors.New("func1 error")
		}
	})

	if err == nil {
		t.Fatal("err should return")
	}

	<-time.After(time.Second * 2)
}

func TestSafeGo(t *testing.T) {
	err := SafeGo(nil, nil)
	if err != nil {
		t.Fatal("err should be nil", err)
	}

	err = SafeGo(nil, func(ctx context.Context) error {
		return nil
	}, nil)
	if err != nil {
		t.Fatal("err should be nil", err)
	}

	err = SafeGo(nil, func(ctx context.Context) error {
		time.Sleep(time.Second)
		fmt.Println("A")
		return nil
	}, func(ctx context.Context) error {
		fmt.Println("B")
		return nil
	})
	if err != nil {
		t.Fatal("err should be nil", err)
	}

	err = SafeGo(nil, func(ctx context.Context) error {
		time.Sleep(time.Second * 60)
		return nil
	}, nil, func(ctx context.Context) error {
		return errors.New("xxx")
	})
	if err == nil {
		t.Fatal("err should not be nil")
	}
}

func TestHystrix(t *testing.T) {
	for i := 0; i < 1000; i++ {
		Hystrix("test1", func() error {
			return errors.New("test error")
		}, func(e error) error {
			return nil
		})
	}

	for i := 0; i < 1000; i++ {
		Hystrix("test2", func() error {
			return errors.New("test error")
		}, nil)
	}
}

func BenchmarkMultiRun(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MultiRun(func() error {
			return nil
		}, func() error {
			return nil
		})
	}
}

func BenchmarkSafeGo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SafeGo(nil, func(ctx context.Context) error {
			return nil
		}, func(ctx context.Context) error {
			return nil
		})
	}
}
