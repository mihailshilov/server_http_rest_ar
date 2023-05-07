package store

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	GetUserRights(uint64) ([]model.UserRightsArr, error)
}

//data repository
type DataRepository interface {
	QueryInsertOrders(model.Orders) error
	QueryInsertConsOrders(model.ConsOrders) error
	QueryInsertRequests(model.Requests) error
	QueryInsertStatuses(model.Statuses) error
	QueryInsertParts(model.Parts) error
	QueryInsertWorks(model.Works) error
	QueryInsertInforms(model.Informs) error
	QueryInsertCarsForSite(model.CarsForSite) error
	//QueryUpdateCarsForSite([]model.ISKStatus) error
	QueryInsertMssql(model.CarsForSite) ([]model.ISKStatus, error)
	IsOrderReal(idOrder string) error
	IsRequestReal(idOrder string) error
	//RequestAzgaz(data []model.DataAzgaz, config *model.Service) (*model.ResponseAzgaz, error)
	//QueryInsertLogistic(jsonLogistic) error
	IsRequestUnic(model.Requests) error
	IsOrderUnic(model.Orders) error
	IsConsOrderUnic(model.ConsOrders) error
}
