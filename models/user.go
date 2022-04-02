package models

import (
	"mime/multipart"
)

func (IdToken) TableName() string {
	return "user_info"
}

type IdToken struct {
	User     string `gorm:"name"`
	Id       string `gorm:"-"`
	Phone    string `gorm:"number"`
	Balance  string `gorm:"balance"`
	ImageUrl string `gorm:"image"`
}

type UserInfo struct {
	Uid     int
	HubId   int `gorm:"column:hub_id"`
	Name    string
	Word    string
	Number  string
	Balance float64
	Image   string
}

// BasicInfo 用户基准信息 一般用于 token解析数据绑定
type BasicInfo struct {
	Uid      int    `gorm:"column:uid"`
	Username string `gorm:"column:username"`
}

// Balance 用户余额充值
type Balance struct {
	BasicInfo
	Balance float64 `json:"balance,omitempty" form:"balance"`
}

// User token 附带用户信息
type User struct {
	BasicInfo
}

// UserOrder 用户订单返回
type UserOrder struct {
	BasicInfo
	Allorder []OrmUserOrder
}

// AllOrder 单条订单

// Commit 用户提交评论
type Commit struct {
	BasicInfo
	Oid    int    `json:"oid,omitempty" form:"oid"`
	Commit string `json:"commit,omitempty" form:"commit"`
}

// UserImage 用户头像上传
type UserImage struct {
	BasicInfo
	Image *multipart.FileHeader `json:"image,omitempty" form:"image"`
}

// Info 用于用户个人信息 订单分类展示

type MyInfo struct {
	BasicInfo
	Balance  float64 `gorm:"column:balance"`
	ImageUrl string  `gorm:"column:image"`
	Category []Category
}
type AllOrder struct {
	Gid   int    `gorm:"gid"`
	Oid   int    `gorm:"oid"`
	Count int    `gorm:"count"`
	State string `gorm:"state"`
}

type Category struct {
	State string `gorm:"column:State"`
	Order []AllOrder
}
