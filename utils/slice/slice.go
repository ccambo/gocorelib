/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package slice

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Types ] =====
type ()

// ===== [ Implementations ] =====
// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// RemoveString - 지정한 문자열 배열의 각 문자열들을 지정한 제거 함수에 전달하고 결과가 true 인 경우에 제거한 문자열 배열 반환
func RemoveString(slice []string, remove func(item string) bool) []string {
	for i := 0; i < len(slice); i++ {
		if remove(slice[i]) {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

// HasString - 지정한 문자열 배열에 지정한 문자열이 존재하는지를 반환
func HasString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
