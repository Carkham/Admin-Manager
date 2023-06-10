package utils

import (
	"admin/conf"
	"admin/model"
	UserUser "admin/model/model"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/argon2"
	"log"
	"os"
	"time"
)

type EncodeParam struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func GenerateFromPassword(password string, p EncodeParam) (encodedHash string, err error) {
	salt, err := GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("argon2$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// RSADecrypt RSA解密
func RSADecrypt(cipherText []byte) (string, error) {
	//打开文件
	file, err := os.Open("private.pem")
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)
	//获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	//pem解码
	block, _ := pem.Decode(buf)
	//X509解码
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	//对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		return "", err
	}
	//返回明文
	return string(plainText), nil
}

func DecodeBase64(encoded string) ([]byte, error) {
	encrypt, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return encrypt, nil
}

func CreateAdmin() {
	_, _ = model.Q.UserUser.
		Where(model.Q.UserUser.ID.Eq(1)).
		Or(model.Q.UserUser.Username.Eq(
			conf.Config.Admin.Username,
		)).
		Or(model.Q.UserUser.Username.Eq(
			conf.Config.Admin.Email,
		)).
		Delete()

	now := time.Now()

	// 密码存入数据库前加密
	encodeParam := EncodeParam{
		Memory:      102400,
		Iterations:  2,
		Parallelism: 8,
		SaltLength:  22,
		KeyLength:   32,
	}
	hash, _ := GenerateFromPassword(conf.Config.Admin.Password, encodeParam)

	admin := UserUser.UserUser{
		ID:          1,
		LastLogin:   &now,
		IsSuperuser: true,
		FirstName:   "admin",
		LastName:    "admin",
		IsStaff:     true,
		IsActive:    true,
		//DateJoined:  now,
		Username: conf.Config.Admin.Username,
		Password: hash,
		Email:    conf.Config.Admin.Email,
		Avatar:   "0",
	}

	err := model.Q.UserUser.Create(&admin)

	if err != nil {
		panic(err)
	}

}
