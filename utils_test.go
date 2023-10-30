package utils

import (
	"fmt"
	"testing"
)

func TestCPUID(t *testing.T) {
	fmt.Println("cpuid", CPUID())
}

func TestMap(t *testing.T) {
	failed := &Counter{Data: make(map[string]int)}
	failed.Incr("hello", 3)
	failed.Incr("hello", 1)
	fmt.Println(failed.Get("hello"))
}

func TestAesEnc(t *testing.T) {
	res := AesEncStr("hello world", "1q2w3e4r5t6y7u8i")
	fmt.Println("aes encrypt", res)

	text := AesDecStr(res, "1q2w3e4r5t6y7u8i")
	fmt.Println("aes decrypt", text)
}

func TestDesEnc(t *testing.T) {
	res := DesEncStr("hello world", "1q2w3e4r")
	fmt.Println("des encrypt", res)

	text := DesDecStr(res, "1q2w3e4r")
	fmt.Println("des decrypt", text)
}

func TestHash(t *testing.T) {
	fmt.Println(RandomString(5, "hello world"))
}
