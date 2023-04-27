package sqlstore

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	_ "reflect"
	"strconv"
	"strings"
	"time"

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

	query := `insert into requests ("ИдЗаявки", "ДатаВремяЗаявки", "ДатаВремяЗаписи", "Ответственный", "ИдОрганизации", "ИдПодразделения", "ДатаВремяОбновления") values($1, $2, $3, $4, $5, $6, $7)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	dt := time.Now()

	_, err = tx.Exec(ctx, query,
		data.DataRequest.ИдЗаявки,
		data.DataRequest.ДатаВремяЗаявки,
		data.DataRequest.ДатаВремяЗаписи,
		data.DataRequest.Ответственный,
		data.DataRequest.ИдОрганизации,
		data.DataRequest.ИдПодразделения,
		dt.Format("2006-01-02T15:04:05"),
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

//informs
func (r *DataRepository) QueryInsertInforms(data model.Informs) error {

	query := `insert into informs ("ТипДокумента", "ИдДокумента", "ИдОрганизации", "ИдПодразделения", "ДатаВремяОтправки", "ДатаВремяДоставки", "ДатаВремяОбновления") values($1, $2, $3, $4, $5, $6, $7)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	dt := time.Now()

	_, err = tx.Exec(ctx, query,
		data.DataInform.ТипДокумента,
		data.DataInform.ИдДокумента,
		data.DataInform.ИдОрганизации,
		data.DataInform.ИдПодразделения,
		data.DataInform.ДатаВремяОтправки,
		data.DataInform.ДатаВремяДоставки,
		dt.Format("2006-01-02T15:04:05"),
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

	query := `insert into orders ("ИдЗаказНаряда", "ИдЗаявки", "ИдСводногоЗаказНаряда", "ДатаВремяСоздания", "ДатаВремяОткрытия", "ВидОбращения", "ПовторныйРемонт", "ПричинаОбращения", "VINбазовый", "VINТекущий", "Ответственный", "ИдОрганизации", "ИдПодразделения", "ГосНомерТС", "ПробегТС", "ДатаВремяОбновления") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	reason := strings.Replace(data.DataOrder.ПричинаОбращения, "\n", ", ", -1)

	dt := time.Now()

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

//cons_orders
func (r *DataRepository) QueryInsertConsOrders(data model.ConsOrders) error {

	query := `insert into cons_orders ("ИдСводногоЗаказНаряда", "ИдЗаявки", "ДатаВремяСоздания", "Ответственный", "ИдОрганизации", "ИдПодразделения", "ДатаВремяОбновления") values($1, $2, $3, $4, $5, $6, $7)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	dt := time.Now()

	_, err = tx.Exec(ctx, query,
		data.DataConsOrder.ИдСводногоЗаказНаряда,
		data.DataConsOrder.ИдЗаявки,
		data.DataConsOrder.ДатаВремяСоздания,
		data.DataConsOrder.Ответственный,
		data.DataConsOrder.ИдОрганизации,
		data.DataConsOrder.ИдПодразделения,
		dt.Format("2006-01-02T15:04:05"),
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

		iter = append(
			iter,
			data.DataStatus.ИдЗаказНаряда,
			data.DataStatus.ИдОрганизации,
			data.DataStatus.ИдПодразделения,
			k.Статус,
			k.ДатаВремя,
			dt.Format("2006-01-02T15:04:05"),
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

		// str := "500"
		number, err := strconv.Atoi(k.КоличествоОпераций)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
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
			number,
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
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"works"}, tableWorks, pgx.CopyFromRows(works))
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
