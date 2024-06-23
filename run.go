package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	passHash := "$bitcoin$64$fb795cc23f19bb90fdd102da017d173f2078c4bd632c066825b8a9e2f9223aa2$16$1751d7bfe05f388c$452956$2$00$2$00"
	passHash = passHash[9:]                       // 64$12c098515dc4f4140786e352f05d3065f17a2ca8f15c5f1c93923dc7146380c6$16$146b99a74fa7b536$135174$2$00$2$00"
	passHashArray := strings.Split(passHash, "$") // [64 12c098515dc4f4140786e352f05d3065f17a2ca8f15c5f1c93923dc7146380c6 16 146b99a74fa7b536 135174 2 00 2 00]
	passToHashTmp, _ := hex.DecodeString(passHashArray[3])
	for i := 0; i < 1000; i++ {
		passphrase := "123"
		passToHash := []byte(passphrase + string(passToHashTmp))

		num, _ := strconv.Atoi(passHashArray[4])
		for i := 0; i < num; i++ {
			passToHash = sha512Hash(passToHash)
		}

		key, iv := passToHash[0:32], passToHash[32:48]

		// aes-256-cbc 解密
		deStr, _ := hex.DecodeString(passHashArray[1])
		block, err := aes.NewCipher(key)
		if err != nil {
			fmt.Println("Error creating AES cipher:", err)
			return
		}
		mode := cipher.NewCBCDecrypter(block, iv)
		decryptedData := make([]byte, len(deStr))
		mode.CryptBlocks(decryptedData, deStr)

		decryptedData, rebool := pkcs7Unpad(decryptedData)
		if rebool == true {
			fmt.Println("ok~:", passphrase)
			CreateHis(passphrase)
		} else {
			// fmt.Println("c")
			continue
		}
	}

	elapsed := time.Since(start)
	fmt.Println("程序执行耗时:", elapsed)
}

func sha512Hash(data []byte) []byte {
	hasher := sha512.New()
	_, err := io.WriteString(hasher, string(data))
	if err != nil {
		panic(err)
	}
	return hasher.Sum(nil)
}

// pkcs7Unpad 去除PKCS7填充
func pkcs7Unpad(data []byte) ([]byte, bool) {
	length := len(data)
	unpadding := int(data[length-1])
	if length > unpadding {
		return data[:(length - unpadding)], true
	} else {
		return nil, false
	}
}

func CreateHis(strContent string) {
	fd, _ := os.OpenFile("privateKey.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fdTime := time.Now().Format("2006-01-02 15:04:05")
	fdContent := strings.Join([]string{"======", fdTime, "=====", strContent, "\n"}, "")
	buf := []byte(fdContent)
	_, _ = fd.Write(buf)
	_ = fd.Close()
}
