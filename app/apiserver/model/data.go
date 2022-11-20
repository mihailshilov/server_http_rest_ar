package model

//предварительная запись
type Requests struct {
	DataRequest DataRequest `json:"requests"`
}

type DataRequest struct {
	ИдЗаявки        string `json:"id_request" validate:"required"`
	ДатаВремяЗаявки string `json:"date_time_req" validate:"required,yyyy-mm-ddThh:mm:ss"`
	ДатаВремяЗаписи string `json:"date_time_reс,omitempty"`
	Ответственный   string `json:"responsible" validate:"required"`
	ИдОрганизации   string `json:"id_org" validate:"required,number"`
	ИдПодразделения string `json:"id_dep" validate:"required,number"`
}

//Информирование
type Informs struct {
	DataInform DataInform `json:"informs"`
}

type DataInform struct {
	ТипДокумента      string `json:"type_doc" validate:"required"`
	ИдДокумента       string `json:"id_doc" validate:"required"`
	ИдОрганизации     string `json:"id_org" validate:"required,number"`
	ИдПодразделения   string `json:"id_dep" validate:"required,number"`
	ДатаВремяОтправки string `json:"date_time_create" validate:"required,yyyy-mm-ddThh:mm:ss"`
	ДатаВремяДоставки string `json:"date_time_delivery"`
}

//orders
type Orders struct {
	DataOrder DataOrder `json:"orders"`
}

type DataOrder struct {
	ИдЗаказНаряда         string `json:"id_order" validate:"required"`
	ИдЗаявки              string `json:"id_request"`
	ИдСводногоЗаказНаряда string `json:"id_cons_order"`
	ДатаВремяСоздания     string `json:"date_time_create" validate:"required,yyyy-mm-ddThh:mm:ss"`
	ДатаВремяОткрытия     string `json:"date_time_open" validate:"required,yyyy-mm-ddThh:mm:ss"`
	ВидОбращения          string `json:"order_type" validate:"required"`
	ПовторныйРемонт       string `json:"re_repair" validate:"required"`
	ПричинаОбращения      string `json:"reason" validate:"required"` //reason
	VINбазовый            string `json:"vin0" validate:"required"`
	VINпослеДоработки     string `json:"vin1"`
	Ответственный         string `json:"responsible" validate:"required"`
	ИдОрганизации         string `json:"id_org" validate:"required,number"`
	ИдПодразделения       string `json:"id_dep" validate:"required,number"`
	ГосНомерТС            string `json:"g_num"`
	ПробегТС              string `json:"mileage"`
}

//cons_orders
type ConsOrders struct {
	DataConsOrder DataConsOrder `json:"cons_orders"`
}

type DataConsOrder struct {
	ИдСводногоЗаказНаряда string `json:"id_cons_order" validate:"required"`
	ИдЗаявки              string `json:"id_request"`
	ДатаВремяСоздания     string `json:"date_time_create" validate:"required,yyyy-mm-ddThh:mm:ss"`
	Ответственный         string `json:"responsible" validate:"required"`
	ИдОрганизации         string `json:"id_org" validate:"required,number"`
	ИдПодразделения       string `json:"id_dep" validate:"required,number"`
}

//статусы
type Statuses struct {
	DataStatus DataStatus `json:"order_statuses"`
}

type DataStatus struct {
	ИдЗаказНаряда   string `json:"id_order" validate:"required"`
	ИдОрганизации   string `json:"id_org" validate:"required,number"`
	ИдПодразделения string `json:"id_dep" validate:"required,number"`
	OrderStatuses   `json:"statuses" validate:"required,min=1,dive,required"`
}

type OrderStatuses []struct {
	Статус    string `json:"status" validate:"required"`
	ДатаВремя string `json:"date_time_status" validate:"required,yyyy-mm-ddThh:mm:ss"`
}

//запчасти
type Parts struct {
	DataPart DataPart `json:"order_parts"`
}

type DataPart struct {
	ИдЗаказНаряда   string `json:"id_order" validate:"required"`
	ИдОрганизации   string `json:"id_org" validate:"required,number"`
	ИдПодразделения string `json:"id_dep" validate:"required,number"`
	OrderParts      `json:"parts" validate:"required,min=1,dive,required"`
}

type OrderParts []struct {
	Наименование    string `json:"name" validate:"required"`
	КаталожныйНомер string `json:"code_catalog" validate:"required"`
	ЧертежныйНомер  string `json:"code_drawing" validate:"required"`
	Количество      string `json:"number" validate:"required"`
	ЕдИзм           string `json:"units" validate:"required"`
	Стоимость       string `json:"price" validate:"required"`
	НДС             string `json:"vat" validate:"required"`
}

//работы
type Works struct {
	DataWork DataWork `json:"order_works"`
}

type DataWork struct {
	ИдЗаказНаряда   string `json:"id_order" validate:"required"`
	ИдОрганизации   string `json:"id_org" validate:"required,number"`
	ИдПодразделения string `json:"id_dep" validate:"required,number"`
	OrderWorks      `json:"works" validate:"required,min=1,dive,required"`
}

type OrderWorks []struct {
	Наименование            string `json:"name" validate:"required"`
	КодОперации             string `json:"code"`
	НормативнаяТрудоёмкость string `json:"complexity" validate:"required"`
	СтоимостьНЧ             string `json:"price_hour" validate:"required"`
}

//Машины для сайта
type CarsForSite struct {
	DataCarForSite DataCarForSite `json:"cars_for_site"`
}

type DataCarForSite struct {
	Id_org string `json:"id_org" validate:"required,number"`
	Cars   `json:"cars"`
}

type Cars []struct {
	Vin    string `json:"vin"`
	Id_isk string `json:"id_isk"`
	Flag   string `json:"flag"`
}

//заказ наряд создание
// type DataOrder []struct {
// 	ИдЗаказНаряд                 string `json:"id_order" validate:"required"`
// 	ВремяФомрированияЗаказНаряда string `json:"date_time_order" validate:"required,yyyy-mm-ddThh:mm:ss"`
// 	ВидОбращения                 string `json:"type" validate:"required"`
// 	ПовторныйРемонт              bool   `json:"repeated" validate:"required"`
// 	ПричинаОбращения             string `json:"reason" validate:"required"`
// 	ЗаявкаИлиРасширение          string `json:"sign" validate:"required"`
// 	Works                        `json:"works" validate:"required,dive,required"`
// 	Parts                        `json:"parts" validate:"required,dive,required"`
// }

// type Works []struct {
// 	НаименованиеРабот       string `json:"work_name" validate:"required"`
// 	НормативнаяТрудоёмкость int    `json:"normativ" validate:"required"`
// 	СтоимостьНормоЧаса      int    `json:"price_hour" validate:"required"`
// }

// type Parts []struct {
// 	НаименованияЗапаснойЧасти string  `json:"repairs" validate:"required"`
// 	КаталожныйНомер           string  `json:"number" validate:"required"`
// 	Количество                int     `json:"qua" validate:"required"`
// 	ЕдИзмерения               string  `json:"unit" validate:"required"`
// 	Стоимость                 float32 `json:"price_repairs" validate:"required"`
// 	Поставщик                 string  `json:"provider" validate:"required"`
// }
