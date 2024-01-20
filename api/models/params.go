package models

type GetAllParams struct {
	Limit      int64  `json:"limit" binding:"required" default:"10"`
	Page       int64  `json:"page" binding:"required" default:"1"`
	Search     string `json:"search"`
	UserID     int64  `json:"user_id"`
	SortByDate string `json:"sort" enums:"desc,asc" default:"desc"`
}
