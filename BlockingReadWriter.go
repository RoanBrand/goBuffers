package goBuffers

import "bytes"

type BlockingReadWriter struct {
	buf   *bytes.Buffer
	read  chan []byte
	write chan []byte

	readError  chan error
	writeError chan error
	nRead      chan int
	nWrite     chan int
}

func NewBlockingReadWriter() *BlockingReadWriter {
	b := &BlockingReadWriter{
		buf:        bytes.NewBuffer(nil),
		read:       make(chan []byte),
		write:      make(chan []byte),
		readError:  make(chan error),
		writeError: make(chan error),
		nRead:      make(chan int),
		nWrite:     make(chan int),
	}
	go func(b *BlockingReadWriter) {
		commit := func(p []byte) {
			wN, wE := b.buf.Write(p)
			b.nWrite <- wN
			b.writeError <- wE
		}
		for {
			select {
			case rx := <-b.read:
				if b.buf.Len() == 0 {
					commit(<-b.write)
				}
				rN, rE := b.buf.Read(rx)
				b.nRead <- rN
				b.readError <- rE
			case tx := <-b.write:
				commit(tx)
			}
		}
	}(b)
	return b
}

func (b *BlockingReadWriter) Read(p []byte) (n int, err error) {
	if p == nil {
		return 0, nil
	}
	b.read <- p
	n = <-b.nRead
	err = <-b.readError
	return
}

func (b *BlockingReadWriter) Write(p []byte) (n int, err error) {
	if p == nil {
		return 0, nil
	}
	b.write <- p
	n = <-b.nWrite
	err = <-b.writeError
	return
}
