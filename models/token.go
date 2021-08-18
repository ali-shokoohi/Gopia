package models

import jwt "github.com/dgrijalva/jwt-go"

/*
Token JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}
