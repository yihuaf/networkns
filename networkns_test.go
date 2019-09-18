package networkns

import (
	"runtime"
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
