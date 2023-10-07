package utils

import (
	"fmt"
	libutils "github.com/EscanBE/go-lib/utils"
	"io/fs"
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
