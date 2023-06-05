package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	_ "reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"

	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

//Data repository
type DataRepository struct {
	store *Store
}

//requests
func (r *DataRepository) QueryInsertRequests(data model.Requests) error {

	query := `insert into requests ("ИдЗаявки", "ДатаВремяЗаявки", "ДатаВремяЗаписи", "Ответственный", "ИдОрганизации", "ИдПодразделения", "ДатаВремяОбновления", "Uid_request") values($1, $2, $3, $4, $5, $6, $7, $8)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	dt := time.Now()

	var uid_req uuid.UUID
	if data.DataRequest.Uid_request != "" {
		uid_req, err = uuid.FromString(data.DataRequest.Uid_request)
		if err != nil {
			logger.ErrorLogger.Println("Неверный формат Guid req")
			return err
		}
	}

	_, err = tx.Exec(ctx, query,
		data.DataRequest.ИдЗаявки,
		data.DataRequest.ДатаВремяЗаявки,
		data.DataRequest.ДатаВремяЗаписи,
		data.DataRequest.Ответственный,
		data.DataRequest.ИдОрганизации,
		data.DataRequest.ИдПодразделения,
		dt.Format("2006-01-02T15:04:05"),
		uid_req,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

func (r *DataRepository) IsRequestUnic(data model.Requests) error {
	query := `select * from requests where "Uid_request" = $1 order  by id desc limit 1`

	//db table
	type request struct {
		Id           int
		IdRequest    string
		DateTimeReq  string
		DateTimeRec  string
		Rresponsible string
		IdOrg        string
		IdDep        string
		DateTimeUp   string
		UidRequest   uuid.UUID
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	var rows pgx.Rows

	var LastRecRow request

	rows, err = tx.Query(ctx, query, data.DataRequest.Uid_request)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&LastRecRow.Id, &LastRecRow.IdRequest, &LastRecRow.DateTimeReq, &LastRecRow.DateTimeRec, &LastRecRow.Rresponsible, &LastRecRow.IdOrg, &LastRecRow.IdDep, &LastRecRow.DateTimeUp, &LastRecRow.UidRequest)
		if err != nil {
			return err
		}

	}

	lastreqstring := LastRecRow.IdRequest + LastRecRow.DateTimeReq + LastRecRow.DateTimeRec + LastRecRow.Rresponsible + LastRecRow.IdOrg + LastRecRow.IdDep

	//logger.InfoLogger.Println(lastreqstring)

	newrecstring := data.DataRequest.ИдЗаявки + data.DataRequest.ДатаВремяЗаявки + data.DataRequest.ДатаВремяЗаписи + data.DataRequest.Ответственный + data.DataRequest.ИдОрганизации + data.DataRequest.ИдПодразделения

	//logger.InfoLogger.Println(newrecstring)

	if lastreqstring == newrecstring {

		err := errors.New("заявка повторилась")

		//logger.InfoLogger.Println(err)

		return err
	}

	return nil

}

//informs
func (r *DataRepository) QueryInsertInforms(data model.Informs) error {

	query := `insert into informs ("ТипДокумента", "ИдДокумента", "ИдОрганизации", "ИдПодразделения", "ДатаВремяОтправки", "ДатаВремяДоставки", "ДатаВремяОбновления", "Uid_doc") values($1, $2, $3, $4, $5, $6, $7, $8)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	dt := time.Now()

	var uid_doc uuid.UUID
	//var err error
	if data.DataInform.Uid_doc != "" {
		uid_doc, err = uuid.FromString(data.DataInform.Uid_doc)
		if err != nil {
			logger.ErrorLogger.Println("Неверный формат Guid документа")

			return err
		}
	}

	_, err = tx.Exec(ctx, query,
		data.DataInform.ТипДокумента,
		data.DataInform.ИдДокумента,
		data.DataInform.ИдОрганизации,
		data.DataInform.ИдПодразделения,
		data.DataInform.ДатаВремяОтправки,
		data.DataInform.ДатаВремяДоставки,
		dt.Format("2006-01-02T15:04:05"),
		uid_doc,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

//orders
func (r *DataRepository) QueryInsertOrders(data model.Orders) error {

	query := `insert into orders ("ИдЗаказНаряда", "ИдЗаявки", "ИдСводногоЗаказНаряда", "ДатаВремяСоздания", "ДатаВремяОткрытия", "ВидОбращения", "ПовторныйРемонт", "ПричинаОбращения", "VINбазовый", "VINТекущий", "Ответственный", "ИдОрганизации", "ИдПодразделения", "ГосНомерТС", "ПробегТС", "ДатаВремяОбновления", "Uid_order", "Uid_request", "Uid_consorder") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	reason := strings.Replace(data.DataOrder.ПричинаОбращения, "\n", ". ", -1)

	dt := time.Now()

	var uid_ord, uid_req, uid_cons uuid.UUID
	if data.DataOrder.Uid_order != "" {
		uid_ord, err = uuid.FromString(data.DataOrder.Uid_order)
		if err != nil {
			logger.ErrorLogger.Println("Неверный формат Guid ord")
			return err
		}
	}

	if data.DataOrder.Uid_request != "" {
		uid_req, err = uuid.FromString(data.DataOrder.Uid_request)
		if err != nil {
			logger.ErrorLogger.Println("Неверный формат Guid req")
			return err
		}
	}

	if data.DataOrder.Uid_consorder != "" {
		uid_cons, err = uuid.FromString(data.DataOrder.Uid_consorder)
		if err != nil {
			logger.ErrorLogger.Println("Неверный формат Guid cons")
			return err
		}
	}

	_, err = tx.Exec(ctx, query,
		data.DataOrder.ИдЗаказНаряда,
		data.DataOrder.ИдЗаявки,
		data.DataOrder.ИдСводногоЗаказНаряда,
		data.DataOrder.ДатаВремяСоздания,
		data.DataOrder.ДатаВремяОткрытия,
		data.DataOrder.ВидОбращения,
		data.DataOrder.ПовторныйРемонт,
		reason,
		data.DataOrder.VINбазовый,
		data.DataOrder.VINТекущий,
		data.DataOrder.Ответственный,
		data.DataOrder.ИдОрганизации,
		data.DataOrder.ИдПодразделения,
		data.DataOrder.ГосНомерТС,
		data.DataOrder.ПробегТС,
		dt.Format("2006-01-02T15:04:05"),
		uid_ord,
		uid_req,
		uid_cons,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

func (r *DataRepository) IsOrderUnic(data model.Orders) error {

	query := `select * from orders where "Uid_order" = $1 order  by id desc limit 1`

	type order struct {
		Id           int
		IdOrder      string
		IdOrg        string
		IdDep        string
		IdConsOrder  string
		IdRequest    string
		DateTimeRec  string
		DateTimeOpen string
		OrderType    string
		ReRepair     string
		Reason       string
		Vin0         string
		Vin1         string
		Rresponsible string
		GNum         string
		Mileage      string
		DateTimeUp   string
		UidOrder     uuid.UUID
		UidRequest   uuid.UUID
		UidConsOrder uuid.UUID
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	var rows pgx.Rows

	var LastRecRow order

	rows, err = tx.Query(ctx, query, data.DataOrder.Uid_order)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&LastRecRow.Id, &LastRecRow.IdOrder, &LastRecRow.IdOrg, &LastRecRow.IdDep, &LastRecRow.IdConsOrder, &LastRecRow.IdRequest, &LastRecRow.DateTimeRec, &LastRecRow.DateTimeOpen, &LastRecRow.OrderType, &LastRecRow.ReRepair, &LastRecRow.Reason, &LastRecRow.Vin0, &LastRecRow.Vin1, &LastRecRow.Rresponsible, &LastRecRow.GNum, &LastRecRow.Mileage, &LastRecRow.DateTimeUp, &LastRecRow.UidOrder, &LastRecRow.UidRequest, &LastRecRow.UidConsOrder)
		if err != nil {
			logger.InfoLogger.Println(err)
			return err
		}

	}

	lastorederstring := LastRecRow.IdOrder + LastRecRow.IdOrg + LastRecRow.IdDep + LastRecRow.IdConsOrder + LastRecRow.IdRequest + LastRecRow.DateTimeRec + LastRecRow.DateTimeOpen + LastRecRow.OrderType + LastRecRow.ReRepair + LastRecRow.Reason + LastRecRow.Vin0 + LastRecRow.Vin1 + LastRecRow.Rresponsible + LastRecRow.GNum + LastRecRow.Mileage

	neworderstring := data.DataOrder.ИдЗаказНаряда + data.DataOrder.ИдОрганизации + data.DataOrder.ИдПодразделения + data.DataOrder.ИдСводногоЗаказНаряда + data.DataOrder.ИдЗаявки + data.DataOrder.ДатаВремяСоздания + data.DataOrder.ДатаВремяОткрытия + data.DataOrder.ВидОбращения + data.DataOrder.ПовторныйРемонт + data.DataOrder.ПричинаОбращения + data.DataOrder.VINбазовый + data.DataOrder.VINТекущий + data.DataOrder.Ответственный + data.DataOrder.ГосНомерТС + data.DataOrder.ПробегТС

	logger.InfoLogger.Println(lastorederstring)
	logger.InfoLogger.Println(neworderstring)

	if lastorederstring == neworderstring {
		err := errors.New("З-Н повторился")

		//logger.InfoLogger.Println(err)

		return err
	}

	return nil

}

//cons_orders
func (r *DataRepository) QueryInsertConsOrders(data model.ConsOrders) error {

	query := `insert into cons_orders ("ИдСводногоЗаказНаряда", "ИдЗаявки", "ДатаВремяСоздания", "Ответственный", "ИдОрганизации", "ИдПодразделения", "ДатаВремяОбновления", "Uid_consorder", "Uid_request") values($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	logger.InfoLogger.Println("Запущен блок транзакций") //поиск бага

	dt := time.Now()

	var uid_req, uid_cons uuid.UUID
	if data.DataConsOrder.Uid_request != "" {
		uid_req, err = uuid.FromString(data.DataConsOrder.Uid_request)
		if err != nil {
			logger.ErrorLogger.Println("Неверный формат Guid req")
			return err
		}
	}

	if data.DataConsOrder.Uid_consorder != "" {
		uid_cons, err = uuid.FromString(data.DataConsOrder.Uid_consorder)
		if err != nil {
			logger.ErrorLogger.Println("Неверный формат Guid cons")
			return err
		}
	}

	_, err = tx.Exec(ctx, query,
		data.DataConsOrder.ИдСводногоЗаказНаряда,
		data.DataConsOrder.ИдЗаявки,
		data.DataConsOrder.ДатаВремяСоздания,
		data.DataConsOrder.Ответственный,
		data.DataConsOrder.ИдОрганизации,
		data.DataConsOrder.ИдПодразделения,
		dt.Format("2006-01-02T15:04:05"),
		uid_cons,
		uid_req,
	)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	logger.InfoLogger.Println("Запущен Exec") //поиск бага

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

func (r *DataRepository) IsConsOrderUnic(data model.ConsOrders) error {

	query := `select * from cons_orders where "Uid_consorder" = $1 order  by id desc limit 1`

	type consorder struct {
		Id           int
		IdConsOrder  string
		IdOrg        string
		IdDep        string
		IdRequest    string
		DateTimeRec  string
		Rresponsible string
		DateTimeUp   string
		UidConsOrder uuid.UUID
		UidRequest   uuid.UUID
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	logger.InfoLogger.Println("Запущен блок транзакций") //поиск бага

	var rows pgx.Rows

	var LastRecRow consorder

	rows, err = tx.Query(ctx, query, data.DataConsOrder.Uid_consorder)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	defer rows.Close()

	logger.InfoLogger.Println("tx.Query") //поиск бага

	for rows.Next() {

		err := rows.Scan(&LastRecRow.Id, &LastRecRow.IdConsOrder, &LastRecRow.IdOrg, &LastRecRow.IdDep, &LastRecRow.IdRequest, &LastRecRow.DateTimeRec, &LastRecRow.Rresponsible, &LastRecRow.DateTimeUp, &LastRecRow.UidConsOrder, &LastRecRow.UidRequest)
		if err != nil {
			logger.InfoLogger.Println(err)
			return err
		}

	}

	lastconsorederstring := LastRecRow.IdConsOrder + LastRecRow.IdOrg + LastRecRow.IdDep + LastRecRow.IdRequest + LastRecRow.DateTimeRec + LastRecRow.Rresponsible

	newconsorrderstring := data.DataConsOrder.ИдСводногоЗаказНаряда + data.DataConsOrder.ИдОрганизации + data.DataConsOrder.ИдПодразделения + data.DataConsOrder.ИдЗаявки + data.DataConsOrder.ДатаВремяСоздания + data.DataConsOrder.Ответственный

	if lastconsorederstring == newconsorrderstring {
		err := errors.New("сводный з-н повторился")

		//logger.InfoLogger.Println(err)

		logger.InfoLogger.Println("Дубли найдены") //поиск бага

		return err
	}

	logger.InfoLogger.Println("Дубли не найдены") //поиск бага

	return nil

}

//statuses old
// func (r *DataRepository) QueryInsertStatuses_old(data model.Statuses) error {

// 	//query := `insert into statuses values($1, $2, $3, $4, $5)`

// 	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancelFunc()

// 	tx, err := r.store.dbPostgres.Begin(context.Background())
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	_, err = tx.Exec(ctx, query,
// 		data.DataStatus.ИдЗаказНаряд,
// 		data.DataStatus.ИдОрганизации,
// 		data.DataStatus.ИдПодразделения,
// 		data.DataStatus.Статус,
// 		data.DataStatus.ДатаВремя,
// 	)
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	return nil

// }

func (r *DataRepository) QueryInsertStatuses(data model.Statuses) error {

	var statuses [][]interface{}

	statusesValSlice := reflect.ValueOf(data.DataStatus).FieldByName("OrderStatuses").Interface().(model.OrderStatuses)

	//fmt.Println(reflect.ValueOf(data))

	for _, k := range statusesValSlice {

		var iter []interface{}

		dt := time.Now()

		var uid_ord uuid.UUID
		var err error
		if data.DataStatus.Uid_order != "" {
			uid_ord, err = uuid.FromString(data.DataStatus.Uid_order)
			if err != nil {
				logger.ErrorLogger.Println("Неверный формат Guid ord")

				return err
			}
		}

		iter = append(
			iter,
			data.DataStatus.ИдЗаказНаряда,
			data.DataStatus.ИдОрганизации,
			data.DataStatus.ИдПодразделения,
			k.Статус,
			k.ДатаВремя,
			dt.Format("2006-01-02T15:04:05"),
			uid_ord,
		)

		statuses = append(statuses, iter)

	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	//ctx.Value("test")
	//context.WithValue(ctx, "1", "123")

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//Statuses
	tableStatuses := []string{
		"ИдЗаказНаряда",
		"ИдОрганизации",
		"ИдПодразделения",
		"Статус",
		"ДатаВремя",
		"ДатаВремяОбновления",
		"Uid_order",
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"statuses"}, tableStatuses, pgx.CopyFromRows(statuses))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

func (r *DataRepository) IsOrderReal(idOrder string) error {
	query := `select count(*) from orders where ИдЗаказНаряда like $1`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	type Row struct {
		count int
	}

	var rows pgx.Rows

	rows, err = tx.Query(ctx, query, idOrder)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	defer rows.Close()

	var rowSlice []Row
	for rows.Next() {
		var r Row
		err := rows.Scan(&r.count)
		if err != nil {
			return err
		}
		rowSlice = append(rowSlice, r)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if rowSlice[0].count == 0 {
		err := errors.New("нет заказ-наряда")
		logger.ErrorLogger.Println("нет заказ-наряда")
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func (r *DataRepository) IsRequestReal(idRequest string) error {
	query := `select count(*) from requests where ИдЗаявки like $1`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	type Row struct {
		count int
	}

	var rows pgx.Rows

	rows, err = tx.Query(ctx, query, idRequest)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	defer rows.Close()

	var rowSlice []Row
	for rows.Next() {
		var r Row
		err := rows.Scan(&r.count)
		if err != nil {
			return err
		}
		rowSlice = append(rowSlice, r)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if rowSlice[0].count == 0 {
		err := errors.New("нет заявки")
		logger.ErrorLogger.Println("нет заявки")
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func (r *DataRepository) QueryInsertParts(data model.Parts) error {

	var parts [][]interface{}

	partsValSlice := reflect.ValueOf(data.DataPart).FieldByName("OrderParts").Interface().(model.OrderParts)

	//fmt.Println(reflect.ValueOf(data))

	for _, k := range partsValSlice {

		var iter []interface{}

		dt := time.Now()

		var uid_ord uuid.UUID
		var err error
		if data.DataPart.Uid_order != "" {
			uid_ord, err = uuid.FromString(data.DataPart.Uid_order)
			if err != nil {
				logger.ErrorLogger.Println("Неверный формат Guid ord")

				return err
			}
		}

		iter = append(
			iter,
			data.DataPart.ИдЗаказНаряда,
			data.DataPart.ИдОрганизации,
			data.DataPart.ИдПодразделения,
			k.Наименование,
			k.КаталожныйНомер,
			k.ЧертежныйНомер,
			k.Количество,
			k.ЕдИзм,
			k.Стоимость,
			k.НДС,
			dt.Format("2006-01-02T15:04:05"),
			uid_ord,
			k.Скидка,
		)

		parts = append(parts, iter)

	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//Parts
	tableParts := []string{
		"ИдЗаказНаряда",
		"ИдОрганизации",
		"ИдПодразделения",
		"Наименование",
		"КаталожныйНомер",
		"ЧертежныйНомер",
		"Количество",
		"ЕдИзм",
		"Стоимость",
		"НДС",
		"ДатаВремяОбновления",
		"Uid_order",
		"discount",
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"parts"}, tableParts, pgx.CopyFromRows(parts))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

func (r *DataRepository) QueryInsertWorks(data model.Works) error {

	var works [][]interface{}

	worksValSlice := reflect.ValueOf(data.DataWork).FieldByName("OrderWorks").Interface().(model.OrderWorks)

	//fmt.Println(reflect.ValueOf(data))

	for _, k := range worksValSlice {

		var iter []interface{}

		dt := time.Now()

		var uid_ord uuid.UUID
		var err error
		if data.DataWork.Uid_order != "" {
			uid_ord, err = uuid.FromString(data.DataWork.Uid_order)
			if err != nil {
				logger.ErrorLogger.Println("Неверный формат Guid ord")

				return err
			}
		}

		iter = append(
			iter,
			data.DataWork.ИдЗаказНаряда,
			data.DataWork.ИдОрганизации,
			data.DataWork.ИдПодразделения,
			k.Наименование,
			k.КодОперации,
			k.НормативнаяТрудоёмкость,
			k.СтоимостьНЧ,
			dt.Format("2006-01-02T15:04:05"),
			k.КоличествоОпераций,
			uid_ord,
			k.НДС,
			k.Скидка,
			k.ЕдИзм,
		)

		works = append(works, iter)

	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//Works
	tableWorks := []string{
		"ИдЗаказНаряда",
		"ИдОрганизации",
		"ИдПодразделения",
		"Наименование",
		"КодОперации",
		"НормативнаяТрудоёмкость",
		"СтоимостьНЧ",
		"ДатаВремяОбновления",
		"КоличествоОпераций",
		"Uid_order",
		"НДС",
		"discount",
		"ЕдИзм",
	}

	xx, err := tx.CopyFrom(ctx, pgx.Identifier{"works"}, tableWorks, pgx.CopyFromRows(works))
	if err != nil {
		logger.ErrorLogger.Println(err)
		logger.InfoLogger.Println(xx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

func (r *DataRepository) QueryInsertCarsForSite(data model.CarsForSite) error {

	var carsforsite [][]interface{}

	carsforsiteValSlice := reflect.ValueOf(data.DataCarForSite).FieldByName("Cars").Interface().(model.Cars)

	//fmt.Println(reflect.ValueOf(data))

	for _, k := range carsforsiteValSlice {

		var iter []interface{}

		dt := time.Now()

		iter = append(
			iter,
			data.DataCarForSite.Id_org,
			k.Vin,
			k.Id_isk,
			k.Flag,
			dt.Format("2006-01-02T15:04:05"),
			k.Vin_current,
		)

		carsforsite = append(carsforsite, iter)

	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	//ctx.Value("test")
	//context.WithValue(ctx, "1", "123")

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//CarsForSite
	tableCarsForSite := []string{
		"id_org",
		"vin",
		"id_isk",
		"flag",
		"date_rec",
		"vin_current",
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"carsforsite"}, tableCarsForSite, pgx.CopyFromRows(carsforsite))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

//query insert mssql

func (r *DataRepository) QueryInsertMssql(data model.CarsForSite) ([]model.ISKStatus, error) {

	//var carsforsite [][]interface{}

	carsforsiteValSlice := reflect.ValueOf(data.DataCarForSite).FieldByName("Cars").Interface().(model.Cars)

	//var response *model.ResponseCarsForSite

	var mssql_respond string
	var mssql_mess string
	var mssql_errors []string
	mssql_responds := []model.ISKStatus{}

	for _, k := range carsforsiteValSlice {

		//var iter []interface{}
		iter := &model.ISKStatus{}

		//request mssql

		_, err := r.store.dbMssql.Exec(r.store.config.Spec.Queryies.Booking,
			sql.Named("VIN", k.Vin),
			sql.Named("НомернойТовар", k.Id_isk),
			sql.Named("Значение", k.Flag),
			sql.Named("Результат", sql.Out{Dest: &mssql_respond}),
			sql.Named("Сообщение", sql.Out{Dest: &mssql_mess}),
		)
		if err != nil {
			//return "", err
			logger.ErrorLogger.Println(err)
			err_mes := "ошибка в" + k.Id_isk
			mssql_errors = append(mssql_errors, err_mes)
			logger.ErrorLogger.Println(mssql_respond)
			return nil, err
		}

		iter = &model.ISKStatus{
			Vin:    k.Vin,
			Id_isk: k.Id_isk,
			Flag:   k.Flag,
			MsResp: mssql_respond,
			MsMess: mssql_mess,
		}

		mssql_responds = append(mssql_responds, *iter)

		logger.InfoLogger.Println(mssql_respond)
		logger.InfoLogger.Println(mssql_mess)

		//response = append(response, )

		//return mssql_respond, nil
	}

	return mssql_responds, nil

}

func (r *DataRepository) QueryUpdateCarsForSite(data []model.ISKStatus) error {
	return nil
}

//request Azgaz catalog
/*
func (r *DataRepository) RequestAzgaz(data []model.DataAzgaz, config *model.Service) (*model.ResponseAzgaz, error) {

	for _, car := range data {
		if car.MsResp == "0" {
			logger.InfoLogger.Println(car.Vin)
		}
	}
*/

//var dataset model.DataAzgaz
//var response *model.ResponseAzgaz

// bodyJson := &model.DataAzgazReq{
// 	Visible: data.Visible,
// }

// dataset.Data = append(dataset.Data, bodyJson)

// bodyBytesReq, err := json.Marshal(dataset)
// if err != nil {
// 	return nil, err
// }

// resp, err := http.Put(config.Spec.Client.UrlAzgazTest+data.Vin, "application/json", bytes.NewBuffer(bodyBytesReq))
// if err != nil {
// 	logger.ErrorLogger.Println(err)
// 	return nil, err
// }

// defer resp.Body.Close()

// bodyBytesResp, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	logger.ErrorLogger.Println(err)
// 	return nil, err
// }

// if err := json.Unmarshal(bodyBytesResp, &response); err != nil {
// 	logger.ErrorLogger.Println(err)
// 	return nil, err
// }

//return nil, nil

//}

//Logistic
// func (r *DataRepository) QueryInsertMssql(jsonLogistic json) error {

// }

/*archive*/

// func (r *DataRepository) QueryInsertOrders(data model.Orders) error {

// 	var orders [][]interface{}
// 	ordersValSlice := reflect.ValueOf(data.ListOrders).FieldByName("DataOrder").Interface().(model.DataOrder)

// 	for _, k := range ordersValSlice {
// 		var iter []interface{}
// 		iter = append(
// 			iter,
// 			data.ListOrders.ИдЗаявки,
// 			data.ListOrders.VINбазовый,
// 			data.ListOrders.VINТекущий,
// 			k.ИдЗаказНаряд,
// 			k.ВремяФомрированияЗаказНаряда,
// 			k.ВидОбращения,
// 			k.ПовторныйРемонт,
// 			k.ПричинаОбращения,
// 			k.ЗаявкаИлиРасширение,
// 		)
// 		orders = append(orders, iter)
// 	}

// 	//parts
// 	var parts [][]interface{}
// 	for _, k := range ordersValSlice {
// 		for _, l := range k.Parts {
// 			var iter []interface{}
// 			iter = append(
// 				iter,
// 				data.ListOrders.ИдЗаявки,
// 				data.ListOrders.VINбазовый,
// 				data.ListOrders.VINТекущий,
// 				k.ИдЗаказНаряд,
// 				l.НаименованияЗапаснойЧасти,
// 				l.КаталожныйНомер,
// 				l.Количество,
// 				l.ЕдИзмерения,
// 				l.Стоимость,
// 				l.Поставщик,
// 			)
// 			parts = append(parts, iter)
// 		}
// 	}

// 	//works
// 	var works [][]interface{}
// 	for _, k := range ordersValSlice {
// 		for _, l := range k.Works {
// 			var iter []interface{}
// 			iter = append(
// 				iter,
// 				data.ListOrders.ИдЗаявки,
// 				data.ListOrders.VINбазовый,
// 				data.ListOrders.VINТекущий,
// 				k.ИдЗаказНаряд,
// 				l.НаименованиеРабот,
// 				l.НормативнаяТрудоёмкость,
// 				l.СтоимостьНормоЧаса,
// 			)
// 			works = append(works, iter)
// 		}
// 	}

// 	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancelFunc()

// 	tx, err := r.store.dbPostgres.Begin(context.Background())
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	//orders
// 	tableOrders := []string{
// 		"ИдЗаявки",
// 		"vinбазовый",
// 		"VINТекущий",
// 		"ИдЗаказНаряд",
// 		"ВремяФомрированияЗаказНаряда",
// 		"ВидОбращения",
// 		"ПовторныйРемонт",
// 		"ПричинаОбращения",
// 		"ЗаявкаИлиРасширение",
// 	}

// 	_, err = tx.CopyFrom(ctx, pgx.Identifier{"orders"}, tableOrders, pgx.CopyFromRows(orders))
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	//parts
// 	tableParts := []string{
// 		"ИдЗаявки",
// 		"vinбазовый",
// 		"VINТекущий",
// 		"ИдЗаказНаряд",
// 		"НаименованияЗапаснойЧасти",
// 		"КаталожныйНомер",
// 		"Количество",
// 		"ЕдИзмерения",
// 		"Стоимость",
// 		"Поставщик",
// 	}

// 	_, err = tx.CopyFrom(ctx, pgx.Identifier{"parts"}, tableParts, pgx.CopyFromRows(parts))
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	//works
// 	tableWorks := []string{
// 		"ИдЗаявки",
// 		"vinбазовый",
// 		"VINТекущий",
// 		"ИдЗаказНаряд",
// 		"НаименованиеРабот",
// 		"НормативнаяТрудоёмкость",
// 		"СтоимостьНормоЧаса",
// 	}

// 	_, err = tx.CopyFrom(ctx, pgx.Identifier{"works"}, tableWorks, pgx.CopyFromRows(works))
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		logger.ErrorLogger.Println(err)
// 		return err
// 	}

// 	return nil

// }
