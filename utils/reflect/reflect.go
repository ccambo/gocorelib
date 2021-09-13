/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package reflect

import "reflect"

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// In - 지정한 값이 지정한 형식에 존재하는지 검증
// condition:
// - Slice, Array, Map
func In(value interface{}, container interface{}) bool {
	containerVal := reflect.ValueOf(container)
	switch reflect.TypeOf(container).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < containerVal.Len(); i++ {
			if containerVal.Index(i).Interface() == value {
				return true
			}
		}
	case reflect.Map:
		if containerVal.MapIndex(reflect.ValueOf(value)).IsValid() {
			return true
		}
	default:
		return false
	}

	return false
}

// Override - 지정한 source를 target에 설정
// condition:
// - source / target 중에 하나라도 nil인 경우는 처리하지 않는다.
// - source / target 중 Ptr이 아니거나 둘의 형식이 다른 경우는 처리하지 않는다.
func Override(target interface{}, source interface{}) {
	// Nil 검사
	if reflect.ValueOf(target).IsNil() || reflect.ValueOf(source).IsNil() {
		return
	}

	// Prt 형식 및 동일 형식 검사
	if reflect.ValueOf(target).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(source).Type().Kind() != reflect.Ptr ||
		reflect.ValueOf(target).Kind() != reflect.ValueOf(source).Kind() {
		return
	}

	targetVal := reflect.ValueOf(target).Elem()
	sourceVal := reflect.ValueOf(source).Elem()

	for i := 0; i < targetVal.NumField(); i++ {
		val := sourceVal.Field(i).Interface()
		if !reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
			targetVal.Field(i).Set(reflect.ValueOf(val))
		}
	}
}
