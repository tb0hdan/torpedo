package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"gopkg.in/h2non/filetype.v1"
)

func GetURLBytes(url string) (result []byte, err error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (https://github.com/tb0hdan/torpedo; tb0hdan@gmail.com) Go-http-client/1.1")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func GetMIMEType(fname string) (mimetype, extension string, is_image bool, err error) {
	// Read a file
	buf, err := ioutil.ReadFile(fname)

	if err != nil {
		fmt.Printf("GetMIMEType could not read file %s", fname)
		return
	}

	// We only have to pass the file header = first 261 bytes
	head := buf[:261]

	kind, err := filetype.Match(head)
	if err != nil {
		fmt.Printf("Mimetype unkwown: %s", err)
		return
	}

	mimetype = kind.MIME.Value
	extension = kind.Extension
	is_image = filetype.IsImage(head)
	return
}

func DownloadToTmp(url string) (fname string, mimetype string, is_image bool, err error) {
	img, _ := GetURLBytes(url)
	tmpfile, err := ioutil.TempFile("/tmp", "torpedo")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(img); err != nil {
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	fname = tmpfile.Name()
	mimetype, _, is_image, err = GetMIMEType(fname)
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

func ChannelsUploadImage(channels []string, fname, fpath, ftype string, api *slack.Client) {
	parameters := slack.FileUploadParameters{File: fpath, Filetype: ftype,
		Filename: fname, Title: fname,
		Channels: channels}
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
