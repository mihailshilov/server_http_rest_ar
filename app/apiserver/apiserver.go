package apiserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/store/sqlstore"

	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

func Start(config *model.Service) error {

	dbPostgres, err := newDbPostgres(config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	defer dbPostgres.Close()

	// dbMssql, err := newDbMssql(config.Spec.DBms.Url)
	// if err != nil {
	// 	logger.ErrorLogger.Println(err)
	// 	return err
	// }

	// defer dbMssql.Close()

	// store_db := sqlstore.New(dbPostgres, dbMssql)
	store_db := sqlstore.New(dbPostgres)

	//cert, key files
	fcert, err := filepath.Abs("/root/cert/carsrv.shilov.pro.crt")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	fkey, err := filepath.Abs("/root/cert/carsrv.shilov.pro.key")
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//cert, key load
	cer, err := tls.LoadX509KeyPair(fcert, fkey)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	configCert := &tls.Config{Certificates: []tls.Certificate{cer}}

	caCert, err := ioutil.ReadFile(fcert)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	//setup HTTPS client
	clt := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cer},
			},
		},
	}

	server := newServer(store_db, config, clt)

	//setup HTTPS server
	srv := &http.Server{
		Addr:      config.Spec.Ports.Addr,
		TLSConfig: configCert,
		Handler:   server.router,
	}

	return srv.ListenAndServeTLS(fcert, fkey)
	//return srv.ListenAndServe()
}

//connect to postgres
func newDbPostgres(conf *model.Service) (*pgxpool.Pool, error) {

	config, _ := pgx.ParseConfig("")
	config.Host = conf.Spec.DBpg.Host
	config.Port = conf.Spec.DBpg.Port
	config.User = conf.Spec.DBpg.User
	config.Password = conf.Spec.DBpg.Password
	config.Database = conf.Spec.DBpg.Database
	config.LogLevel = pgx.LogLevelDebug
	config.Logger = logrusadapter.NewLogger(logger.PgLog())
	config.TLSConfig = nil

	poolConfig, _ := pgxpool.ParseConfig("")
	poolConfig.ConnConfig = config
	poolConfig.MaxConnLifetime = time.Duration(conf.Spec.DBpg.MaxConnLifetime) * time.Minute
	poolConfig.MaxConnIdleTime = time.Duration(conf.Spec.DBpg.MaxConnIdletime) * time.Minute
	poolConfig.MaxConns = conf.Spec.DBpg.MaxConns
	poolConfig.MinConns = conf.Spec.DBpg.MinConns
	poolConfig.HealthCheckPeriod = time.Duration(conf.Spec.DBpg.HealthCheckPeriod) * time.Minute

	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return conn, nil
}

//connect to mssql
// func newDbMssql(databaseUrl string) (*sql.DB, error) {

// 	db, err := sql.Open("sqlserver", databaseUrl)
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return nil, err

// 	}

// 	if err := db.Ping(); err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return nil, err
// 	}

// 	return db, nil

// }

//12
