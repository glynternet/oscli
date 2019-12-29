package file

import (
	"log"
	"os"

	"github.com/glynternet/oscli/internal/record"
	"github.com/pkg/errors"
)

func readFromFile(logger *log.Logger, oscFile string) (record.Recording, error) {
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
		logger.Println(cErr)
	}
	return recording, err
}

func writeToFile(logger *log.Logger, r record.Recording, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "creating file at %s", path)
	}
	logger.Printf("Writing to file at %s", path)
	_, err = r.WriteTo(file)
	err = errors.Wrap(err, "writing recording to writer")
	cErr := errors.Wrap(file.Close(), "closing file")
	if err == nil {
		return cErr
	}
	if cErr != nil {
		logger.Println(cErr)
	}
	return err
}
