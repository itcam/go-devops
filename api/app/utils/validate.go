package utils

import "github.com/asaskevich/govalidator"

func HasLower(s string) bool {

	for _, c := range s {
		if 'a' <= c && c <= 'z' {
			return true
		}
	}
	return false
}

func HasUpper(s string) bool {

	for _, c := range s {
		if 'A' <= c && c <= 'Z' {
			return true
		}
	}
	return false
}

func HasShuzi(s string) bool {

	for _, c := range s {
		if '0' <= c && c <= '9' {
			return true
		}
	}
	return false
}

func ValidatePassPolicy(str string) {

	//自定义验证函数是否包含大小写和数字 并且长度不小于8
	govalidator.TagMap["PassPolicy"] = govalidator.Validator(func(str string) bool {
		if HasLower(str) && HasUpper(str) && HasShuzi(str) && len(str) >= 8 {
			return true
		}
		return false
	})
}
