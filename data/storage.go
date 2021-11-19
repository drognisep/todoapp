package data

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	dataFile = ".todo_file"
)

func LoadTaskData() (*Model, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	home = filepath.ToSlash(home)
	data, err := os.ReadFile(path.Join(home, dataFile))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Model{}, nil
		}
		return nil, errors.Wrap(err, "Failed to read from data file")
	}
	if len(data) == 0 {
		return &Model{}, nil
	}
	var model *Model
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, errors.Wrap(err, "Invalid data file format")
	}
	return model, nil
}

func SaveTaskData(model *Model) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	home = filepath.ToSlash(home)
	jsonData, err := json.Marshal(&model)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize JSON data")
	}
	if err := os.WriteFile(path.Join(home, dataFile), jsonData, 0777); err != nil {
		return errors.Wrap(err, "Failed to write task data")
	}
	return nil
}
