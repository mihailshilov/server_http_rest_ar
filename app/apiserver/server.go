package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/gorilla/docs"

	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/store"

	"github.com/go-playground/validator"

	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

//errors
var (
	errIncorrectEmailOrPassword = errors.New("incorrect auth")
	errReg                      = errors.New("service registration error")
	errJwt                      = errors.New("token error")
	errFindUser                 = errors.New("user not found")
	//errMssql                    = errors.New("mssql error")
)

//server configure
type server struct {
	router   *mux.Router
	validate *validator.Validate
	store    store.Store
	config   *model.Service
	client   *http.Client
}

func newServer(store store.Store, config *model.Service, client *http.Client) *server {
	s := &server{
		router:   mux.NewRouter(),
		validate: validator.New(),
		store:    store,
		config:   config,
		client:   client,
	}
	s.configureRouter()
	return s
}

//custome validate date format
func IsDateCorrect(fl validator.FieldLevel) bool {
	DateRegexString := "^(19|20)\\d\\d-(0[1-9]|1[012])-([012]\\d|3[01])T([01]\\d|2[0-3]):([0-5]\\d):([0-5]\\d)$"
	DateRegex := regexp.MustCompile(DateRegexString)
	return DateRegex.MatchString(fl.Field().String())
}

//write new token struct
func newToken(token string, exp time.Time) *model.Token_exp {
	return &model.Token_exp{
		Token: token,
		Exp:   exp,
	}
}

//write response struct
func newResponse(status string, response string) *model.Response {
	return &model.Response{
		Status:   status,
		Response: response,
	}
}

//write http error
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

//write http response
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	//open

	//s.router.HandleFunc("/swagger", s.handleSwagger()).Methods("GET")
	s.router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	s.router.HandleFunc("/authentication", s.handleAuth()).Methods("POST")

	//private
	auth := s.router.PathPrefix("/auth").Subrouter()
	auth.Use(s.middleWare)
	//booking, forms submit
	auth.HandleFunc("/servicerequests", s.handleRequests()).Methods("POST")
	auth.HandleFunc("/serviceconsorders", s.handleConsOrders()).Methods("POST")
	auth.HandleFunc("/serviceorders", s.handleOrders()).Methods("POST")
	auth.HandleFunc("/servicestatuses", s.handleStatuses()).Methods("POST")
	auth.HandleFunc("/serviceparts", s.handleParts()).Methods("POST")
	auth.HandleFunc("/serviceworks", s.handleWorks()).Methods("POST")

}

//handle Auth
func (s *server) handleAuth() http.HandlerFunc {

	var req model.User

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, errReg)
			logger.ErrorLogger.Println(err)
			return
		}

		u, err := s.store.User().FindUser(req.Login, req.Secret)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			logger.ErrorLogger.Println(err)
			return
		}

		token, datetime_exp, err := s.store.User().CreateToken(uint64(u.ID), s.config)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errJwt)
			logger.ErrorLogger.Println(err)
			return
		}
		token_data := newToken(token, datetime_exp)
		s.respond(w, r, http.StatusOK, token_data)
		logger.InfoLogger.Println("token issued success")

	}

}

//Middleware
func (s *server) middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//extract user_id
		user_id, err := s.store.User().ExtractTokenMetadata(r, s.config)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errJwt)
			logger.ErrorLogger.Println(err)
			return
		}

		if err := s.store.User().FindUserid(user_id.UserId); err != nil {
			s.error(w, r, http.StatusUnauthorized, errFindUser)
			logger.ErrorLogger.Println(err)
			return
		}

		next.ServeHTTP(w, r)

	})

}

//handle service requests
func (s *server) handleRequests() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Requests{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//fmt.Println(req) //debug

		logger.InfoLogger.Println("good request )")
		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Data().QueryInsertRequests(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

//handle service consolidated orders
func (s *server) handleConsOrders() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.ConsOrders{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("good request )")
		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Data().QueryInsertConsOrders(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

//handle service orders
func (s *server) handleOrders() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Orders{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		fmt.Println(r) //debug

		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
		}

		if err := s.store.Data().QueryInsertOrders(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))
	}

}

//handle service statuses
func (s *server) handleStatuses() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Statuses{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}
		fmt.Println(req) //debug

		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Data().QueryInsertStatuses(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("good request )")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

//handle service parts
func (s *server) handleParts() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Parts{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		fmt.Println('-') //debug
		fmt.Println(req) //debug
		fmt.Println('-') //debug

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Data().QueryInsertParts(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("good request )")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

//handle service works
func (s *server) handleWorks() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Works{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}
		//fmt.Println(req) //debug

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Data().QueryInsertWorks(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("good request )")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}
