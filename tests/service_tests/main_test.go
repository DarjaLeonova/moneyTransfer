package service_tests

import (
	"moneyTransfer/pkg/logger"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	logger.Init()

	code := m.Run()
	os.Exit(code)
}
