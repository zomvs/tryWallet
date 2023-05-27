package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"strings"
	"time"
)



func main() {
	for i := 0; i < 24; i++ {
		go func() {
			for true {
				privateKey, _ := crypto.GenerateKey()
				publicKeyECDSA, _ := privateKey.Public().(*ecdsa.PublicKey)
				address := crypto.PubkeyToAddress(*publicKeyECDSA).String()

				//pA ,a := address[2:6],"FFFF" // 前缀
				aB ,b := address[38:],"FFFF"   // 后缀 38 = 长度42 减去你需要计算的末位

				if strings.EqualFold(aB,b) {  // 忽略大小写  aB==b 校验大小写
					privateKeyBytes := crypto.FromECDSA(privateKey)
					privateStr := hexutil.Encode(privateKeyBytes)[2:]

					CreateHis(fmt.Sprintf("%s,%s", address, privateStr))
				}
			}
		}()
	}

	time.Sleep(2400 * time.Hour)
}
func CreateHis(strContent string) {
	fd, _ := os.OpenFile("privateKey.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fdTime := time.Now().Format("2006-01-02 15:04:05")
	fdContent := strings.Join([]string{"======", fdTime, "=====", strContent, "\n"}, "")
	buf := []byte(fdContent)
	_, _ = fd.Write(buf)
	_ = fd.Close()
}
