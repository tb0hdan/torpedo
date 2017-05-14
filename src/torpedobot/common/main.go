package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"io"

	"github.com/nlopes/slack"
	"gopkg.in/h2non/filetype.v1"
)

type Utils struct {
	logger *log.Logger
}

func (cu *Utils) SetLogger(logger *log.Logger) {
	cu.logger = logger
	return
}

func (cu *Utils) GetURLBytes(url string) (result []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		cu.logger.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (https://github.com/tb0hdan/torpedo; tb0hdan@gmail.com) Go-http-client/1.1")

	resp, err := client.Do(req)
	if err != nil {
		cu.logger.Fatalln(err)
	}

	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		cu.logger.Fatalln(err)
	}
	return
}

func (cu *Utils) GetMIMEType(fname string) (mimetype, extension string, is_image bool, err error) {
	// Read a file
	buf, err := ioutil.ReadFile(fname)

	if err != nil {
		cu.logger.Printf("GetMIMEType could not read file %s", fname)
		return
	}

	// We only have to pass the file header = first 261 bytes
	head := buf[:261]

	kind, err := filetype.Match(head)
	if err != nil {
		cu.logger.Printf("Mimetype unkwown: %s", err)
		return
	}

	mimetype = kind.MIME.Value
	extension = kind.Extension
	is_image = filetype.IsImage(head)
	return
}

func (cu *Utils) DownloadToTmp(url string) (fname string, mimetype string, is_image bool, err error) {
	img, _ := cu.GetURLBytes(url)
	tmpfile, err := ioutil.TempFile("/tmp", "torpedo")
	if err != nil {
		cu.logger.Fatal(err)
	}

	if _, err := tmpfile.Write(img); err != nil {
		cu.logger.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		cu.logger.Fatal(err)
	}
	fname = tmpfile.Name()
	mimetype, _, is_image, err = cu.GetMIMEType(fname)
	return
}

func GetRequestedFeature(full_command string, usage ...string) (requestedFeature, command, message string) {
	// Support multiple commands within single function
	requestedFeature = strings.Split(full_command, " ")[0]
	command = strings.TrimSpace(strings.TrimLeft(full_command, requestedFeature))
	if len(usage) == 0 {
		message = fmt.Sprintf("Usage: %s string\n", requestedFeature)
	} else {
		message = fmt.Sprintf("Usage: %s %s\n", requestedFeature, usage[0])
	}
	return
}

func ChannelsUploadImage(channels []string, fname, fpath, ftype string, api_i interface{}) {
	parameters := slack.FileUploadParameters{File: fpath, Filetype: ftype,
		Filename: fname, Title: fname,
		Channels: channels}
	api := api_i.(slack.Client)
	api.UploadFile(parameters)
}

func UnformatURL(url string) (newurl string) {
	re := regexp.MustCompile("[<>]")
	newurl = strings.TrimSpace(re.ReplaceAllString(url, ""))
	return
}

func FileExists(fpath string) (exists bool) {
	// TODO: Find a way around this, os.IsExist expects an error and we don't have one yet
	exists = true
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		exists = false
	}
	return
}

func MD5Hash(message string) (result string) {
	my_hash := md5.New()
	io.WriteString(my_hash, message)
	result = fmt.Sprintf("%x", my_hash.Sum(nil))
	return
}

func SHA1Hash(message string) (result string) {
	my_hash := sha1.New()
	io.WriteString(my_hash, message)
	message = fmt.Sprintf("%x", my_hash.Sum(nil))
	return
}

func SHA256Hash(message string) (result string) {
	my_hash := sha256.New()
	io.WriteString(my_hash, message)
	message = fmt.Sprintf("%x", my_hash.Sum(nil))
	return
}

func SHA512Hash(message string) (result string) {
	my_hash := sha512.New()
	io.WriteString(my_hash, message)
	message = fmt.Sprintf("%x", my_hash.Sum(nil))
	return
}
