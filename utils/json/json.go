/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package json

import (
	"encoding/json"
	"strings"

	"k8s.io/klog"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ========== [ JsonRawMessage START ] =========

// JsonRawMessage - raw json message bytes

type JsonRawMessage []byte

// Find - 관리 중인 Json Mesages에 대해 지정한 키에 해당하는 Json Message 반환
// conditions:
// - 내부 처리 중 오류가 발생한 경우는 nil 반환
func (m JsonRawMessage) Find(key string) JsonRawMessage {
	var objmap map[string]json.RawMessage

	err := json.Unmarshal(m, &objmap)
	if err != nil {
		klog.Errorf("Resolve JSON Key failed, find key=%s, err=$s", key, err)
		return nil
	}

	return JsonRawMessage(objmap[key])
}

// ToList - 관리 중인 Json Message의 리스트 반환
func (m JsonRawMessage) ToList() []JsonRawMessage {
	var lists []json.RawMessage

	err := json.Unmarshal(m, &lists)
	if err != nil {
		klog.Errorf("Resolve JSON List failed, err=$s", err)
		return nil
	}

	var res []JsonRawMessage
	for _, v := range lists {
		res = append(res, JsonRawMessage(v))
	}

	return res
}

// ToString - 관리 중인 Json Message를 문자열로 반환
func (m JsonRawMessage) ToString() string {
	res := strings.Replace(string(m[:]), "\"", "", -1)
	return res
}

// ========== [ JsonRawMessage END ] =========

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====
