package apps

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/adinandradrs/cezbek-engine/internal/model"
	"github.com/sethvargo/go-password/password"
	"go.uber.org/zap"
	"io"
	"math/big"
)

func RandomOtp(d int) (string, error) {
	ns := "012345679"
	b := make([]byte, d)
	for i := 0; i < d; i++ {
		max := big.NewInt(int64(len(ns)))
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = ns[num.Int64()]
	}
	return string(b), nil
}

func RandomPassword(len int, d int, sym int, logger *zap.Logger) (res string, e *model.TechnicalError) {
	res, err := password.Generate(len, d, sym, false, false)
	if err != nil {
		return res, Exception("failed to generate password", err, zap.Int("d", d), logger)
	}
	logger.Info("success generate password", zap.String("generated", res))
	return res, e
}

func Hash(key string) string {
	sha256 := sha256.New()
	sha256.Write([]byte(key))
	return hex.EncodeToString(sha256.Sum(nil))
}

func Encrypt(d string, h string, logger *zap.Logger) (res []byte, ex *model.TechnicalError) {
	c, _ := aes.NewCipher([]byte(h))
	o, err := cipher.NewGCM(c)
	if err != nil {
		return res, Exception("failed to encrypt GCM data", err, zap.String("h", h), logger)
	}
	nsz := make([]byte, o.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nsz); err != nil {
		return res, Exception("failed to encrypt I/O reader", err, zap.String("h", h), logger)
	}
	return o.Seal(nsz, nsz, []byte(d), nil), ex
}

func Decrypt(data []byte, hash string, logger *zap.Logger) (res string, ex *model.TechnicalError) {
	c, err := aes.NewCipher([]byte(hash))
	if err != nil {
		return res, Exception("failed to decrypt chiper", err, zap.String("hash", hash), logger)
	}
	aead, err := cipher.NewGCM(c)
	if err != nil {
		return res, Exception("failed to decrypt init GCM", err, zap.String("hash", hash), logger)
	}
	nsz := aead.NonceSize()
	n, cbytes := data[:nsz], data[nsz:]
	o, err := aead.Open(nil, n, cbytes, nil)
	if err != nil {
		return res, Exception("failed to decrypt nonce size parse logic", err, zap.String("hash", hash), logger)
	}
	return string(o), ex
}