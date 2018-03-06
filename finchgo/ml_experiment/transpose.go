package main

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type readWriteSeekCloser interface {
	io.ReadWriteCloser
	io.Seeker
}

// FileBuffer saves the columns in temporary files.
type FileBuffer struct {
	names []string
	rwcs  []readWriteSeekCloser
	err   error

	size int
	sep  byte
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// Remove removes all temporary files
func (b *FileBuffer) Remove() {
	for _, n := range b.names {
		if fileExists(n) {
			if err := os.Remove(n); err != nil {
				log.Printf("error removing '%s': %v", n, err)
			}
		}
	}
}

// WriteTo writes the content from the temporary files
// into the result file
func (b *FileBuffer) WriteTo(w io.Writer) (int64, error) {
	size := b.size
	if size == 0 {
		size = 32 * 1024
	}
	var sum int64
	buf := make([]byte, size)
	for _, r := range b.rwcs {
		_, err := r.Seek(0, io.SeekStart)
		if err != nil {
			return 0, err
		}
		n, err := io.CopyBuffer(w, r, buf)
		if err != nil {
			return n, err
		}
		sum += n

		i, err := w.Write([]byte("\n"))
		if err != nil {
			return int64(i), err
		}
		sum += int64(i)
	}
	return sum, nil
}

// Store stores the content of the csv.Reader in
// temporary files
func (b *FileBuffer) Store(r *csv.Reader) error {
	for {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		err = b.append(line)
		if err != nil {
			return err
		}
	}
}

func (b *FileBuffer) append(line []string) error {
	if len(b.rwcs) == 0 {
		for _ = range line {
			f, err := ioutil.TempFile("", "transposer")
			if err != nil {
				return err
			}
			b.names = append(b.names, f.Name())
			b.rwcs = append(b.rwcs, f)
		}
	} else if len(line) != len(b.rwcs) {
		return errors.New("")
	}

	for i, s := range line {
		_, err := b.rwcs[i].Write(append([]byte(s), b.sep))
		if err != nil {
			return err
		}
	}
	return nil
}

func transposeCsv(csvFile io.Reader, w io.Writer) error {
	r := csv.NewReader(csvFile)
	buf := &FileBuffer{
		size: 32 * 1024,
		sep:  byte(','),
	}
	defer buf.Remove()

	err := buf.Store(r)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	return err
}
