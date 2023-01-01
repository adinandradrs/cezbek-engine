package model

import "database/sql"

type Parameter struct {
	Id         sql.NullInt64  `json:"id"`
	ParamGroup sql.NullString `json:"param_group"`
	ParamName  sql.NullString `json:"param_name"`
	ParamValue sql.NullString `json:"value_param"`
	BaseEntity
}
