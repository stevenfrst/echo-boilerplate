package encoder

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func EncodeUrlEmailVerify(str string) string {
	encode := base64.StdEncoding.EncodeToString([]byte(str))
	encodeEncrypt := base64.StdEncoding.EncodeToString([]byte(str + viper.GetString(`encrypt.additional`)))
	log.Println(encodeEncrypt, viper.GetString(`encrypt.keystring`))
	return fmt.Sprintf(viper.GetString(`server.address.frontend`)+"/api/v1/verify?u=%s&v=%s", encode, AesEncrypt(encodeEncrypt, viper.GetString(`encrypt.keystring`)))
}

func DecodeEmailVerify(email string, verify string) (string, error) {
	decode, err := base64.StdEncoding.DecodeString(email)
	if err != nil {
		return "", err
	}
	decrypt, err := AesDecrypt(verify, viper.GetString(`encrypt.keystring`))
	if err != nil {
		return "", err
	}
	decodeDecrypt, err := base64.StdEncoding.DecodeString(decrypt)
	if err != nil {
		return "", err
	}
	decodeStr := string(decodeDecrypt)
	decodeStr = strings.ReplaceAll(decodeStr, viper.GetString(`encrypt.additional`), "")
	if decodeStr == string(decode) {
		return decodeStr, nil
	}
	return "", nil
}
