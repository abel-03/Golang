//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

type lineReader struct {
	ioReader    io.Reader
	bytesReader *bufio.Reader
}

func (lr *lineReader) ReadLine() (string, error) {
	var sb strings.Builder

	for {
		b, err := lr.bytesReader.ReadByte()
		if err != nil {
			return sb.String(), err
		}

		if b == '\n' {
			return sb.String(), nil
		}

		sb.WriteString(string(b))
	}

}

func NewReader(r io.Reader) LineReader {
	return &lineReader{
		ioReader:    r,
		bytesReader: bufio.NewReader(r),
	}
}

type lineWriter struct {
	ioWriter io.Writer
}

func (lw *lineWriter) Write(l string) error {
	_, err := lw.ioWriter.Write([]byte(l + "\n"))
	return err
}

func NewWriter(w io.Writer) LineWriter {
	return &lineWriter{
		ioWriter: w,
	}
}

type heapItem struct {
	lr  LineReader
	top string
}

type heapSlice []heapItem

func (h heapSlice) Len() int           { return len(h) }
func (h heapSlice) Less(i, j int) bool { return strings.Compare(h[i].top, h[j].top) < 0 }
func (h heapSlice) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *heapSlice) Push(x interface{}) {
	*h = append(*h, x.(heapItem))
}

func (h *heapSlice) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &heapSlice{}
	for _, r := range readers {
		str, err := r.ReadLine()
		if err != nil && (err != io.EOF) || (err == io.EOF && str == "") {
			continue
		}
		heap.Push(h, heapItem{
			top: str,
			lr:  r,
		})
	}
	heap.Init(h)

	for h.Len() > 0 {
		minLineReaderHeapItem := heap.Pop(h).(heapItem)

		err := w.Write(minLineReaderHeapItem.top)

		if err != nil {
			return err
		}

		str, err := minLineReaderHeapItem.lr.ReadLine()
		if err == nil || (str != "" && err == io.EOF) {
			heap.Push(h, heapItem{
				top: str,
				lr:  minLineReaderHeapItem.lr,
			})
			heap.Fix(h, h.Len()-1)
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	var readers []LineReader

	for _, filename := range in {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		lr := NewReader(f)

		var lines []string
		for {
			str, rlErr := lr.ReadLine()

			if rlErr == io.EOF && str != "" {
				lines = append(lines, str)
			}

			if rlErr != nil {
				break
			}

			lines = append(lines, str)
		}

		sort.Strings(lines)

		f.Close()

		f, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return err
		}

		lw := NewWriter(f)

		for _, str := range lines {
			err := lw.Write(str)
			if err != nil {
				return err
			}
		}
		f.Close()
	}

	lw := NewWriter(w)

	for _, filename := range in {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		lr := NewReader(f)
		readers = append(readers, lr)
	}

	err := Merge(lw, readers...)

	if err != nil {
		return err
	}

	return nil
}
