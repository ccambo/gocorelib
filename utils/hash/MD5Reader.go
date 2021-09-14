/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package hash

import (
	"crypto/md5"
	"hash"
	"io"
)

// ===== [ Constants and Variables ] =====
const ()

var ()

// ========== [ MD5Reader START ] =========

// MD5Reader - MD5Reader 정보 관리용
type MD5Reader struct {
	md5  hash.Hash
	body io.Reader
}

// Read - 지정한 byte배열에 관리 중인 정보를 MD5로 출력하고 크기와 오류를 반환
func (r *MD5Reader) Read(b []byte) (int, error) {
	n, err := r.body.Read(b)
	if err != nil {
		return n, err
	}
	return r.md5.Write(b[:n])
}

// MD5 - 관리 중인 데이터를 MD5 byte 배열로 반환
func (r *MD5Reader) MD5() []byte {
	return r.md5.Sum(nil)
}

// ========== [ MD5Reader END ] =========

// ===== [ Private Functions ] =====
// ===== [ Public Functions ] =====

// NewMD5Reader - 지정한 io.Reader를 기반으로 MD5 Reader 생성
func NewMD5Reader(reader io.Reader) *MD5Reader {
	return &MD5Reader{
		md5:  md5.New(),
		body: reader,
	}
}
