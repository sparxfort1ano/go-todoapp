package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

func (r *WebRepository) GetFile(filePath string) ([]byte, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf(
				"file: %s: %w",
				filePath,
				errs.ErrNotFound,
			)
		}

		return nil, fmt.Errorf(
			"get file: %s: %w",
			filePath,
			err,
		)
	}

	return file, nil
}
