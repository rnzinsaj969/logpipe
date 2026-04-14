package checkpoint_test

import (
	"sync"
	"testing"

	"github.com/yourorg/logpipe/internal/checkpoint"
)

func TestConcurrentSetAndFlush(t *testing.T) {
	path := tempPath(t)
	s, err := checkpoint.New(path)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "svc"
			s.Set(key, checkpoint.State{Offset: int64(i), Inode: uint64(i)})
			_ = s.Flush()
		}(i)
	}
	wg.Wait()

	// After all goroutines finish, a reload must succeed without error.
	_, err = checkpoint.New(path)
	if err != nil {
		t.Fatalf("reload after concurrent flush failed: %v", err)
	}
}

func TestMultipleSourcesRoundTrip(t *testing.T) {
	path := tempPath(t)
	s, _ := checkpoint.New(path)

	sources := []string{"web", "worker", "cron", "db"}
	for idx, src := range sources {
		s.Set(src, checkpoint.State{Offset: int64(idx * 10), Inode: uint64(idx + 1)})
	}
	if err := s.Flush(); err != nil {
		t.Fatal(err)
	}

	s2, _ := checkpoint.New(path)
	for idx, src := range sources {
		st := s2.Get(src)
		if st.Offset != int64(idx*10) || st.Inode != uint64(idx+1) {
			t.Errorf("source %s: got %+v", src, st)
		}
	}
}
