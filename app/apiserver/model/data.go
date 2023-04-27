package model

//предварительная запись
type Requests struct {
	DataRequest DataRequest `json:"requests"`
}

type DataRequest struct {
	ИдЗаявки        string `json:"id_request" validate:"required"`
	Uid_request     string `json:"uid_request"`
	ДатаВремяЗаявки string `json:"date_time_req" validate:"required,yyyy-mm-ddThh:mm:ss"`
	ДатаВремяЗаписи string `json:"date_time_rec,omitempty"`
	Ответственный   string `json:"responsible" validate:"required"`
	ИдОрганизации   string `json:"id_org" validate:"required,number"`
	ИдПодразделения string `json:"id_dep" validate:"required,number"`
}

//Информирование
type Informs struct {
	DataInform DataInform `json:"informs"`
}

type DataInform struct {
	ТипДокумента      string `json:"type_doc" validate:"required,oneof=Заявка Заказ-наряд"`
	ИдДокумента       string `json:"id_doc" validate:"required"`
	Uid_doc           string `json:"uid_doc"`
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
	Uid_order             string `json:"uid_order"`
	Uid_request           string `json:"uid_request"`
	Uid_consorder         string `json:"uid_cons_order"`
	ДатаВремяСоздания     string `json:"date_time_create" validate:"required,yyyy-mm-ddThh:mm:ss"`
	ДатаВремяОткрытия     string `json:"date_time_open" validate:"required,yyyy-mm-ddThh:mm:ss"`
	ВидОбращения          string `json:"order_type" validate:"required"`
	ПовторныйРемонт       string `json:"re_repair" validate:"required,oneof=Да Нет"`
	ПричинаОбращения      string `json:"reason"` //reason
	VINбазовый            string `json:"vin0"`
	VINТекущий            string `json:"vin1"`
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
	Uid_consorder         string `json:"uid_cons_order"`
	Uid_request           string `json:"uid_request"`
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
	Uid_order       string `json:"uid_order"`
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
	Uid_order       string `json:"uid_order"`
	ИдОрганизации   string `json:"id_org" validate:"required,number"`
	ИдПодразделения string `json:"id_dep" validate:"required,number"`
	OrderParts      `json:"parts" validate:"required,min=1,dive,required"`
}

type OrderParts []struct {
	Наименование    string `json:"name" validate:"required"`
	КаталожныйНомер string `json:"code_catalog" validate:"required"`
	ЧертежныйНомер  string `json:"code_drawing" validate:"required"`
	Количество      string `json:"number" validate:"required,numeric"`
	ЕдИзм           string `json:"units" validate:"required"`
	Стоимость       string `json:"price" validate:"required,numeric"`
	НДС             string `json:"vat" validate:"required,numeric"`
	Скидка          string `json:"discount" validate:""`
}

//работы
type Works struct {
	DataWork DataWork `json:"order_works"`
}

type DataWork struct {
	ИдЗаказНаряда   string `json:"id_order" validate:"required"`
	Uid_order       string `json:"uid_order"`
	ИдОрганизации   string `json:"id_org" validate:"required,number"`
	ИдПодразделения string `json:"id_dep" validate:"required,number"`
	OrderWorks      `json:"works" validate:"required,min=1,dive,required"`
}

type OrderWorks []struct {
	Наименование            string `json:"name" validate:"required"`
	КодОперации             string `json:"code"`
	НормативнаяТрудоёмкость string `json:"complexity" validate:"required,numeric"`
	КоличествоОпераций      string `json:"number" validate:"required,numeric"`
	ЕдИзм                   string `json:"units"`
	СтоимостьНЧ             string `json:"price_hour" validate:"required,numeric"`
	НДС                     string `json:"vat" validate:""`
	Скидка                  string `json:"discount" validate:""`
}

//Машины для сайта
type CarsForSite struct {
	DataCarForSite DataCarForSite `json:"cars_for_site"`
}

type DataCarForSite struct {
	Id_org string `json:"id_org" validate:"required,number"`
	Cars   `json:"cars" validate:"required,min=1,dive,required"`
}

type Cars []struct {
	Vin         string `json:"vin" validate:"required"`
	Vin_current string `json:"vin_current"`
	Id_isk      string `json:"id_isk" validate:"required,number"`
	Flag        string `json:"flag" validate:"required,oneof=0 1"`
}

//Data booking
type DataBooking struct {
	Vin      string
	Id_isk   string
	Значение string
}

//resp struct api gaz crm
type ResponseCarsForSite struct {
	StatusMs   StatusMs   `json:"StatusISK"`
	StatusSite StatusSite `json:"StatusSite"`
}

type StatusMs []struct {
	Vin    string `json:"vin"`
	Status string `json:"status"`
}

type StatusSite []struct {
	Vin    string `json:"vin"`
	Status string `json:"status"`
}

// type ISKStatus struct {
// 	ISKStatuses
// }

type ISKStatus struct {
	Vin    string
	Id_isk string
	Flag   string
	MsResp string
	MsMess string
}

//resp struct api gaz crm
type ResponseAzgaz struct {
	Visible bool `json:"visible"`
}
