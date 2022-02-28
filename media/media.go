package media

import (
	"github.com/skip2/go-qrcode"
	"net/url"
)

type Qrcode struct {
	Value        string
	Size         int
	SaveFilePath string
}

func (this *Qrcode) Encode() (string, error) {
	if this.Size == 0 {
		this.Size = 256
	}
	var png []byte
	png, err := qrcode.Encode(this.Value, qrcode.Medium, this.Size)
	return string(png), err
}

func (this *Qrcode) WriteFile() error {
	if this.Size == 0 {
		this.Size = 256
	}
	if this.SaveFilePath == "" {
		this.SaveFilePath = "/tmp/tmpQrcode.png"
	}
	return qrcode.WriteFile(this.Value, qrcode.Medium, this.Size, this.SaveFilePath)
}

type Media struct {
	MediaService   string
	UpperAttribute map[string]interface{} //外部属性
}

func (this *Media) FormatUrl(uri, filePath string) string {
	parse, err := url.Parse(uri)
	if err != nil {
		return "URL format parsing error"
	}

	if this.MediaService == "" {
		return "/" + filePath + parse.Path
	}

	return this.MediaService + "/" + filePath + parse.Path
}
