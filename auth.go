package main

import (
	"net/http"
	"errors"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

const SecretKey  = "star00star"


type T struct {}


func (t T)  ExtractToken(req *http.Request) (s string, err error) {
	//t实现Extractor接口，获取token的值
	token := req.FormValue("token")
	if token == ""{
		err = errors.New("token is nil")
		return 
	}else{
		return token,nil
	}
}


func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)  {
	token, err := getToken(r)
	if err != nil{
		fmt.Println(err)
	}else {
		if token.Valid {
			next(w,r)
			fmt.Println("验证通过")
		}else {
			fmt.Fprint(w,"验证失败")
			fmt.Println("验证失败")
		}


	}
}

func getToken(r *http.Request) (token *jwt.Token, err error) {//由request获取token
	t := T{}
	return request.ParseFromRequest(r, t,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})
}