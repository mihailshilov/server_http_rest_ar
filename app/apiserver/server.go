package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/mihailshilov/server_http_rest_ar/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/store"

	"github.com/go-playground/validator/v10"

	stats_api "github.com/fukata/golang-stats-api-handler"
	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

type ctxKey string

const keyUserRights ctxKey = "user_rights"

// @title API для сервисных станций СТТ
// @version 1.0
// @oas 3
// @description API-сервер для сбора данных о работе сервисных станций стт
// @contact.name API Support
// @contact.email shilovmo@st.tech
// @host https://carsrv.st.tech
// @BasePath /
// @query.collection.format multi
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

//errors

var (
	errIncorrectEmailOrPassword = errors.New("incorrect auth")
	errReg                      = errors.New("service registration error")
	errJwt                      = errors.New("token error")
	errFindUser                 = errors.New("user not found")
	errMssql                    = errors.New("mssql error")
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

//custom tunslate for erros
func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Поле является обязательным"
	case "number":
		return "Поле должно быть целым числом"
	case "min":
		return "Нельзя передавать пустой масив"
	case "numeric":
		return "Поле должно быть вещественным числом"
	case "oneof":
		return "Допустимо только значение из списка для этого поля"
	case "yyyy-mm-ddThh:mm:ss":
		return "Время указано не верно"
	}
	return fe.Error() // default error
}

func FindRights(a []model.UserRightsArr, b model.UserRightsArr) bool {
	for _, n := range a {
		if b == n {
			return true
		}
	}
	return false
}

