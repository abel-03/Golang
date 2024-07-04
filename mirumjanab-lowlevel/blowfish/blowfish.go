//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <openssl/blowfish.h>
import "C"
import "unsafe"

const block_size = 8

type Blowfish struct {
	key C.BF_KEY
}

func New(key []byte) *Blowfish {
	blow_fish := Blowfish{}
	C.BF_set_key(&blow_fish.key, (C.int)(len(key)), (*C.uchar)(unsafe.Pointer(&key[0])))
	return &blow_fish
}

func (b *Blowfish) Encrypt(dst, src []byte) {
    if dst == nil || src == nil ||  len(src) != len(dst) {
        return 
    }

    C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&src[0])), (*C.uchar)(unsafe.Pointer(&dst[0])), &b.key, C.BF_ENCRYPT)

}

func (b *Blowfish) Decrypt(dst, src []byte) {
    if dst == nil || src == nil || len(src) != len(dst){
        return
    }

    C.BF_ecb_encrypt((*C.uchar)(unsafe.Pointer(&src[0])), (*C.uchar)(unsafe.Pointer(&dst[0])), &b.key, C.BF_DECRYPT)
}

func (b *Blowfish) BlockSize() int {
	return block_size
}