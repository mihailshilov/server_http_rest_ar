package model

import "time"

//User
type User struct {
	ID     int    `json:"-"`
	Login  string `json:"login"`
	Secret string `json:"secret"`
}

//for jwt verify
type User2 struct {
	UserID uint64
}

//for token and exp
type Token_exp struct {
	Token string
	Exp   time.Time
}

type AccessDetails struct {
	UserId uint64
	Exp    uint64
}

//response struct
type Response struct {
	Status   string `json:"status" example:"OK"`
	Response string `json:"response" example:"data_received"`
}

//response struct booking
type ResponseBooking struct {
	StatusMs       string `json:"status_ms"`
	ResponseMs     string `json:"response_ms"`
	StatusGazCrm   string `json:"status_gcrm"`
	ResponseGazCrm string `json:"response_gcrm"`
}

// type UserRights struct {
// 	UserRightsArr UserRightsArr
// }

type UserRights struct {
	URA []UserRightsArr `json:"rights"`
}

type UserRightsArr struct {
	IdOrg int `json:"id_org"`
	IdDep int `json:"id_dep"`
}
