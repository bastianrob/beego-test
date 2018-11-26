package models

//Product db model
type Product struct {
	ID       int64  `orm:"pk;auto;column(id)" json:"id"`
	CreateTs int64  `orm:"column(create_ts)" json:"create_ts"`
	Name     string `orm:"column(name)" json:"name"`
	Category string `orm:"column(category)" json:"category"`
	Price    int64  `orm:"column(price)" json:"price"`
}
