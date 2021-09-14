/*
Copyright 2021 MSFL Authors. All right reserved.
*/
package id

import (
	"errors"
	"net"

	"github.com/ccambo/gocorelib/utils/strings"
	"github.com/sony/sonyflake"
	hashids "github.com/speps/go-hashids"
)

// ===== [ Constants and Variables ] =====
const Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"

var sf *sonyflake.Sonyflake
var upperMachineID uint16

// ===== [ Private Functions ] =====

// init - Called on package load
func init() {
	var st sonyflake.Settings

	// Unique ID Generator 생성
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{
			MachineID: lower16BitIP,
		})
		upperMachineID, _ = upper16BitIP()
	}
}

// lower16BitIP - IP4 주소의 하위 16비트 반환
func lower16BitIP() (uint16, error) {
	ip, err := IPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

// upper16BitIP - IP4 주소의 상위 16비트 반환
func upper16BitIP() (uint16, error) {
	ip, err := IPv4()
	if err != nil {
		return 0, err
	}

	return uint16(ip[0])<<8 + uint16(ip[1]), nil
}

// ===== [ Public Functions ] =====

// GetIntId - Int (uint64) 형식의 UID 생성
// conditions:
// - 생성할 떄 오류 발생시는 Panic 처리
func GetIntId() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

// GetUuid - 지정한 접두어를 포함한 문자열 형식의 UUID 생성
// conditions:
// - 반환 형식은 `B6BZVN3mOPvx...`
// - 생성할 떄 오류 발생시는 Panic 처리
func GetUuid(prefix string) string {
	id := GetIntId()
	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + strings.Reverse(i)
}

// GetUuid36 - 지정한 접두어를 포함한 소문자 기준의 문자열 형식의 UUID 생성
// conditions:
// - 반환 형식은 `300m50zn91nwz5...`
// - 생성할 떄 오류 발생시는 Panic 처리
func GetUuid36(prefix string) string {
	id := GetIntId()
	hd := hashids.NewData()
	hd.Alphabet = Alphabet36
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + strings.Reverse(i)
}

// IPv4 - IP 정보 반환
func IPv4() (net.IP, error) {
	// 시스템의 unicast interface addresses 추출
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		ip := ipnet.IP.To4()
		if ip == nil {
			continue
		}
		return ip, nil

	}
	return nil, errors.New("no ip address")
}
