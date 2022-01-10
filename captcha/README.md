# captcha

图形验证码

Demo:

```go
package main

import (
	"github.com/junbin-yang/golib/captcha"
    "time"
)

func main() {
	(&captcha.Options{Expiration: 15 * time.Minute}).New()
    // 生成验证码
	id, b64s, value, err := captcha.GenerateDefaultCaptcha(false)
    ...
    
    // 验证
    if captcha.Verify(id, value) {
    	...    
    }
}

```

