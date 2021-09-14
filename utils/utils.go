/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package utils

import (
	"fmt"
	"strings"
	"time"
)

// ===== [ Constants and Variables ] =====

const layoutISO = "2006.01.02"

var ()

// ========== [ IndexNames START ] =========

// ResolveIndexNames - 지정한 접두어와 시간 (시작/종료) 기준으로 Index 등에 활용할 수 있는 문자열들을 ','로 구분되는 문자열로 구성 (ex. Elasticsearch log filename, ...)
// conditions:
// - start 시간이 end 보다 이후인 경우는 감안하지 않는다. 이 경우는 이 함수 사용 이전에 검증되어 걸러져야 하는 것으로 판단한다.
func ResolveIndexNames(prefix string, start, end time.Time) string {
	// 종료 시간이 없는 경우는 현재 시각 설정
	if end.IsZero() {
		end = time.Now()
	}

	// 시작 시간이 없거나 종료 시간 기준으로 30일이 넘은 경우는 prefix만 반환
	if start.IsZero() || end.Sub(start).Hours() > 24*30 {
		return fmt.Sprintf("%s*", prefix)
	}

	var indices []string
	days := int(end.Sub(start).Hours() / 24)
	if start.Add(time.Duration(days)*24*time.Hour).UTC().Day() != end.UTC().Day() {
		days++
	}

	for i := 0; i <= days; i++ {
		suffix := end.Add(time.Duration(-i) * 24 * time.Hour).UTC().Format(layoutISO)
		indices = append(indices, fmt.Sprintf("%s-%s", prefix, suffix))
	}

	return strings.Join(indices, ",")
}

// ========== [ IndexNames END ] =========
