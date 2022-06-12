package sqlstore

import (
	"context"
	"reflect"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/mihailshilov/server_http_rest_ar/app/apiserver/model"

	logger "github.com/mihailshilov/server_http_rest_ar/app/apiserver/logger"
)

//Data repository
type DataRepository struct {
	store *Store
}

func (r *DataRepository) QueryInsertOrders(data model.Orders) error {

	var orders [][]interface{}
	ordersValSlice := reflect.ValueOf(data.ListOrders).FieldByName("DataOrder").Interface().(model.DataOrder)

	for _, k := range ordersValSlice {
		var iter []interface{}
		iter = append(
			iter,
			data.ListOrders.ИдЗаявки,
			data.ListOrders.VINбазовый,
			data.ListOrders.VINпослеДоработки,
			k.ИдЗаказНаряд,
			k.ВремяФомрированияЗаказНаряда,
			k.ВидОбращения,
			k.ПовторныйРемонт,
			k.ПричинаОбращения,
			k.ЗаявкаИлиРасширение,
		)
		orders = append(orders, iter)
	}

	//parts
	var parts [][]interface{}
	for _, k := range ordersValSlice {
		for _, l := range k.Parts {
			var iter []interface{}
			iter = append(
				iter,
				data.ListOrders.ИдЗаявки,
				data.ListOrders.VINбазовый,
				data.ListOrders.VINпослеДоработки,
				k.ИдЗаказНаряд,
				l.НаименованияЗапаснойЧасти,
				l.КаталожныйНомер,
				l.Количество,
				l.ЕдИзмерения,
				l.Стоимость,
				l.Поставщик,
			)
			parts = append(parts, iter)
		}
	}

	//works
	var works [][]interface{}

	for _, k := range ordersValSlice {
		for _, l := range k.Works {
			var iter []interface{}
			iter = append(
				iter,
				data.ListOrders.ИдЗаявки,
				data.ListOrders.VINбазовый,
				data.ListOrders.VINпослеДоработки,
				k.ИдЗаказНаряд,
				l.НаименованиеРабот,
				l.НормативнаяТрудоёмкость,
				l.СтоимостьНормоЧаса,
			)
			works = append(works, iter)
		}
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//orders
	tableOrders := []string{
		"ИдЗаявки",
		"vinбазовый",
		"vinпослеДоработки",
		"ИдЗаказНаряд",
		"ВремяФомрированияЗаказНаряда",
		"ВидОбращения",
		"ПовторныйРемонт",
		"ПричинаОбращения",
		"ЗаявкаИлиРасширение",
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"orders"}, tableOrders, pgx.CopyFromRows(orders))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//parts
	tableParts := []string{
		"ИдЗаявки",
		"vinбазовый",
		"vinпослеДоработки",
		"ИдЗаказНаряд",
		"НаименованияЗапаснойЧасти",
		"КаталожныйНомер",
		"Количество",
		"ЕдИзмерения",
		"Стоимость",
		"Поставщик",
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"parts"}, tableParts, pgx.CopyFromRows(parts))
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	//works
	tableWorks := []string{
		"ИдЗаявки",
		"vinбазовый",
		"vinпослеДоработки",
		"ИдЗаказНаряд",
		"НаименованиеРабот",
		"НормативнаяТрудоёмкость",
		"СтоимостьНормоЧаса",
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

//requests
func (r *DataRepository) QueryInsertRequests(data model.Requests) error {

	query := `insert into requests values($1, $2, $3, $4)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.DataRequest.ИдЗаявки,
		data.DataRequest.VINбазовый,
		data.DataRequest.VINпослеДоработки,
		data.DataRequest.ВремяЗаявки,
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

//statuses
func (r *DataRepository) QueryInsertStatuses(data model.Statuses) error {

	query := `insert into statuses values($1, $2, $3, $4, $5)`

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()

	tx, err := r.store.dbPostgres.Begin(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	_, err = tx.Exec(ctx, query,
		data.DataStatus.ИдЗаказНаряд,
		data.DataStatus.VINбазовый,
		data.DataStatus.VINпослеДоработки,
		data.DataStatus.ТекущийСтатусЗаказНаряда,
		data.DataStatus.ВремяПрисвоенияСтатуса,
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
