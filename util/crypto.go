package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/base64"
)

// Md5进行签名
func MD5Sign(secret, msg string) (string, error) {
	mac := hmac.New(md5.New, []byte(secret))
	_, err := mac.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func MD5(secret string) (string, error) {
	hash := md5.New()
	hash.Write([]byte(secret))
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// SHA1计算消息摘要
func SHA1(msg string) (string, error) {
	sha := sha1.New()
	_, err := sha.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sha.Sum(nil)), nil
}

func SHA256Sign(secret, msg string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func SHA256SignByte(secret, msg string) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(msg))
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func SHA512Sign(secret, msg string) (string, error) {
	mac := hmac.New(sha512.New, []byte(secret))
	_, err := mac.Write([]byte(msg))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func Sha384Sign(secret, msg string) (string, error) {
	mac := hmac.New(sha512.New384, []byte(secret))
	_, err := mac.Write([]byte(msg))
	if err != nil {
		return "", nil
	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func SHA256Base64Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	signByte := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signByte), nil
}
