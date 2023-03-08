package model

type User struct {
	Nickname string `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`
	Email    string `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`
	Phone    string `json:"phone" gorm:"column:phone" validate:"omitempty"`
}
