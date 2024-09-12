package authentication

import (
	"GIG/app/constants/error_messages"
	"GIG/app/constants/info_messages"
	"GIG/app/constants/user_roles"
	"GIG/app/repositories"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/lsflk/gig-sdk/libraries"
	"github.com/lsflk/gig-sdk/models"
	"github.com/pkg/errors"
	"github.com/revel/revel"
)

// Authenticate is and method will be called before any authenticate needed action.
// In order to valid the user.
func Authenticate(c *revel.Controller) revel.Result {

	user, authMethod, err := GetAuthUser(c.Request.Header)

	if err != nil {
		log.Println("Failed to authenticate user")
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(err.Error())
	}

	// if ApiKey exist and not requesting access to AdminControllers
	if authMethod == ApiKey && !libraries.StringInSlice(AdminControllers, c.Name) {
		return nil
	}

	log.Println("Paasing ApiKey validation and failed")

	if err != nil { // if Bearer token doesn't exist
		log.Println(error_messages.TokenApiKeyFailed)

		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(error_messages.TokenApiKeyFailed)
	}

	if user.Role != user_roles.Admin && libraries.StringInSlice(AdminControllers, c.Name) { // Only admin users are allowed to access UserController
		log.Println(err)
		c.Response.Status = http.StatusUnauthorized
		return c.RenderJSON(error_messages.InvalidLoginCredentials)
	}

	log.Println(info_messages.LoginSuccess, user.Email)
	return nil
}

func GetAuthUser(header *revel.RevelHeader) (models.User, string, error) {
	log.Println("Request Headers:")
	log.Println(header.Server.GetKeys())
	for _, key := range header.Server.GetKeys() {
		values := header.Get(key)
		log.Printf("%s: %v", key, values)
	}
	tokenString, err := getTokenString(header, AuthHeaderName)
	apiKey, keyErr := getTokenString(header, ApiKeyHeaderName)

	if keyErr == nil { // if ApiKey exist
		log.Println("ApiKey found")
		user, userErr := repositories.UserRepository{}.GetUserBy(ApiKey, apiKey)
		log.Println(user)
		if userErr == nil {
			return user, ApiKey, nil
		}
	}

	log.Println("Got key error")
	if err != nil { // if Bearer token doesn't exist
		return models.User{}, Bearer, errors.New(error_messages.TokenApiKeyFailed)
	}

	var claims jwt.MapClaims
	claims, err = decodeToken(tokenString)
	if err != nil {
		return models.User{}, Bearer, errors.New(error_messages.TokenDecodeError)
	}
	email, found := claims["iss"]
	if !found {
		log.Println(err)
		return models.User{}, Bearer, err
	}
	user, err := repositories.UserRepository{}.GetUserBy(Email, fmt.Sprintf("%v", email))

	if err != nil {
		return models.User{}, Bearer, errors.New(error_messages.InvalidLoginCredentials)
	}

	return user, Bearer, nil
}
