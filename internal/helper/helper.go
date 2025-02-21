package helper

import (
	"crypto/rand"
	"io"
	"os"
)

func GetRandomString(length int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i := 0; i < length; i++ {
		bytes[i] = letters[bytes[i]%byte(len(letters))]
	}
	return string(bytes)
}

func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
