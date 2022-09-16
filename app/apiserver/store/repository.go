package store

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"
)

//user repository
type UserRepository interface {
	//auth methods
	FindUser(string, string) (*model.User, error)
	FindUserid(uint64) error
	//jwt methods
	CreateToken(uint64, *model.Service) (string, time.Time, error)
	ExtractTokenMetadata(*http.Request, *model.Service) (*model.AccessDetails, error)
	VerifyToken(*http.Request, *model.Service) (*jwt.Token, error)
	ExtractToken(*http.Request) string
}

//data repository
type DataRepository interface {
	QueryInsertOrders(model.Orders) error
	QueryInsertConsOrders(model.ConsOrders) error
	QueryInsertRequests(model.Requests) error
	QueryInsertStatuses(model.Statuses) error
}
