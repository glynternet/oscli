package file

import (
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

func writeToFile(logger log.Logger, r record.Recording, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "creating file at %s", path)
	}
	if err := logger.Log(log.Message("Writing to file"),
		log.KV{K: "path", V: path}); err != nil {
		return errors.Wrap(err, "writing log message")
	}

	_, err = r.WriteTo(file)
	err = errors.Wrap(err, "writing recording to writer")
	cErr := errors.Wrap(file.Close(), "closing file")
	if err == nil {
		return cErr
	}
	if cErr != nil {
		_ = logger.Log(
			log.Message("Error closing file"),
			log.Error(cErr))
	}
	return err
}
