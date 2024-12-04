package loader

import (
	"bufio"
	"os"
	"strings"

	"github.com/hatlonely/hello-golang/pkg/kvstore"
	"github.com/hatlonely/hello-golang/pkg/refx"
	"github.com/pkg/errors"
)

func init() {
	refx.Register("loader", "FileLoader", NewFileLoader)
}

type FileLoaderOptions struct {
	Path string
}

type FileLoader struct {
	path string

	listener   kvstore.Listener
	streamChan chan kvstore.KVStream
}

func NewFileLoader(options *FileLoaderOptions) (*FileLoader, error) {
	loader := &FileLoader{
		path:       options.Path,
		streamChan: make(chan kvstore.KVStream, 2),
	}

	go func() {
		for stream := range loader.streamChan {
			loader.listener(stream)
		}
	}()

	return loader, nil
}

func (l *FileLoader) OnChange(listener kvstore.Listener) error {
	l.listener = listener

	file, err := os.Open(l.path)
	if err != nil {
		return errors.Wrap(err, "os.Open failed")
	}
	// defer file.Close()

	l.streamChan <- &FileLoaderKVStream{
		scanner:   bufio.NewScanner(file),
		separator: "\t",
	}

	return nil
}

func (l *FileLoader) Close() error {
	close(l.streamChan)
	return nil
}

type FileLoaderKVStream struct {
	scanner *bufio.Scanner

	separator string
}

func (s *FileLoaderKVStream) HasNext() bool {
	return s.scanner.Scan()
}

func (s *FileLoaderKVStream) Next() (any, any, error) {
	line := s.scanner.Text()
	parts := strings.Split(line, s.separator)

	if len(parts) != 2 {
		return nil, nil, errors.Errorf("invalid line [%v]", line)
	}

	return parts[0], parts[1], nil
}