type ApiError struct {
	Param   string
	Message string
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

//write response struct carsforsite
/*
func newResponseBooking(statusms string, responsems string, statussite string, responsesite string) *model.ResponseCarsForSite {
	return &model.ResponseCarsForSite{
		StatusMs:     statusms,
		ResponseMs:   responsems,
		StatusSite:   statussite,
		ResponseSite: responsesite,
	}
}
*/

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
		httpSwagger.URL("doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		//httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
		//httpSwagger.Plugins([]string),
		httpSwagger.UIConfig(map[string]string{
			"showExtensions":        "true",
			"onComplete":            `() => { window.ui.setBasePath('v3'); }`,
			"defaultModelRendering": `"model"`,
		}),
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
	auth.HandleFunc("/serviceinform", s.handleInforms()).Methods("POST")
	auth.HandleFunc("/carsforsite", s.handleCarsForSite()).Methods("POST")
	auth.HandleFunc("/stats", stats_api.Handler).Methods("GET")

}

// HandleAuth godoc
// @Summary Авторизация
// @Description Auth Login
// @Tags Авторизация
// @ID auth-login
// @Accept  json
// @Produce  json
// @Param input body model.User true "user info"
// @Success 200 {object} model.Token_exp "OK"
// @Router /authentication/ [post]
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

		w.Header().Add("Content-Type", "application/json")

		//extract user_id
		user_id, err := s.store.User().ExtractTokenMetadata(r, s.config)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errJwt)
			logger.ErrorLogger.Println(err)
			return
		}

		//add user_id to context

		user_rights, err := s.store.User().GetUserRights(user_id.UserId)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return
		}

		ctx := context.WithValue(r.Context(), keyUserRights, user_rights)

		if err := s.store.User().FindUserid(user_id.UserId); err != nil {
			s.error(w, r, http.StatusUnauthorized, errFindUser)
			logger.ErrorLogger.Println(err)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

// handleRequests godoc
// @Summary Создать заявку
// @Tags Отправка данных
// @Description Создать заявку
// @ID create-request
// @Accept  json
// @Produce  json
// @Param input body model.DataRequest true "request info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/servicerequests/ [post]
// @Security ApiKeyAuth
func (s *server) handleRequests() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Requests{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)

			errs := err.(validator.ValidationErrors)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		//проверка на дубли

		if err := s.store.Data().IsRequestUnic(req); err != nil {
			logger.ErrorLogger.Println("Заявка " + req.DataRequest.ИдЗаявки + " дублируется. Запись не внесена в БД, гуид: " + req.DataRequest.Uid_request)
			s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))
			return
		}
		logger.InfoLogger.Println("Проверка заявки на уникальность выполнена")

		if err := s.store.Data().QueryInsertRequests(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		logger.InfoLogger.Println("good request - Заявка передана")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

// handleInforms godoc
// @Summary Создать информирование
// @Tags Отправка данных
// @Description Создать информировние
// @ID create-inform
// @Accept  json
// @Produce  json
// @Param input body model.DataInform true "inform info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/serviceinform/ [post]
// @Security ApiKeyAuth
func (s *server) handleInforms() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Informs{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		logger.InfoLogger.Println("good request )")

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)

			errs := err.(validator.ValidationErrors)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		//Проверка наличия заказ-наряда
		if req.DataInform.ТипДокумента == "Заказ-наряд" {
			if err := s.store.Data().IsOrderReal(req.DataInform.ИдДокумента); err != nil {
				logger.ErrorLogger.Println(err)
				out_order := make([]ApiError, 1)
				out_order[0] = ApiError{"id_doc", "Заказ-наряд не найден"}

				s.respond(w, r, http.StatusBadRequest, out_order)
				return
			}
		} else {
			if err := s.store.Data().IsRequestReal(req.DataInform.ИдДокумента); err != nil {
				logger.ErrorLogger.Println(err)
				out_order := make([]ApiError, 1)
				out_order[0] = ApiError{"id_doc", "Заявка не найдена"}

				s.respond(w, r, http.StatusBadRequest, out_order)
				return
			}
		}

		if err := s.store.Data().QueryInsertInforms(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

// handleConsOrders godoc
// @Summary Создать сводный заказ-наряд
// @Tags Отправка данных
// @Description Создать сводный заказ-наряд
// @ID create-consorder
// @Accept  json
// @Produce  json
// @Param input body model.DataConsOrder true "consOrder info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/consorders/ [post]
// @Security ApiKeyAuth
func (s *server) handleConsOrders() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.ConsOrders{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)

			errs := err.(validator.ValidationErrors)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		//проверка на дубли

		if err := s.store.Data().IsConsOrderUnic(req); err != nil {
			logger.ErrorLogger.Println("Сводный З-Н " + req.DataConsOrder.ИдСводногоЗаказНаряда + " дублируется. Запись не внесена в БД, гуид: " + req.DataConsOrder.Uid_consorder)
			s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))
			return
		}

		if err := s.store.Data().QueryInsertConsOrders(req); err != nil {
			logger.ErrorLogger.Println(err)
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		logger.InfoLogger.Println("good request - Сводный заказ-наряд добавлен")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

// handleOrders godoc
// @Summary Создать заказ-наряд
// @Tags Отправка данных
// @Description Создать заказ-наряд
// @ID create-order
// @Accept  json
// @Produce  json
// @Param input body model.DataOrder true "order info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/orders/ [post]
// @Security ApiKeyAuth
func (s *server) handleOrders() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Orders{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)
			logger.ErrorLogger.Printf(" Заказ-наряд № " + req.DataOrder.ИдЗаказНаряда)

			errs := err.(validator.ValidationErrors)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		//Проверка прав

		userRights, ok := r.Context().Value(keyUserRights).([]model.UserRightsArr)
		if !ok {
			logger.ErrorLogger.Println("Ид пользователя не найден")
			return
		}

		logger.InfoLogger.Println(userRights)

		org_id, _ := strconv.Atoi(req.DataOrder.ИдОрганизации)
		dep_id, _ := strconv.Atoi(req.DataOrder.ИдПодразделения)

		reqOrgDep := model.UserRightsArr{IdOrg: org_id, IdDep: dep_id}

		logger.InfoLogger.Println(reqOrgDep)

		haveRights := FindRights(userRights, reqOrgDep)

		if haveRights != true {
			logger.ErrorLogger.Println("Недостаточно прав для данной организации")
			s.respond(w, r, http.StatusOK, newResponse("error", "Недостаточно прав для данной организации (некорректные аттрибуты id_org/id_dep)"))
			return
		}

		// if err := s.store.Data().RightsСheck(req, userID); err != nil {
		// 	logger.ErrorLogger.Println("У пользователя № " + strconv.FormatInt(int64(userID), 10) + "Надостаточно прав на Организацию:" + req.DataOrder.ИдОрганизации + " и подразделения№ " + req.DataOrder.ИдПодразделения)
		// 	s.respond(w, r, http.StatusOK, newResponse("error", "Недостаточно прав"))
		// 	return
		// }

		//проверка на дубли

		if err := s.store.Data().IsOrderUnic(req); err != nil {
			logger.ErrorLogger.Println("З-Н " + req.DataOrder.ИдЗаказНаряда + " дублируется. Запись не внесена в БД, гуид: " + req.DataOrder.Uid_order)
			s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))
			return
		}

		if err := s.store.Data().QueryInsertOrders(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("good request - Заказ-наряд добавлен")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))
	}

}

