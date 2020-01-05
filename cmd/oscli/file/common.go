package file

import (
	"io"
	"os"

	"github.com/glynternet/oscli/internal/record"
	"github.com/glynternet/pkg/log"
	"github.com/pkg/errors"
)

func readFromFile(logger log.Logger, oscFile string) (record.Recording, error) {
	f, err := os.OpenFile(oscFile, os.O_RDONLY, 0400)
	if err != nil {
		return record.Recording{}, errors.Wrap(err, "opening file")
	}

	var recording record.Recording
	_, err = recording.ReadFrom(f)
	err = errors.Wrap(err, "reading recording from file")
	cErr := errors.Wrap(f.Close(), "closing file")
	if err == nil {
		return recording, cErr
	}
	if cErr != nil {
		_ = logger.Log(
			log.Message("Error closing file"),
			log.Error(cErr))
	}
	return recording, err
}

func writeToWriteCloser(r io.WriterTo, wc io.WriteCloser) (error, error) {
	_, wErr := r.WriteTo(wc)
	return errors.Wrap(wErr, "writing to WriteCloser"),
		errors.Wrap(wc.Close(), "closing WriteCloser")
}

func catchFirstLogOthers(logger log.Logger, errs ...error) error {
	var caught error
	for _, err := range errs {
		if err == nil {
			continue
		}
		if caught == nil {
			caught = err
			continue
		}
		_ = logger.Log(
			log.Message("Lower priority error occurred"),
			log.Error(err))
	}
	return caught
}

// fileCreatingWriteCloser creates a new file and provides a WriteCloser implementation that will write and log to it when called
func fileCreatingWriteCloser(logger log.Logger, path string) (io.WriteCloser, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrapf(err, "creating file at %s", path)
	}
	return &loggingFileWriteCloser{
		path:   path,
		file:   file,
		logger: logger,
	}, nil
}

type loggingFileWriteCloser struct {
	// path can be retrieved from create os.File reference but leaving it here for convenience
	path   string
	file   *os.File
	logger log.Logger
}

func (l *loggingFileWriteCloser) Write(p []byte) (n int, err error) {
	if err := l.logger.Log(log.Message("Writing to file"),
		log.KV{K: "path", V: l.path}); err != nil {
		return 0, errors.Wrap(err, "writing log message")
	}

	return l.file.Write(p)
}

func (l *loggingFileWriteCloser) Close() error {
	return errors.Wrap(l.file.Close(), "closing file")
}
