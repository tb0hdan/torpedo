package file

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"torpedobot/common"
)

func GetCreateChannelDir(channel string) (channelDir string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	channelDirPath := fmt.Sprintf("%s%s%s%s%s", wd, string(os.PathSeparator), "data", string(os.PathSeparator), channel)
	err = os.MkdirAll(channelDirPath, 0755)
	if err == nil {
		channelDir = channelDirPath
	}
	return
}

func GetChannelFile(channel, message string) (channelFile, mimetype string, err error) {
	wd, err := GetCreateChannelDir(channel)
	if err != nil {
		return
	}
	// TODO: Add message permutations
	encoded := base64.URLEncoding.EncodeToString([]byte(strings.TrimSpace(message)))
	fname := fmt.Sprintf("%s%s%s", wd, string(os.PathSeparator), encoded)
	if common.FileExists(fname) {
		mimetype, _, _, err = common.GetMIMEType(fname)
		if err != nil {
			return
		}
		channelFile = fname
	}
	return
}

func SetChannelFile(channel, message string) (result string, err error) {
	wd, err := GetCreateChannelDir(channel)
	if err != nil {
		return
	}
	url_formatted := strings.Split(message, " ")[0]
	url := common.UnformatURL(url_formatted)
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		result = "No valid URL found"
		return
	}
	destination := strings.TrimSpace(strings.TrimLeft(message, url_formatted))
	if destination == "" {
		result = "No valid destination found"
		return
	}
	// Check if target already exists before downloading new one
	encoded := base64.URLEncoding.EncodeToString([]byte(destination))
	new_name := fmt.Sprintf("%s%s%s", wd, string(os.PathSeparator), encoded)
	if common.FileExists(new_name) {
		result = "Destination already exists, set skipped. Use `!rmimg destination` to remove."
		return
	}
	fname, _, is_image, err := common.DownloadToTmp(url)
	if is_image {
		err = os.Rename(fname, new_name)
		if err != nil {
			result = "There was an issue with setting image"
		} else {
			result = "Image set"
		}
	}
	return
}

func ListChannelFiles(channel string) (files []string, err error) {
	wd, err := GetCreateChannelDir(channel)
	if err != nil {
		return
	}

	file_names, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, file := range file_names {
		files = append(files, file.Name())
	}
	return
}

