package pinterest

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"torpedobot/common"
	"os"
	"log"
)


type PinterestInterestItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}


type PinterestPage struct {
	Cursor   string `json:"cursor"`
	NextPage string `json:"next"`
}


type PinterestInterests struct {
	Data []*PinterestInterestItem `json:"data"`
	Page *PinterestPage           `json:"page"`
}


type PinterestBoardPinItem struct {
	URL  string `json:"url"`
	Note string `json:"note"`
	Link string `json:"link"`
	ID   string `json:"id"`
}


type PinterestBoardPins struct {
	Data []*PinterestBoardPinItem `json:"data"`
	Page *PinterestPage           `json:"page"`
}


type PinterestPinMedia struct {
	Type string `json:"type"`
}


type PinterestPinOriginalImage struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}


type PinterestPinImage struct {
	Original *PinterestPinOriginalImage `json:"original"`
}


type PinterestPinItem struct {
	Media     *PinterestPinMedia `json:"media"`
	CreatedAt string             `json:"created_at"`
	Image     *PinterestPinImage `json:"image"`
	ID        string             `json:"id"`
}


type PinterestPinResponse struct {
	Data *PinterestPinItem `json:"data"`
}


type PinterestOauthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       []string
}


const PINTEREST_API_BASE = "https://api.pinterest.com/v1"


type PinterestClient struct {
	config struct {
		client_id     string
		client_secret string
		client_token  string
		redirect_url  string
		state         string
	}
	logger *log.Logger
	utils *common.Utils
}


func (api *PinterestClient) GetAcccessToken(code string) (client_token string, err error) {
	var pinterestOauthResponse PinterestOauthResponse

	url_tpl := fmt.Sprintf("%s/oauth/token?", PINTEREST_API_BASE)
	url_tpl += "grant_type=authorization_code&"
	url_tpl += "client_id=%s&"
	url_tpl += "client_secret=%s&"
	url_tpl += "code=%s"
	url := fmt.Sprintf(url_tpl, api.config.client_id, api.config.client_secret, code)
	r, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewReader([]byte{}))
	defer r.Body.Close()
	token, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(token, &pinterestOauthResponse)
	if err == nil {
		client_token = pinterestOauthResponse.AccessToken
	}
	return
}


func (api *PinterestClient) PrepareRequest(state string) {
	url_tpl := "https://api.pinterest.com/oauth/?"
	url_tpl += "response_type=code&"
	url_tpl += "redirect_uri=%s&"
	url_tpl += "client_id=%s&"
	url_tpl += "scope=read_public,read_relationships&"
	url_tpl += "state=%s"
	url := fmt.Sprintf(url_tpl, "https://localhost:1443/pinterest", api.config.client_id, state)
	api.logger.Println(url)
}


func (api *PinterestClient) GetInterests() {
	var pinterestInterests PinterestInterests
	url := fmt.Sprintf("%s/me/following/interests/?access_token=%s", PINTEREST_API_BASE, api.config.client_token)
	response, err := api.utils.GetURLBytes(url)
	err = json.Unmarshal(response, &pinterestInterests)
	if err == nil {
		for _, item := range pinterestInterests.Data {
			api.logger.Printf("%s - %s\n", item.ID, item.Name)
		}
	}
}


func (api *PinterestClient) GetPinsForBoard(board string) (pins []string, err error) {
	var pinterestBoardPins PinterestBoardPins
	url := fmt.Sprintf("https://api.pinterest.com/v1/boards/%s/pins/?access_token=%s", board, api.config.client_token)
	response, err := api.utils.GetURLBytes(url)
	err = json.Unmarshal(response, &pinterestBoardPins)
	if err == nil {
		for _, item := range pinterestBoardPins.Data {
			pins = append(pins, item.ID)
		}
	}
	return
}


func (api *PinterestClient) GetImagesForPin(pin string) (images []string, err error) {
	var pinterestPinResponse PinterestPinResponse
	url := fmt.Sprintf("https://api.pinterest.com/v1/pins/%s/?access_token=%s&fields=media,image,created_at", pin, api.config.client_token)
	response, err := api.utils.GetURLBytes(url)
	if err != nil {
		api.logger.Printf("%+v", err)
		return
	}
	err = json.Unmarshal(response, &pinterestPinResponse)
	if err == nil {
		images = []string{pinterestPinResponse.Data.Image.Original.URL}
	} else {
		api.logger.Printf("%+v", err)
	}
	return
}


func (api *PinterestClient) GetImagesForBoard(board string) (images []string, err error) {
	pins, err := api.GetPinsForBoard(board)
	if err != nil {
		api.logger.Printf("Get pins for board failed with %+v", err)
		return
	}
	for _, pin := range pins {
		pin_images, err := api.GetImagesForPin(pin)
		if err != nil {
			api.logger.Printf("Get images failed for pin %s with %+v", pin, err)
			break
		}
		images = append(images, pin_images[0])

	}
	return
}


func (api *PinterestClient) RandState() (err error) {
	c := 16
	b := make([]byte, c)
	_, err = rand.Read(b)
	if err != nil {
		api.logger.Println("error:", err)
		return
	}
	encoded := base64.URLEncoding.EncodeToString(b)
	api.config.state = strings.Replace(encoded, "==", "", -1)
	return
}


func New(client_token string) (client *PinterestClient) {
	client = &PinterestClient{}
	client.logger = log.New(os.Stdout, "pinterest-plugin: ", log.Lshortfile|log.LstdFlags)
	client.utils = &common.Utils{}
	client.utils.SetLogger(client.logger)
	client.config.client_token = client_token
	return
}
