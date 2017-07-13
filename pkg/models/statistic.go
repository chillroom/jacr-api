package models

import "github.com/jmoiron/sqlx/types"

type Statistic struct {
	Name  string         `db:"name"`
	Value types.JSONText `db:"value"`
}
