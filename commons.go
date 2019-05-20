package common

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SaveImage(imgFile multipart.File, savePath string) (savedPath string) {
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(imgFile)
	savedFileName := GetRandomHashString()
	buffer := &bytes.Buffer{}
	_, err := io.Copy(buffer, imgFile)
	if err != nil {
		panic(err)
	}
	fileBytes, _ := ioutil.ReadAll(buffer)
	headData := fileBytes[:512]

	if contentType := http.DetectContentType(headData); contentType != "" {
		switch contentType {
		case "image/gif":
			savedFileName += ".gif"
			break
		case "image/png":
			savedFileName += ".png"
			break
		case "image/jpg":
			savedFileName += ".jpg"
			break
		case "image/jpeg":
			savedFileName += ".jpg"
			break
		}
	}
	if _, err = os.Stat(savePath); os.IsNotExist(err) {
		err = os.MkdirAll(savePath, 0755)
		if err != nil {
			panic(err)
		}
	}

	savedPath = fmt.Sprintf("%s/%s", savePath, savedFileName)
	err = ioutil.WriteFile(savedPath, fileBytes, 0644)

	if err != nil {
		panic(err)
	}

	return
}

func GetRandomHashString() (str string) {
	u1 := uuid.NewV4()
	str = u1.String()
	str += strconv.FormatInt(time.Now().UnixNano(), 10)
	str = fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
	h := md5.New()
	h.Write([]byte(str))
	str = fmt.Sprintf("%x", h.Sum(nil))

	return
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
