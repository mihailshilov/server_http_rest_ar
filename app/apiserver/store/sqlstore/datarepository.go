package sqlstore

import (
	"context"
	"reflect"
	_ "reflect"
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

	query := `insert into requests ("ИдЗаявки", "ДатаВремяЗаявки", "ДатаВремяИнформирования", "Ответственный", "ИдОрганизации", "ИдПодразделения") values($1, $2, $3, $4, $5, $6)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.DataRequest.ИдЗаявки,
		data.DataRequest.ДатаВремяЗаявки,
		data.DataRequest.ДатаВремяИнформирования,
		data.DataRequest.Ответственный,
		data.DataRequest.ИдОрганизации,
		data.DataRequest.ИдПодразделения,
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

	query := `insert into informs ("ТипДокумента", "ИдДокумента", "ИдОрганизации", "ИдПодразделения", "ДатаВремяОтправки", "ДатаВремяДоставки") values($1, $2, $3, $4, $5, $6)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.DataInform.ТипДокумента,
		data.DataInform.ИдДокумента,
		data.DataInform.ИдОрганизации,
		data.DataInform.ИдПодразделения,
		data.DataInform.ДатаВремяОтправки,
		data.DataInform.ДатаВремяДоставки,
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

	query := `insert into orders ("ИдЗаказНаряда", "ИдЗаявки", "ИдСводногоЗаказНаряда", "ДатаВремяСоздания", "ДатаВремяОткрытия", "ВидОбращения", "ПовторныйРемонт", "ПричинаОбращения", "VINбазовый", "VINпослеДоработки", "ДатаВремяИнформирования", "Ответственный", "ИдОрганизации", "ИдПодразделения", "ГосНомерТС", "ПробегТС") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.DataOrder.ИдЗаказНаряда,
		data.DataOrder.ИдЗаявки,
		data.DataOrder.ИдСводногоЗаказНаряда,
		data.DataOrder.ДатаВремяСоздания,
		data.DataOrder.ДатаВремяОткрытия,
		data.DataOrder.ВидОбращения,
		data.DataOrder.ПовторныйРемонт,
		data.DataOrder.ПричинаОбращения,
		data.DataOrder.VINбазовый,
		data.DataOrder.VINпослеДоработки,
		data.DataOrder.ДатаВремяИнформирования,
		data.DataOrder.Ответственный,
		data.DataOrder.ИдОрганизации,
		data.DataOrder.ИдПодразделения,
		data.DataOrder.ГосНомерТС,
		data.DataOrder.ПробегТС,
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

	query := `insert into cons_orders ("ИдСводногоЗаказНаряда", "ИдЗаявки", "ДатаВремяСоздания", "Ответственный", "ИдОрганизации", "ИдПодразделения") values($1, $2, $3, $4, $5, $6)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.DataConsOrder.ИдСводногоЗаказНаряда,
		data.DataConsOrder.ИдЗаявки,
		data.DataConsOrder.ДатаВремяСоздания,
		data.DataConsOrder.Ответственный,
		data.DataConsOrder.ИдОрганизации,
		data.DataConsOrder.ИдПодразделения,
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

		iter = append(
			iter,
			data.DataStatus.ИдЗаказНаряда,
			data.DataStatus.ИдОрганизации,
			data.DataStatus.ИдПодразделения,
			k.Статус,
			k.ДатаВремя,
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

func (r *DataRepository) QueryInsertParts(data model.Parts) error {

	var parts [][]interface{}

	partsValSlice := reflect.ValueOf(data.DataPart).FieldByName("OrderParts").Interface().(model.OrderParts)

	//fmt.Println(reflect.ValueOf(data))

	for _, k := range partsValSlice {

		var iter []interface{}

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

		iter = append(
			iter,
			data.DataWork.ИдЗаказНаряда,
			data.DataWork.ИдОрганизации,
			data.DataWork.ИдПодразделения,
			k.Наименование,
			k.КодОперации,
			k.НормативнаяТрудоёмкость,
			k.СтоимостьНЧ,
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

		iter = append(
			iter,
			data.DataCarForSite.Id_org,
			k.Vin,
			k.Id_isk,
			k.Flag,
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

// func (r *DataRepository) QueryInsertOrders(data model.Orders) error {

// 	var orders [][]interface{}
// 	ordersValSlice := reflect.ValueOf(data.ListOrders).FieldByName("DataOrder").Interface().(model.DataOrder)

// 	for _, k := range ordersValSlice {
// 		var iter []interface{}
// 		iter = append(
// 			iter,
// 			data.ListOrders.ИдЗаявки,
// 			data.ListOrders.VINбазовый,
// 			data.ListOrders.VINпослеДоработки,
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
// 				data.ListOrders.VINпослеДоработки,
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
// 				data.ListOrders.VINпослеДоработки,
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
// 		"vinпослеДоработки",
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
// 		"vinпослеДоработки",
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
// 		"vinпослеДоработки",
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
