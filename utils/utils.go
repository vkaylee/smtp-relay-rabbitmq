package utils

import (
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

func ErrFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func ErrPrintln(err error) {
	if err != nil {
		log.Println(err)
	}
}
func DownloadFile(URL string) (string, error) {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	mime := mimetype.Detect(body)
	// Gen random file name
	fileName := fmt.Sprintf("attachment-%d%s", rand.Uint64(), mime.Extension())
	// Only allow this type
	allowed := []string{"image/jpeg", "image/png", "image/gif"}
	if !mimetype.EqualsAny(mime.String(), allowed...) {
		return "", err
	}
	// Write file to disk
	err = ioutil.WriteFile(fileName, body, 0777)
	if err != nil {
		return "", err
	}
	return fileName, nil
}