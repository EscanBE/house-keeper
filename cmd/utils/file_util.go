package utils

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"strconv"
)

func ValidatePasswordFileMode(mode fs.FileMode) error {
	str := fmt.Sprintf("%o", int(mode))
	symbolNum, err := strconv.ParseInt(str, 10, 64)
	libutils.PanicIfErr(err, fmt.Sprintf("failed to parse %s", str))
	if symbolNum%100 != 0 {
		return fmt.Errorf("not allowed to have permission for group/other")
	}
	if symbolNum < 400 {
		return fmt.Errorf("require read permission")
	}
	return nil
}

func IsFileAndExists(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.Wrap(err, fmt.Sprintf("problem while checking target file %s", file))
	}

	if fi.IsDir() {
		return false, nil
	}

	return true, nil
}
