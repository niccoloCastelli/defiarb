package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

func WrapErrors(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}
	combinedErr := make([]string, len(errs))
	for i, err := range errs {
		combinedErr[i] = fmt.Sprintf("%d. %s", i, err.Error())
	}
	return errors.New(strings.Join(combinedErr, "\n"))
}
