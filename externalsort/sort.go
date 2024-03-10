//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"errors"
	"io"
	"os"
	"slices"
	"strings"
)

type ReaderImpl struct {
	ioReader    io.Reader
	bufioReader *bufio.Reader
}

func (ri *ReaderImpl) ReadLine() (string, error) {
	var sb strings.Builder

	for {
		b, err := ri.bufioReader.ReadByte()

		if err != nil {
			return sb.String(), err
		}

		if b == '\n' {
			return sb.String(), nil
		}

		sb.WriteByte(b)
	}
}

func NewReader(r io.Reader) LineReader {
	return &ReaderImpl{
		ioReader:    r,
		bufioReader: bufio.NewReader(r),
	}
}

type WriterImpl struct {
	ioWriter io.Writer
}

func (wi *WriterImpl) Write(line string) error {
	_, err := wi.ioWriter.Write([]byte(line))

	if err == nil {
		_, err = wi.ioWriter.Write([]byte{'\n'})
	}

	return err
}

func NewWriter(w io.Writer) LineWriter {
	return &WriterImpl{
		ioWriter: w,
	}
}

type HeapItem struct {
	lr  *LineReader
	top string
}

type Heap []*HeapItem

func (h Heap) Len() int { return len(h) }

func (h Heap) Less(i, j int) bool {
	return strings.Compare(h[i].top, h[j].top) == -1
}

func (h Heap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Heap) Push(x any) {
	*h = append(*h, x.(*HeapItem))
}

func (h *Heap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[0 : n-1]
	return item
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := make(Heap, 0, len(readers))

	for _, reader := range readers {
		str, err := reader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) || errors.Is(err, io.EOF) && len(str) == 0 {
			continue
		}

		heap.Push(&h, &HeapItem{
			lr:  &reader,
			top: str,
		})
	}

	for h.Len() > 0 {
		item := heap.Pop(&h).(*HeapItem)
		err := w.Write(item.top)

		if err != nil {
			return err
		}

		str, err := (*item.lr).ReadLine()
		if err == nil || errors.Is(err, io.EOF) && len(str) > 0 {
			heap.Push(&h, &HeapItem{
				lr:  item.lr,
				top: str,
			})
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	var err error

	for _, filename := range in {
		var f *os.File
		var lines []string

		func() {
			f, err = os.OpenFile(filename, os.O_RDONLY, 0600)

			if err != nil {
				return
			}

			defer func(f *os.File) {
				err = f.Close()
			}(f)

			lr := NewReader(f)

			for {
				line, lrErr := lr.ReadLine()

				if lrErr == nil || errors.Is(lrErr, io.EOF) && len(line) > 0 {
					lines = append(lines, line)
				}

				if lrErr != nil {
					break
				}
			}
		}()

		func() {
			if err != nil {
				return
			}

			if len(lines) > 1 {
				slices.SortFunc(lines, strings.Compare)
			}

			f, err = os.OpenFile(filename, os.O_WRONLY, 0600)

			if err != nil {
				return
			}

			defer func(f *os.File) {
				err = f.Close()
			}(f)

			lw := NewWriter(f)

			for _, line := range lines {
				err = lw.Write(line)

				if err != nil {
					return
				}
			}
		}()

		if err != nil {
			return err
		}
	}

	readers := make([]LineReader, 0, len(in))
	lw := NewWriter(w)

	for _, filename := range in {
		var f *os.File
		f, err = os.OpenFile(filename, os.O_RDONLY, 0600)

		if err != nil {
			return err
		}

		//goland:noinspection GoDeferInLoop
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		lr := NewReader(f)
		readers = append(readers, lr)
	}

	return Merge(lw, readers...)
}
