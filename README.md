# goBuffers
Useful Buffers for Golang

## Blocking ReadWriter
A buffer of bytes that implements the `Reader` and `Writer` interfaces that will block if `Read` is called when the buffer is empty, until another Goroutine calls `Write` with some data.  

I use this to simulate the behaviour of `net.Conn`'s `Read` & `Write` in tests to block until IO is ready. (Fake net.Conn)

## Dynamic Buffered Channel
Buffered channels in Go are fixed size when created. This type is a drop-in replacement for a buffered channel that can grow in size at runtime as needed.  

This means that a Goroutine that sends to this type will never block.
