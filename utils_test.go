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

func TestXor(t *testing.T) {
    println(XorEncStr("https://image.zyh365.com/upload/schoolpics/20230420/20230420142539_425.jpg", 26))
    println(AesEncStr("https://image.zyh365.com/upload/schoolpics/20230420/20230420142539_425.jpg", "1234567812345678"))
    println(DesEncStr("https://image.zyh365.com/upload/schoolpics/20230420/20230420142539_425.jpg", "12345678"))
}

func BenchmarkXor(b *testing.B) {
	for n := 0; n < b.N; n++ {
		XorDecStr("cm5uamkgNTVzd3t9fzRgY3IpLC80eXV3NW9qdnV7fjVpeXJ1dXZqc3lpNSgqKCkqLigqNSgqKCkqLigqKy4oLykjRS4oLzRwan0=", 26)
	}
}

func BenchmarkAes(b *testing.B) {
	for n := 0; n < b.N; n++ {
        AesDecStr("sbGYj3AFPMRZFuFtSNADA8Ay+bjp7J0/JQP2VGvvMFVVzY3bx8vCjhCq6It7JuyFSqFP4wttX57l8CegN9rDONzpk6Ou7to/iC4N4FxaBTU=", "1234567812345678")
	}
}

func BenchmarkDes(b *testing.B) {
	for n := 0; n < b.N; n++ {
        DesDecStr("bsrIGBPKq/ogXAiMrr8X1BGr8XY1deq+YIhW9K8r0HkxJ9ioV6cr/kZJBusvEZdcaCoQlLecGDKWuQ+2IvZslE+CC8epLWq7aAHb2UngP/o=", "12345678")
	}
}
