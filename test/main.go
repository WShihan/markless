package main

import (
	"encoding/json"
	"fmt"
	"time"

	jose "gopkg.in/square/go-jose.v2"
)

func createAndEncryptJWT(secretKey []byte) (string, error) {
	// 创建 JWT 载荷
	payload := map[string]interface{}{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  time.Now().Unix(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// 创建 JWE 加密器
	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       secretKey,
		},
		nil,
	)
	if err != nil {
		return "", err
	}

	// 加密 JWT
	jwe, err := encrypter.Encrypt(payloadBytes)
	if err != nil {
		return "", err
	}

	serializedJWE := jwe.FullSerialize()
	return serializedJWE, nil
}

func decryptJWT(encryptedJWT string, secretKey []byte) (map[string]interface{}, error) {
	// 解析 JWE
	jwe, err := jose.ParseEncrypted(encryptedJWT)
	if err != nil {
		return nil, err
	}

	// 解密 JWE
	decryptedBytes, err := jwe.Decrypt(secretKey)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	err = json.Unmarshal(decryptedBytes, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func main() {
	// 对称密钥 (32 bytes for AES-256)
	secretKey := []byte("14280907969981547271211133197484")

	encryptedJWT, err := createAndEncryptJWT(secretKey)
	if err != nil {
		fmt.Println("Error creating and encrypting JWT:", err)
		return
	}

	fmt.Println("Encrypted JWE:", encryptedJWT)

	decryptedPayload, err := decryptJWT(encryptedJWT, secretKey)
	if err != nil {
		fmt.Println("Error decrypting JWT:", err)
		return
	}

	fmt.Println("Decrypted Payload:", decryptedPayload["sub"])
}
