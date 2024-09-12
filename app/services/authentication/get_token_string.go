package authentication

import (
	"GIG/app/constants/error_messages"
	"errors"
	"log"
	"strings"

	"github.com/revel/revel"
)

func getTokenString(header *revel.RevelHeader, headerName string) (tokenString string, err error) {
	authHeader := header.Get(headerName)
	log.Println(authHeader)
	if authHeader == "" {
		return "", errors.New(error_messages.AuthHeaderNotFound)
	}

	tokenSlice := strings.Split(authHeader, " ")
	if len(tokenSlice) != 2 {
		return "", errors.New(error_messages.InvalidTokenFormat)
	}
	tokenString = tokenSlice[1]
	log.Println(tokenSlice)
	return tokenString, nil

}
