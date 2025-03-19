package pkg_utils

import (
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwt"
)

type Jwt struct {
	secret   string
	expires  int
	_jwtauth *jwtauth.JWTAuth
}

func NewJwt(secret string, expires int) *Jwt {
	return &Jwt{
		secret:   secret,
		expires:  expires,
		_jwtauth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

func (j *Jwt) Generate(claims map[string]interface{}) (string, error) {

	_, stringToken, jwterr := j._jwtauth.Encode(claims)

	if jwterr != nil {
		return "", jwterr
	}

	return stringToken, nil
}

func (j *Jwt) Decode(tokenString string) jwt.Token {
	token, _ := j._jwtauth.Decode(tokenString)
	return token
}

func (j *Jwt) VerifyToken(tokenString string) error {
	if _, e := jwtauth.VerifyToken(j._jwtauth, tokenString); e != nil {
		return e
	}

	return nil
}

func (l *Jwt) ExtractStringValue(token jwt.Token, key string) string {
	if v, e := token.Get(key); e {
		return v.(string)
	}
	return ""
}

func (l *Jwt) ExtractIntValue(token jwt.Token, claimName string, defaultValue int) int {

	value := defaultValue
	if v, e := token.Get(claimName); e {
		value = int(v.(float64))
	}
	return value
}
