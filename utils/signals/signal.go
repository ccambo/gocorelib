//go:build !windows
// +build !windows

/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package signals

import (
	"os"
	"os/signal"
)

// ===== [ Constants and Variables ] =====
const ()

var onlyOneSignalHandler = make(chan struct{}) // 단일 시그널 처리기 관리용

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// SetupSignalHandler - SIGTERM and SIGINT 등록
// conditions:
// - 등록된 Signal에 따라 종료될 채널 반환
// - 두번째 시그널이 전달되면 프로그램을 -1 반환 코드로 종료된다.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // 두번 호출되면 Panic 발생

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
