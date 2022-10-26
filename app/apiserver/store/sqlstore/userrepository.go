package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/store"

	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

//User repository
type UserRepository struct {
	store *Store
}

type AccessDetails struct {
	UserId uint64
	Exp    uint64
}

//Find jwt email password (create token)
func (r *UserRepository) FindUser(login string, secret string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.dbPostgres.QueryRow(context.Background(),
		"SELECT id, login, secret FROM users WHERE login = $1 AND secret = $2",
		login, secret).Scan(&u.ID, &u.Login, &u.Secret); err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Println(err)
			return nil, store.ErrRecordNotFound
		}
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return u, nil
}

//Find jwt user id (verify token)
func (r *UserRepository) FindUserid(userid uint64) error {
	u := &model.User2{}

	if err := r.store.dbPostgres.QueryRow(context.Background(),
		"SELECT id FROM users WHERE id = $1",
		userid).Scan(&u.UserID); err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Println(err)
			return store.ErrRecordNotFound
		}
		logger.ErrorLogger.Println(err)

		return err
	}
	return nil
}

//creating token
//create token
func (r *UserRepository) CreateToken(userid uint64, config *model.Service) (string, time.Time, error) {
	var err error

	datetime_exp_unix := time.Now().Add(time.Hour * 24 * time.Duration(config.Spec.Jwt.LifeTerm)).Unix()
	datetime_exp := time.Unix(datetime_exp_unix, 0)
	t := new(time.Time)

	//os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = datetime_exp_unix
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(config.Spec.Jwt.TokenDecode))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return "", *t, err
	}

	return token, datetime_exp, nil
}

//extract token from header
func (r *UserRepository) ExtractToken(req *http.Request) string {
	bearToken := req.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return bearToken
}

//verify token
func (r *UserRepository) VerifyToken(req *http.Request, config *model.Service) (*jwt.Token, error) {
	tokenString := r.ExtractToken(req)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Spec.Jwt.TokenDecode), nil
	})
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}
	return token, nil
}

//extract data from token
func (r *UserRepository) ExtractTokenMetadata(req *http.Request, config *model.Service) (*model.AccessDetails, error) {

	//var accessDetails model.AccessDetails

	token, err := r.VerifyToken(req, config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		return &model.AccessDetails{
			UserId: userId,
		}, nil
	}
	logger.ErrorLogger.Println(err)
	return nil, err
}
