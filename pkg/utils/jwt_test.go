package pkg_utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _jwt = NewJwt("SECRET", 500)

func TestExtractKey(t *testing.T) {

	claims := map[string]interface{}{
		"key": "rene.oliveira",
		"exp": time.Now().Add(time.Minute * time.Duration(_jwt.expires)).Unix(),
	}

	tokenString, e := _jwt.Generate(claims)

	assert.Nil(t, e)

	token := _jwt.Decode(tokenString)

	value := _jwt.ExtractStringValue(token, "key")

	assert.Equal(t, "rene.oliveira", value)
}

func TestExtractIntValue(t *testing.T) {

	claims := map[string]interface{}{
		"test": 9,
		"exp":  time.Now().Add(time.Minute * time.Duration(_jwt.expires)).Unix(),
	}

	tokenString, e := _jwt.Generate(claims)

	assert.Nil(t, e)

	token := _jwt.Decode(tokenString)

	value := _jwt.ExtractIntValue(token, "test", 0)

	assert.Equal(t, 9, value)
}

func TestExtractIntDefaultValue(t *testing.T) {

	claims := map[string]interface{}{
		"exp": time.Now().Add(time.Minute * time.Duration(_jwt.expires)).Unix(),
	}

	tokenString, e := _jwt.Generate(claims)

	assert.Nil(t, e)

	token := _jwt.Decode(tokenString)

	value := _jwt.ExtractIntValue(token, "test", 9)

	assert.Equal(t, 9, value)
}

func TestValidateExpiredToken(t *testing.T) {

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDE4Mzg2MzcsInJsLW1heC1yZXF1ZXN0cyI6Miwicmwtc2Vjb25kcy1ibG9ja2VkIjo2MCwidXNlciI6ImFkbWluIn0.OEpFjMJURF9jmZASE1n3TOEIFCqh_56xeVYLeigFO_w"

	e := _jwt.VerifyToken(expiredToken)

	assert.Equal(t, "token is unauthorized", e.Error())
}
