package model

//предварительная запись
type Requests struct {
	DataRequest DataRequest `json:"requests"`
}

type DataRequest struct {
	ИдЗаявки                string `json:"id_request" validate:"required"`    //*
	ДатаВремяЗаявки         string `json:"date_time_req" validate:"required"` //*
	ДатаВремяИнформирования string `json:"date_time_inf,omitempty"` //*
	Ответственный           string `json:"responsible" validate:"required"`
	ИдОрганизации           string `json:"id_org" validate:"required"`
	ИдПодразделения         string `json:"id_dep" validate:"required"`
}

//orders
type Orders struct {
	DataOrder DataOrder `json:"orders"`
}

type DataOrder struct {
	ИдЗаказНаряда           string `json:"id_order" validate:"required"`
	ИдЗаявки                string `json:"id_request,omitempty"`
	ИдСводногоЗаказНаряда   string `json:"id_cons_order,omitempty"`
	ДатаВремяСоздания       string `json:"date_time_create" validate:"required"`
	ДатаВремяОткрытия       string `json:"date_time_open" validate:"required"`
	ВидОбращения            string `json:"order_type" validate:"required"`
	ПовторныйРемонт         string `json:"re_repair" validate:"required"`
	ПричинаОбращения        string `json:"reason" validate:"required"` //reason
	VINбазовый              string `json:"vin0" validate:"required"`   //*
	VINпослеДоработки       string `json:"vin1,omitempty"`
	ДатаВремяИнформирования string `json:"date_time_form,omitempty"`
	Ответственный           string `json:"responsible" validate:"required"`
	ИдОрганизации           string `json:"id_org" validate:"required"`
	ИдПодразделения         string `json:"id_dep" validate:"required"`
	ГосНомерТС              string `json:"g_num,omitempty"`
	ПробегТС                string `json:"mileage,omitempty"`
}

//cons_orders
type ConsOrders struct {
	DataConsOrder DataConsOrder `json:"cons_orders"`
}

type DataConsOrder struct {
	ИдСводногоЗаказНаряда string `json:"id_cons_order" validate:"required"`    //*
	ИдЗаявки              string `json:"id_request,omitempty"`       //*
	ДатаВремяСоздания     string `json:"date_time_create" validate:"required"` //*
	Ответственный         string `json:"responsible" validate:"required"`
	ИдОрганизации         string `json:"id_org" validate:"required"`
	ИдПодразделения       string `json:"id_dep" validate:"required"`
}

//статусы
type Statuses struct {
	DataStatus DataStatus `json:"order_statuses"`
}

type DataStatus struct {
	ИдЗаказНаряда   string `json:"id_order"` //*
	ИдОрганизации   string `json:"id_org"`
	ИдПодразделения string `json:"id_dep"`
	OrderStatuses   `json:"statuses"`
}

type OrderStatuses []struct {
	Статус    string `json:"status"`
	ДатаВремя string `json:"date_time_status"` //*
}

//запчасти
type Parts struct {
	DataPart DataPart `json:"order_parts"`
}

type DataPart struct {
	ИдЗаказНаряда   string `json:"id_order"` //*
	ИдОрганизации   string `json:"id_org"`
	ИдПодразделения string `json:"id_dep"`
	OrderParts   	`json:"parts"`
}

type OrderParts []struct {
	Наименование		string `json:"name"`
	КаталожныйНомер		string `json:"code_catalog"`
	ЧертежныйНомер		string `json:"code_drawing"`
	Количество			string `json:"number"`
	ЕдИзм				string `json:"units"`
	Стоимость			string `json:"price"`
	НДС					string `json:"vat"`
}

//работы
type Works struct {
	DataWork DataWork `json:"order_works"`
}

type DataWork struct {
	ИдЗаказНаряда   string `json:"id_order"` //*
	ИдОрганизации   string `json:"id_org"`
	ИдПодразделения string `json:"id_dep"`
	OrderWorks  	`json:"works"`
}

type OrderWorks []struct {
	Наименование			   	string `json:"name"`
	КодОперации    				string `json:"code"`  /// не об
	НормативнаяТрудоёмкость    	string `json:"complexity"`
	СтоимостьНЧ    				string `json:"price_hour"`
}

//заказ наряд создание
// type DataOrder []struct {
// 	ИдЗаказНаряд                 string `json:"id_order" validate:"required"`                            //*
// 	ВремяФомрированияЗаказНаряда string `json:"date_time_order" validate:"required,yyyy-mm-ddThh:mm:ss"` //*
// 	ВидОбращения                 string `json:"type" validate:"required"`                                //*
// 	ПовторныйРемонт              bool   `json:"repeated" validate:"required"`
// 	ПричинаОбращения             string `json:"reason" validate:"required"` //*
// 	ЗаявкаИлиРасширение          string `json:"sign" validate:"required"`   //*
// 	Works                        `json:"works" validate:"required,dive,required"`
// 	Parts                        `json:"parts" validate:"required,dive,required"`
// }

// type Works []struct {
// 	НаименованиеРабот       string `json:"work_name" validate:"required"`  //*
// 	НормативнаяТрудоёмкость int    `json:"normativ" validate:"required"`   //*
// 	СтоимостьНормоЧаса      int    `json:"price_hour" validate:"required"` //*
// }

// type Parts []struct {
// 	НаименованияЗапаснойЧасти string  `json:"repairs" validate:"required"`       //*
// 	КаталожныйНомер           string  `json:"number" validate:"required"`        //*
// 	Количество                int     `json:"qua" validate:"required"`           //*
// 	ЕдИзмерения               string  `json:"unit" validate:"required"`          //*
// 	Стоимость                 float32 `json:"price_repairs" validate:"required"` //*
// 	Поставщик                 string  `json:"provider" validate:"required"`      //*
// }
