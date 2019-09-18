// +build linux

package networkns

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

const nsRunDir = "/var/run/ns"

// NetworkNs ...
type NetworkNs struct {
	f *os.File
}

// New ....
func New() (*NetworkNs, error) {
	if err := unix.Unshare(unix.CLONE_NEWNET); err != nil {
		return nil, err
	}

	return Get()
}

// NewWithName ...
func NewWithName(name string) (*NetworkNs, error) {
	if err := unix.Unshare(unix.CLONE_NEWNET); err != nil {
		return nil, err
	}

	nsPath := filepath.Join(nsRunDir, name)

	if err := unix.Mount(GetCurrentThreadNsPath(), nsPath, "none", unix.MS_BIND,
		""); err != nil {
		return nil, err
	}

	return GetFromPath(nsPath)
}

// Close ...
func (ns *NetworkNs) Close() error {
	return ns.f.Close()
}

// GetCurrentThreadNsPath gets the network namespace path of the current thread.
func GetCurrentThreadNsPath() string {
	return (GetThreadNsPath(os.Getpid(), unix.Gettid()))
}

// GetThreadNsPath gets the path of the network namespace under /proc for a
// given pid and tid
func GetThreadNsPath(pid, tid int) string {
	return fmt.Sprintf("/proc/%d/task/%d/ns/net", pid, tid)
}

// GetFromPath gets a handle to a network namespace
// identified by the path
func GetFromPath(path string) (*NetworkNs, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return &NetworkNs{
		f: f,
	}, nil
}

// GetFromName gets a handle to a named network namespace such as one
// created by `ip netns add`.
func GetFromName(name string) (*NetworkNs, error) {
	return GetFromPath(fmt.Sprintf("/var/run/netns/%s", name))
}

// GetFromPid gets a handle to the network namespace of a given pid.
func GetFromPid(pid int) (*NetworkNs, error) {
	return GetFromPath(fmt.Sprintf("/proc/%d/ns/net", pid))
}

// GetFromThread gets a handle to the network namespace of a given pid and tid.
func GetFromThread(pid, tid int) (*NetworkNs, error) {
	return GetFromPath(GetThreadNsPath(pid, tid))
}

// Get gets a handle to the current threads network namespace.
func Get() (*NetworkNs, error) {
	return GetFromThread(os.Getpid(), syscall.Gettid())
}
