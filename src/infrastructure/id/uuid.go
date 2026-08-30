package id

import (
	"crypto/rand"
	"fmt"
)

// newUUID は保存前に採番するための UUID を作る。
// DB の自動採番だと保存するまで ID がなく、セッターが必要になるため使わない。
func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	// UUID v4 のバージョン／バリアントビットを立てる
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
