package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtHandler struct{
	SecretKey []byte
	DefaultMaxDuration int
}

func pickFirstOne(secretKeys ...string) []byte {
	for _, key := range(secretKeys){
		if key != ""{
			return []byte(key)
		}
	}
	return []byte{}
}

func NewJwtHandler(secretkey string) *JwtHandler{
	maxduration, err := strconv.Atoi(os.Getenv("DEFAULT_JWT_EXPIRY_TIME"))
	if err != nil{
		return nil
	}
	return &JwtHandler{
		SecretKey: pickFirstOne(secretkey, os.Getenv("JWT")),
		DefaultMaxDuration: maxduration,
	}
}

type UserClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func (jt *JwtHandler) Encode(details map[string]string, durationInSeconds int) (string, error) {
	if durationInSeconds <= 0 && durationInSeconds >= jt.DefaultMaxDuration {
		return "", errors.New("duration is invalid")
	}

	uid, _ := strconv.ParseUint(details["id"], 10, 64)

	expirationTime := time.Now().Add(time.Duration(durationInSeconds) * time.Second)
	uc := UserClaims{
		UserID: uint(uid),
		Email:  details["email"],
		Name:   details["name"],
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	return token.SignedString(jt.SecretKey)
}

func (jt *JwtHandler) Decode(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jt.SecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}