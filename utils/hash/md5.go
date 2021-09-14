/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package hash

import (
	"encoding/hex"
	"io"

	"code.cloudfoundry.org/bytefmt"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// GetMD5 - 지정한 Reader의 값을 MD5 값으로 변환해서 반환
func GetMD5(reader io.ReadCloser) (string, error) {
	md5Reader := NewMD5Reader(reader)
	data := make([]byte, bytefmt.KILOBYTE)

	for {
		_, err := md5Reader.Read(data)
		if err != nil {
			// reader에 더 이상 읽을 내용이 없는 경우
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	err := reader.Close()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(md5Reader.MD5()), nil
}
