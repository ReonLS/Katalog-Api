package utils

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// enforce security ()
type ContextKey string

const (
	ClaimsKey ContextKey = "claims"
)

// custom claims, designed by including what information's necessary to authorize the user
type JWTclaim struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

var jwtkey []byte

func Init() {
	jwtkey = []byte(os.Getenv("JWT_KEY"))
}

func GenerateJWT(id string, email string, role string) (string, error) {
	//Alur: Initialize custom struct, lalu generate JWT Token, lalu sign dengan secret key
	claims := &JWTclaim{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtkey)

	return signedToken, err
}

func ParseToken(signedToken string) (*JWTclaim, error) {
	//Alur: Parsing signedtoken jd header,payload, signature, dan populate param2 dgn payload
	//Mereturn secret key dengan callback di param 3, supaya bisa validasi signingmenthod
	//populate token.valid apabila no issue, including cek expiresAt
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTclaim{},
		func(token *jwt.Token) (any, error) {
			//bandingin token method dengan familynya (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Different Signing Method")
			}
			return jwtkey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Token Not Valid")
	}

	//typecast claims jadi *jwtclaims supaya bisa akses custom claims
	claims, ok := token.Claims.(*JWTclaim)
	if !ok {
		return nil, errors.New("Could not parse claims")
	}

	//berarti aman
	return claims, nil
}

func GetClaimsFromContext(ctx context.Context) (*JWTclaim, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*JWTclaim)
	return claims, ok
}
