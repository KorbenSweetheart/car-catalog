# Ideas

- [ ] can we try to use Composit design pattern to form a tree structure? Potentially it can be used to filter data from categories
- [ ] try zero-copy techniques

1. `io.Copy()` with file → socket
If both `src` and `dst` implement certain interfaces (like `*os.File` and `*net.TCPConn`), Go’s runtime can use `sendfile()` behind the scenes:

```
f, _ := os.Open("video.mp4")
conn, _ := net.Dial("tcp", "example.com:9000")

io.Copy(conn, f) // uses sendfile() → kernel-to-kernel copy
```

> Data goes directly from file descriptor to socket descriptor — your Go code never touches the bytes.
> This is a true zero-copy path.


✅ “The event system runs as asynchronous (usage of goroutines and channels)”

This is where your approach actually shines.

You can cleanly implement:
- Background goroutine:
    - Periodically refresh Cars API data
- Channels:
    - Swap cached datasets safely
    - Signal refresh completion
    - Handle API failures gracefully

This strongly satisfies the requirement.