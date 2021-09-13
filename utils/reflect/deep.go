/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package reflect

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

// ===== [ Constants and Variables ] =====
const ()

var (
	FloatPrecision          = 10    // 비교할 떄 반올림할 소수 자리 수
	MaxDiff                 = 10    // 반환할 최대 차이의 갯수
	MaxDepth                = 0     // struct 형식의 최대 recursive depth 수로 0을 지정하면 무제한
	LogErrors               = false // 오류 발생을 표준 오류 출력 (STDERR)으로 처리할지 여부
	CompareUnexportedFields = false // 예상치 못한 struct 필드들 (ex. T{s int}) 이 발생한 경우 비교 여부
)

var (
	ErrMaxRecursion = errors.New("recursed to MaxDepth")                 // 최대 Recursive Depth 보다 비교 요청 Depth가 큰 경우
	ErrTypeMismatch = errors.New("variables are different reflect.Type") // 값들이 서로 다른 형식인 경우
	ErrNotHandled   = errors.New("cannot compare the reflect.Kind")      // primitive 유형을 처리하지 못하는 경우
)

var errorType = reflect.TypeOf((*error)(nil)).Elem() // Error 형식

// ========== [ Compare START ] =========

// cmp - 비교를 위한 정보 관리용
type cmp struct {
	diff        []string // 차이점들
	buff        []string // recursive가 발생하는 변수 처리 버퍼
	floatFormat string   // 실수 형의 포맷 정보
}

// equals - 지정한 a, b를 지정한 Depth만큼 비교
// codition
// - depth 가 0인 경우는 무제한
// - 비교된 결과는 cmd struct diff 로 설정
// - struct field 에 `deep:"-"` 설정된 경우는 비교 생략
func (c *cmp) equals(a, b reflect.Value, level int) {
	// check depth
	if MaxDepth > 0 && level > MaxDepth {
		logError(ErrMaxRecursion)
		return
	}

	// check value is nil
	if !a.IsValid() || !b.IsValid() {
		if a.IsValid() && !b.IsValid() {
			c.saveDiff(a.Type(), "<nil pointer>")
		} else if !a.IsValid() && b.IsValid() {
			c.saveDiff("<nil pointer>", b.Type())
		}

		return
	}

	// check different type
	aType := a.Type()
	bType := b.Type()
	if aType != bType {
		c.saveDiff(aType, bType)
		logError(ErrTypeMismatch)
		return
	}

	// check primitive type https://golang.org/pkg/reflect/#Kind
	aKind := a.Kind()
	bKind := b.Kind()

	// check underlying elements? ptr or interface
	aElem := aKind == reflect.Ptr || aKind == reflect.Interface
	bElem := bKind == reflect.Ptr || bKind == reflect.Interface

	// 오류 형식의 구현체인 경우는 메시지 비교
	if aType.Implements(errorType) && bType.Implements(errorType) {
		if (!aElem || !a.IsNil()) && (!bElem || !b.IsNil()) {
			aString := a.MethodByName("Error").Call(nil)[0].String()
			bString := b.MethodByName("Error").Call(nil)[0].String()
			if aString != bString {
				c.saveDiff(aString, bString)
				return
			}
		}
	}

	// pointer, interface 인 경우 참조해제 방식으로 비교
	if aElem || bElem {
		if aElem {
			a = a.Elem()
		}
		if bElem {
			b = b.Elem()
		}
		c.equals(a, b, level+1)
		return
	}

	switch aKind {
	// 구조체 재귀처리
	case reflect.Struct:
		// Equal 함수 지원 여부
		if eqFunc := a.MethodByName("Equal"); eqFunc.IsValid() && eqFunc.CanInterface() {
			funcType := eqFunc.Type()
			// 비교를 위한 변수 지정 여부
			if funcType.NumIn() == 1 && funcType.In(0) == bType {
				retVals := eqFunc.Call([]reflect.Value{b})
				if !retVals[0].Bool() {
					c.saveDiff(a, b)
				}
				return
			}
		}

		for i := 0; i < a.NumField(); i++ {
			if aType.Field(i).PkgPath != "" && !CompareUnexportedFields {
				continue // 예상치 못한 필드는 생략. ex. s in t struct {s string}
			}

			if aType.Field(i).Tag.Get("deep") == "-" {
				continue // ignore
			}

			c.push(aType.Field(i).Name) // 대상 필드명 관리 추가

			// 필드 값 추출
			af := a.Field(i)
			bf := b.Field(i)

			// 재귀적으로 비교
			c.equals(af, bf, level+1)

			c.pop() // 관리중인 대사 필드 제거

			// 최대 비교 수 초과하면 종료
			if len(c.diff) >= MaxDiff {
				break
			}
		}
		// 맵 처리
	case reflect.Map:
		// nil check
		if a.IsNil() || b.IsNil() {
			if a.IsNil() && !b.IsNil() {
				c.saveDiff("<nil map>", b)
			} else if !a.IsNil() && b.IsNil() {
				c.saveDiff(a, "<nil map>")
			}
			return
		}

		// 동일한 주소를 가지는 경우
		if a.Pointer() == b.Pointer() {
			return
		}

		// a map 기준
		for _, key := range a.MapKeys() {
			// 필드 저장
			c.push(fmt.Sprintf("map[%s]", key))

			aVal := a.MapIndex(key)
			bVal := b.MapIndex(key)

			if bVal.IsValid() {
				c.equals(aVal, bVal, level+1)
			} else {
				c.saveDiff(aVal, "<does not have key")
			}

			c.pop()

			if len(c.diff) >= MaxDiff {
				return
			}
		}

		// b map 기준
		for _, key := range b.MapKeys() {
			if aVal := a.MapIndex(key); aVal.IsValid() {
				continue
			}

			c.push(fmt.Sprintf("map[%s]", key))
			c.saveDiff("<does not have key", b.MapIndex(key))
			c.pop()

			if len(c.diff) >= MaxDiff {
				return
			}
		}
		// Array 검증
	case reflect.Array:
		n := a.Len()
		for i := 0; i < n; i++ {
			c.push(fmt.Sprintf("array[%d]", i))
			c.equals(a.Index(i), b.Index(i), level+1)
			c.pop()

			if len(c.diff) >= MaxDiff {
				break
			}
		}
		// Slice
	case reflect.Slice:
		if a.IsNil() || b.IsNil() {
			if a.IsNil() && !b.IsNil() {
				c.saveDiff("<nil slice>", b)
			} else if !a.IsNil() && b.IsNil() {
				c.saveDiff(a, "<nil slice>")
			}
			return
		}

		aLen := a.Len()
		bLen := b.Len()

		if a.Pointer() == b.Pointer() && aLen == bLen {
			return
		}

		n := aLen
		if bLen > aLen {
			n = bLen
		}

		for i := 0; i < n; i++ {
			c.push(fmt.Sprintf("slice[%d]", i))
			if i < aLen && i < bLen {
				c.equals(a.Index(i), b.Index(i), level+1)
			} else if i < aLen {
				c.saveDiff(a.Index(i), "<no value>")
			} else {
				c.saveDiff("<no value>", b.Index(i))
			}
			c.pop()
			if len(c.diff) >= MaxDiff {
				break
			}
		}
		// Float
	case reflect.Float32, reflect.Float64:
		// 소수점 6잘 비교면 충분
		aVal := fmt.Sprintf(c.floatFormat, a.Float())
		bVal := fmt.Sprintf(c.floatFormat, b.Float())

		if aVal != bVal {
			c.saveDiff(a.Float(), b.Float())
		}
		// Boolean
	case reflect.Bool:
		if a.Bool() != b.Bool() {
			c.saveDiff(a.Bool(), b.Bool())
		}
		// Int
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if a.Int() != b.Int() {
			c.saveDiff(a.Int(), b.Int())
		}
		// Uint
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if a.Uint() != b.Uint() {
			c.saveDiff(a.Uint(), b.Uint())
		}
		// String
	case reflect.String:
		if a.String() != b.String() {
			c.saveDiff(a.String(), b.String())
		}

	default:
		logError(ErrNotHandled)
	}
}