// handleStatuses godoc
// @Summary Добавить статус
// @Tags Отправка данных
// @Description Добавить статус заказ-наряда
// @ID create-status
// @Accept  json
// @Produce  json
// @Param input body model.DataStatus true "status info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/statuses/ [post]
// @Security ApiKeyAuth
func (s *server) handleStatuses() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Statuses{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}
		fmt.Println(req) //debug

		//Проверка наличия заказ-наряда
		if err := s.store.Data().IsOrderReal(req.DataStatus.ИдЗаказНаряда); err != nil {
			logger.ErrorLogger.Println(err)
			out_order := make([]ApiError, 1)
			out_order[0] = ApiError{"id_order", "Заказ-наряд не найден"}

			s.respond(w, r, http.StatusBadRequest, out_order)
			return
		}

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)

			errs := err.(validator.ValidationErrors)

			fmt.Println(errs)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		if err := s.store.Data().QueryInsertStatuses(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("good request - статусы добавлены")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

// handleParts godoc
// @Summary Добавить запчасти
// @Tags Отправка данных
// @Description Добавить запчасти заказ-наряда
// @ID create-parts
// @Accept  json
// @Produce  json
// @Param input body model.DataPart true "parts info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/parts/ [post]
// @Security ApiKeyAuth
func (s *server) handleParts() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Parts{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//Проверка наличия заказ-наряда
		if err := s.store.Data().IsOrderReal(req.DataPart.ИдЗаказНаряда); err != nil {
			logger.ErrorLogger.Println(err)
			out_order := make([]ApiError, 1)
			out_order[0] = ApiError{"id_order", "Заказ-наряд не найден"}

			s.respond(w, r, http.StatusBadRequest, out_order)
			return
		}

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)

			errs := err.(validator.ValidationErrors)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		if err := s.store.Data().QueryInsertParts(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}

		logger.InfoLogger.Println("good request - запчасти добавлены")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

// handleWorks godoc
// @Summary Добавить работы
// @Tags Отправка данных
// @Description Добавить работы заказ-наряда
// @ID create-works
// @Accept  json
// @Produce  json
// @Param input body model.DataWork true "works info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/works/ [post]
// @Security ApiKeyAuth
func (s *server) handleWorks() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.Works{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//Проверка наличия заказ-наряда
		if err := s.store.Data().IsOrderReal(req.DataWork.ИдЗаказНаряда); err != nil {
			logger.ErrorLogger.Println(err)
			out_order := make([]ApiError, 1)
			out_order[0] = ApiError{"id_order", "Заказ-наряд не найден"}

			s.respond(w, r, http.StatusBadRequest, out_order)
			return
		}

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)

			errs := err.(validator.ValidationErrors)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		if err := s.store.Data().QueryInsertWorks(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}

		logger.InfoLogger.Println("good request - работы переданы")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

// handleCarsForSite godoc
// @Summary Данные по проданным автомобилям
// @Tags Отправка данных
// @Description Добавить статус заказ-наряда
// @ID create-carsforsite
// @Accept  json
// @Produce  json
// @Param input body model.CarsForSite true "cars for site info"
// @Success 200 {object} model.Response "OK"
// @Router /auth/carsforsite/ [post]
// @Security ApiKeyAuth
func (s *server) handleCarsForSite() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		req := model.CarsForSite{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			logger.ErrorLogger.Println(err)
			return
		}

		//Валидация
		_ = s.validate.RegisterValidation("yyyy-mm-ddThh:mm:ss", IsDateCorrect)

		s.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		if err := s.validate.Struct(req); err != nil {
			logger.ErrorLogger.Println(err)

			errs := err.(validator.ValidationErrors)

			out := make([]ApiError, len(errs))

			for i, e := range errs {

				out[i] = ApiError{e.Field(), msgForTag(e)}

			}

			s.respond(w, r, http.StatusBadRequest, out)

			return
		}

		if err := s.store.Data().QueryInsertCarsForSite(req); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}

		// QueryInsertCarsForSiteToIsk (построчно)

		/*
					хп_УстановитьВидимостьАвтомобиляДляСайта @VIN = vin, @НомернойТовар = id_isk, @Значение = flag, @Сообщение varchar(900) out, @Результат int out
			Любой @Результат != 0 это ошибка, комментарий в переменной @Сообщение
			@Результат = 0, @Сообщение = 'Успешно'

		*/

		resp_isk, err := s.store.Data().QueryInsertMssql(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errMssql)
			logger.ErrorLogger.Println(err)

			s.respond(w, r, http.StatusBadRequest, err)

			return
		}
		logger.ErrorLogger.Println(resp_isk)

		// вернуть ответы по каждой строке и
		// записать ответы в бд (date_rec_isk & status_rec_isk) 6-12
		////
		/*
			if err := s.store.Data().QueryUpdateCarsForSite(resp_isk); err != nil {
				logger.ErrorLogger.Println(err)
				return
			}
		*/
		// при успешном ответе по каждой строке дернуть азгаз
		// записать ответ в бд

		//request gazcrm api

		//resp_isk //тут будет цикл

		/*

			resAzgaz, err := s.store.Data().RequestAzgaz(resp_isk, s.config)
			if err != nil {
				logger.ErrorLogger.Println(err)
			}


		*/
		// if resAzgaz.Visible {
		// 	logger.ErrorLogger.Println(resAzgaz)
		// 	s.respond(w, r, http.StatusBadRequest, newResponse("Error", resAzgaz.Visible))
		// } else {
		// 	logger.InfoLogger.Println("gazcrm form data transfer success")
		// 	s.respond(w, r, http.StatusOK, newResponse("Ok", resAzgaz.Visible))
		// }

		logger.InfoLogger.Println("good request - автомобили обновлены")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))

	}

}

func (s *server) handleLogistic() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		logger.InfoLogger.Println("good request )")
		s.respond(w, r, http.StatusOK, newResponse("ok", "data_received"))
	}

}
