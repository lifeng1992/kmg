package kmgCaptcha

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/dchest/captcha"
	"net/http"
)

//TODO 可用
func SimpleCaptcha(w http.ResponseWriter) (code string, err error) {
	digits := captcha.RandomDigits(4)
	image := captcha.NewImage(digits, 240, 80)
	w.Header().Set("Content-Type", "image/png")
	_, err := image.WriteTo(w)
	if err != nil {
		return
	}
	//取值范围很小,暂不在此做过多安全性
	code := base64.URLEncoding.EncodeToString(sha256.Sum256(digits)[:])
	return
}

func VerifyCode(code string, input string) (ok bool) {
	return base64.URLEncoding.EncodeToString(sha256.Sum256([]byte(input))[:]) == code
}
