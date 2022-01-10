package captcha

import (
	"github.com/mojocn/base64Captcha"
	"time"
)

type Options struct {
	// 验证码存储空间长度
	LimitNumber int
	// 验证码过期时间
	Expiration time.Duration
}

var store base64Captcha.Store

func (this *Options) New() {
	if this.LimitNumber == 0 {
		this.LimitNumber = 10240
	}
	if this.Expiration == 0 {
		this.Expiration = 10 * time.Minute
	}
	store = base64Captcha.NewMemoryStore(this.LimitNumber, this.Expiration)
}

func NewDriver(h, w, l int) *base64Captcha.DriverString {
	driver := new(base64Captcha.DriverString)
	driver.Height = h
	driver.Width = w
	driver.NoiseCount = 10
	//driver.ShowLineOptions = base64Captcha.OptionShowSineLine | base64Captcha.OptionShowSlimeLine | base64Captcha.OptionShowHollowLine
	driver.ShowLineOptions = base64Captcha.OptionShowHollowLine
	driver.Length = l
	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
	driver.Fonts = []string{"wqy-microhei.ttc"}
	return driver
}

// 生成验证码
func GenerateDefaultCaptcha(clear bool) (string, string, string, error) {
	return GenerateCaptcha(40, 120, 4, clear)
}

func GenerateCaptcha(height, width, length int, clear bool) (string, string, string, error) {
	var driver = NewDriver(height, width, length).ConvertFonts()
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := c.Generate()
	if err != nil {
		return "", "", "", err
	}
	value := store.Get(id, clear)
	return id, b64s, value, nil
}

// 验证验证码
func Verify(id, val string) bool {
	if id == "" || val == "" {
		return false
	}
	// 同时在内存清理掉这个图片
	return store.Verify(id, val, true)
}
