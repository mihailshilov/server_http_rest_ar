package model

import "time"

//User
type User struct {
	ID     int
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
	Status   string `json:"status"`
	Response string `json:"response"`
}

//response struct booking
type ResponseBooking struct {
	StatusMs       string `json:"status_ms"`
	ResponseMs     string `json:"response_ms"`
	StatusGazCrm   string `json:"status_gcrm"`
	ResponseGazCrm string `json:"response_gcrm"`
}
