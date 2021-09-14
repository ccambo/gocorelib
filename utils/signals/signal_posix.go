//go:build !windows

/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package signals

import (
	"os"
	"syscall"
)

// ===== [ Constants and Variables ] =====
const ()

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM} // POSIX 환경인 경우 종료 시그널

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====
