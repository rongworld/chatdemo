package main

import (
	"strings"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"fmt"
)

func substring(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		return ""
	}
	if end < 0 || end > length {
		return ""
	}
	return string(rs[start:end])
}

func unicodeIndex(str,substr string) int {
	// 子串在字符串的字节位置
	result := strings.Index(str,substr)
	if result >= 0 {
		// 获得子串之前的字符串并转换成[]byte
		prefix := []byte(str)[0:result]
		// 将子串之前的字符串转换成[]rune
		rs := []rune(string(prefix))
		// 获得子串之前的字符串的长度，便是子串在字符串的字符位置
		result = len(rs)
	}
	return result
}

func getIdFromClaims(key string,claims jwt.Claims) string {
	v:=reflect.ValueOf(claims)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			value := v.MapIndex(k)

			if fmt.Sprintf("%s",k.Interface()) == key{
				return fmt.Sprintf("%v",value.Interface())
			}
		}
	}
	return ""
}