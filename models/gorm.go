package models

func (OrmGoodsLIst) TableName() string {
	return "goods_list"
}

type OrmGoodsLIst struct {
	Gid  int
	Name string
	Url  string
	Type string
}

type OrmGoodInfo struct {
	Gid       int
	Sales     string
	Commit    int
	Grate     int
	Introduce string
	Price     float64
}

func (OrmGoodInfo) TableName() string {
	return "goods_info"
}

func (OrmUserInfo) TableName() string {
	return "user_info"
}

type OrmUserInfo struct {
	Uid     int     `gorm:"uid"`
	Hub_id  int     `gorm:"hub_Id"`
	Name    string  `gorm:"name"`
	Word    string  `gorm:"word"`
	Number  string  `gorm:"number"`
	Balance float64 `gorm:"balance"`
	Image   string  `gorm:"image"`
}

func (OrmOrderState) TableName() string {
	return "order_state_type"
}

type OrmOrderState struct {
	Stateid int    `gorm:"Stateid"`
	State   string `gorm:"State"`
}

func (OrmUserOrder) TableName() string {
	return "user_order"
}

type OrmUserOrder struct {
	Oid   int
	Uid   int
	State string
	Gid   int
	Count int
}

func (OrmShopChart) TableName() string {
	return "user_order"
}

type OrmShopChart struct {
	Id    int `gorm:"chart_id"`
	Uid   int
	Gid   int
	Count int
}

func (OrmCommit) TableName() string {
	return "goods_commit"
}

type OrmCommit struct {
	Cid    int
	Gid    int
	Oid    int
	Url    string
	Commit string
}

func (OrmClass) TableName() string {
	return "goods_class"
}

type OrmClass struct {
	Type  string
	Goods []OrmOneGood `gorm:"-"`
}

type OrmOneGood struct {
	OrmGoodInfo
	OrmGoodsLIst
}
