/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package strings

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
)

// ===== [ Constants and Variables ] =====

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// Diff - 지정한 Source Slice에서 Target Slice를 포함하지 않는 Slice 반환
func Diff(source, target []string) (result []string) {
	excludeMap := make(map[string]bool)

	// 제외할 맵 생성
	for _, s := range target {
		excludeMap[s] = true
	}

	// 반환 대상 맵 정보 구성
	for _, s := range source {
		if !excludeMap[s] {
			result = append(result, s)
		}
	}

	return result
}

// Unique - 지정한 문자열 배열에서 중복된 값을 제거한 유일 값들만 반환
// conditions
// - 동일한 키인 경우는 가장 먼저 찾아지는 것을 기준으로 하고 나머지는 제거
func Unique(source []string) (result []string) {
	uniqueMap := make(map[string]bool)

	// 유일한 경우만 맵 구성
	for _, s := range source {
		uniqueMap[s] = true
	}

	// 반환 대상 구성
	for s := range uniqueMap {
		result = append(result, s)
	}

	return result
}

// CamelCaseToUnderscore - CamelCase로 문장을 Underscore 문장으로 변환
func CamelCaseToUnderscore(str string) string {
	return govalidator.CamelCaseToUnderscore(str)
}

// UndercoreToCamelCase - Underscore 문장을 CamelCase 문장으로 변환
func UndercoreToCamelCase(str string) string {
	return govalidator.UnderscoreToCamelCase(str)
}

// FindString - 지정한 문자열 배열에서 지정한 문자열 검색 후 인덱스 반환
// conditions
// - 존재하지 않으면 `-1` 반환
func FindString(source []string, str string) int {
	for idx, s := range source {
		if str == s {
			return idx
		}
	}

	return -1
}

// StringIn - 지정한 문자열이 지정한 문자열 배열에 존재하는지 여부 반환
// conditions:
// - 한글처럼 2byte 문자는 utf8 DecodeRuneInString / utf8.EncodeRune 함수 사용해서 처리
func StringIn(str string, target []string) bool {
	return FindString(target, str) > -1
}

// Reverse - 지정한 문자열을 역순으로 재 구성하여 반환
func Reverse(str string) string {
	size := len(str)
	buf := make([]byte, size)

	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(str[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

// Split - 지정한 문자열을 지정한 구분자를 기준으로 분리해서 배열로 반환
// conditions:
// - 빈 문자열 `""` 인 경우는 nil 반환
func Split(str string, sep string) []string {
	if str == "" {
		return nil
	}
	return strings.Split(str, sep)
}

// StripAnsi - 지정한 문자열에서 Ansi 코드를 제거한 문자열 반환
func StripAnsi(str string) string {
	return re.ReplaceAllString(str, "")
}

// ShortenString - 지정한 문자열을 지정한 길이로 잘라서 반환
// conditions:
// - 문자열 길이가 지정한 길이보다 작거나 같은 경우는 그대로 반환, 그 외는 처음부터 지정한 길이까지 반환
func ShortenString(str string, n int) string {
	if len(str) <= n {
		return str
	}
	return str[:n]
}
