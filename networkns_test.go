package networkns

import (
	"runtime"
	"sync"
	"testing"
)

func IsOpen(ns *NetworkNs) bool {
	return int(ns.f.Fd()) != -1
}

func TestGetNewSetDelete(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	origns, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	newns, err := New()
	if err != nil {
		t.Fatal(err)
	}

	if IsSame(newns, origns) {
		t.Fatal("New ns failed")
	}

	if err := Set(origns); err != nil {
		t.Fatal(err)
	}

	Close(newns)

	if IsOpen(newns) {
		t.Fatal("newns still open after close", newns)
	}

	ns, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	if !IsSame(ns, origns) {
		t.Fatal("Reset ns failed", origns, newns, ns)
	}
}

func TestThreaded(t *testing.T) {
	ncpu := runtime.GOMAXPROCS(-1)
	if ncpu < 2 {
		t.Skip("-cpu=2 or larger required")
	}

	// Lock this thread simply to ensure other threads get used.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	wg := &sync.WaitGroup{}
	for i := 0; i < ncpu; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			TestGetNewSetDelete(t)
		}()
	}
	wg.Wait()
}
