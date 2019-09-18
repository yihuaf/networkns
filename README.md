# networkns
A small library to manage network namespace in Linux. Taken from `cri-o`,
`runc` and `vishvananda/netns`.

## Test

`make test`

## NOTE

The library can be safely used only with Go >= 1.10 due to [golang/go#20676](https://github.com/golang/go/issues/20676).

After locking a goroutine to its current OS thread with `runtime.LockOSThread()`
and changing its network namespace, any new subsequent goroutine won't be
scheduled on that thread while it's locked. Therefore, the new goroutine
will run in a different namespace leading to unexpected results.

See [here](https://www.weave.works/blog/linux-namespaces-golang-followup) for more details.
