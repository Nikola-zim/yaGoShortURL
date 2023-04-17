package cash

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type AuthUser struct {
	users  map[uint64][]byte // ключ - id; значение - ключ(для расшифровки)
	lastId uint64
}

func (aU *AuthUser) FindUser(idMsg string) (uint64, bool) {
	data, err := hex.DecodeString(idMsg)
	if err != nil {
		return 0, false
	}
	// достаем id
	id := binary.LittleEndian.Uint64(data[:8])
	// по id ищем ключ для расшифровки
	secretKey, ok := aU.users[id]
	if !ok {
		return 0, false
	}
	h := hmac.New(sha256.New, secretKey)
	h.Write(data[:8])
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[8:]) {
		return id, true
	} else {
		return id, false
	}
}

func (aU *AuthUser) AddUser() (string, error) {
	secretKey, err := generateRandom(24)
	if err != nil {
		return "", err
	}
	// запись в мапу для идентификации
	aU.users[aU.lastId+1] = secretKey
	// подписываем алгоритмом HMAC, используя SHA256
	h := hmac.New(sha256.New, secretKey)
	// слайс байт содержащий id (8 байт) и подпись
	cookieByte := make([]byte, 8)
	// запись id для куки и для шифровки
	binary.LittleEndian.PutUint64(cookieByte, aU.lastId+1)
	h.Write(cookieByte)
	dst := h.Sum(nil)

	cookieByte = append(cookieByte, dst...)
	cookieStr := hex.EncodeToString(cookieByte)
	return cookieStr, nil
}

// генератор вектора для секретного ключа
func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func NewAuthUser() *AuthUser {
	return &AuthUser{
		users:  make(map[uint64][]byte),
		lastId: 0,
	}
}
