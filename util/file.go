package util

import (
	"archive/zip"
	"io"
	"os"
)

func WriteToFile(filePath string, data []byte) error {
	fileHandler, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fileHandler.Close()
	_, err = fileHandler.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func Zip(filePath, savePath string) error {
	archive, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	fw, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fw.Close()

	zw := zip.NewWriter(archive)
	w, err := zw.Create(savePath)
	if err != nil {
		return err
	}
	if _, err = io.Copy(w, fw); err != nil {
		return err
	}

	return zw.Close()
}