// saveDiff - 지정한 a, b 값의 차이점을 기록
func (c *cmp) saveDiff(a, b interface{}) {
	if len(c.buff) > 0 {
		varName := strings.Join(c.buff, ".")
		c.diff = append(c.diff, fmt.Sprintf("%s: %v != %v", varName, a, b))
	} else {
		c.diff = append(c.diff, fmt.Sprintf("%v != %v", a, b))
	}
}

// pop - 관리 중인 buff에서 필드명 추출
func (c *cmp) pop() {
	if len(c.buff) > 0 {
		c.buff = c.buff[0 : len(c.buff)-1]
	}
}

// push - 관리중인 buff에 지정한 필드 추가
func (c *cmp) push(name string) {
	c.buff = append(c.buff, name)
}

// ========== [ Compare END ] =========

// ===== [ Private Functions ] =====

// logError - 지정한 오류 출력
func logError(err error) {
	if LogErrors {
		log.Println(err)
	}
}

// ===== [ Public Functions ] =====

// Equal - 지정한 대상들을 비교하고 다른 점들을 모두 문자 배열로 반환
// returns:
// - []string : 다른 점을 모두 문자열 형식으로 반환
// - nil : a,b 가 모두 nil인 경우, 동일한 경우
// - error : 오류 발생한 경우
// condition
// - strtuct 형식은 재귀적으로 비교 진행
// - struct 형식인 경우에 `deep:"-"` 태그가 존재하는 필드는 비교 생략
// - 해당 형식에 `Equal` 함수가 존재하면 호출
func Equal(a, b interface{}) []string {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	c := &cmp{
		diff:        []string{},
		buff:        []string{},
		floatFormat: fmt.Sprintf("%%.%df", FloatPrecision),
	}

	// nil 비교
	if a == nil && b == nil {
		return nil
	} else if a == nil && b != nil {
		c.saveDiff("<nil pointer>", b)
	} else if a != nil && b == nil {
		c.saveDiff(a, "<nil pointer>")
	}

	// nil 비교 결과가 존재하면 반환
	if len(c.diff) > 0 {
		return c.diff
	}

	// 비교
	c.equals(aVal, bVal, 0)
	if len(c.diff) > 0 {
		return c.diff
	}

	return nil
}
