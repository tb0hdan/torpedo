package bashim

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	"torpedobot/common"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)


type BashIM struct {
	logger *log.Logger
	utils  *common.Utils
}


func (bi *BashIM) windows1251_to_utf8(reader_in io.Reader) (reader_out *transform.Reader) {
	reader_out = transform.NewReader(reader_in, charmap.Windows1251.NewDecoder())
	return
}


func (bi *BashIM) windows1251_to_utf8_bytes_reader(input []byte) (output *bytes.Reader, err error) {
	// TODO: Well, this is obviusly effed up
	result, err := ioutil.ReadAll(bi.windows1251_to_utf8(bytes.NewReader(input)))
	if err != nil {
		bi.logger.Printf("windows1251_to_utf8_bytes_reader failed with %+v", err)
	}
	output = bytes.NewReader(result)
	return
}


func (bi *BashIM) Get_html(url string) (result *html.Node) {
	res, err := bi.utils.GetURLBytes(url)
	if err != nil {
		return
	}
	reader, err := bi.windows1251_to_utf8_bytes_reader(res)

	result, err = html.Parse(reader)
	if err != nil {
		bi.logger.Printf("Error %s parsing html", err)
	}
	return
}


func NewClient() (client *BashIM) {
	client = &BashIM{}
	client.utils = &common.Utils{}
	client.logger = client.utils.SetLoggerPrefix("bashim-plugin")
	return
}
