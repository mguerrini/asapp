package helpers

import (
	"fmt"
	"github.com/challenge/pkg/modules/logger"
	"runtime/debug"
)

func Recover() error {
	if err := recover(); err != nil {
		logger.Errorf("[Custom Recovery] panic recovered: %s %s", fmt.Errorf("%s", err), debug.Stack(), err)
		return fmt.Errorf("%s", err)
	}

	return nil
}

