package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"reflect"
	"runtime"

	"github.com/google/uuid"
	"github.com/super-l/machine-code/machine"
)

func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func Getenv(args ...string) string {
	v := os.Getenv(args[0])
	if v == "" && len(args) == 2 {
		return args[1]
	}
	return v
}

func LocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("Can not find the client ip address")
}

func RandomString(size int, charsets ...string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := []byte{}
	if len(charsets) == 0 {
		charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	} else {
		for i := 0; i < len(charsets); i++ {
			for j := 0; j < len(charsets[i]); j++ {
				charset = append(charset, charsets[i][j])
			}
		}
	}
	bytes := make([]byte, size)
	count := len(charset)
	for i := 0; i < size; i++ {
		bytes[i] = charset[r.Intn(count)]
	}
	return string(bytes)
}

func RandomIP() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func ToBytes(input interface{}) []byte {
	switch input.(type) {
	case string:
		return []byte(input.(string))
	case []byte:
		return input.([]byte)
	default:
		return []byte{}
	}
}

func ToString(input interface{}) string {
	switch input.(type) {
	case string:
		return input.(string)
	case []byte:
		return string(input.([]byte))
	default:
		return ""
	}
}

func Md5(input interface{}) string {
	bytes := md5.Sum(ToBytes(input))
	return hex.EncodeToString(bytes[:])
}

func Sha1(input interface{}) string {
	bytes := sha1.Sum(ToBytes(input))
	return hex.EncodeToString(bytes[:])
}

func FileHash(file string, h hash.Hash) string {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return ""
	}
	io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil))
}

func FileMd5(file string) string {
	return FileHash(file, md5.New())
}

func FileSha1(file string) string {
	return FileHash(file, sha1.New())
}

func UUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func CPUID() string {
	uuid, _ := machine.GetPlatformUUID()
	mac, _ := machine.GetMACAddress()
	info := fmt.Sprintf("%s %s", uuid, mac)
	return Md5(info)
}

func URLJoin(base, href string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return href
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return href
	}
	return baseURL.ResolveReference(uri).String()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return Exists(path) && !IsDir(path)
}

func ResetColor() {
	fmt.Printf("%c[0m", 0x1B)
}

func String2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ReadFiles(root string) []string {
	var files []string
	rd, _ := ioutil.ReadDir(root)
	for _, fi := range rd {
		path := filepath.Join(root, fi.Name())
		if fi.IsDir() {
			files = append(files, ReadFiles(path)...)
		} else if fi.Mode()&os.ModeSymlink != 0 {
			path, err := os.Readlink(path)
			if err == nil {
				if IsDir(path) {
					files = append(files, ReadFiles(path)...)
				} else {
					files = append(files, path)
				}
			}
		} else {
			files = append(files, path)
		}
	}
	return files
}
