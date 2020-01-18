package duplicatefinder

import "path/filepath"

// GetFiles finds all .csv files in a specified directory
func GetFiles(path string) ([]string, error) {
	files, err := filepath.Glob(path + "/*.csv")
	if err != nil {
		return nil, err
	}

	return files, nil
}
