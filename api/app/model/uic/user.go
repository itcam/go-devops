package uic

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"type: varchar(64); not null; column:username"`
	FullName string `gorm:"type: varchar(64); column:fullname"`
	PassWord string `gorm:"type: varchar(64); not null; column:password"`
	Salt     string `gorm:"type: varchar(64); not null; column:salt"`
	Email    string `gorm:"type: varchar(64); not null; column:email"`
	Phone    string `gorm:"type: varchar(64); column:phone"`
	Role     string `gorm:"type: varchar(64); not null; column:role"`
	Active   string `gorm:"type: char(2); not null; column:Active"`
}

type UserSerializer struct {
	ID       uint       `json:"id"`
	CreateAt time.Time  `json:"create_at"`
	UpdateAt time.Time  `json:"update_at"`
	DeleteAt *time.Time `json:"delete_at"`
	UserName string     `json:"username"`
	FullName string     `json:"fullname"`
	PassWord string     `json:"password"`
	Salt     string     `json:"salt"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Role     string     `json:"role"`
	Active   string     `json:"Active"`
}

func (p *User) Serializer() UserSerializer {
	return UserSerializer{
		ID:       p.ID,
		CreateAt: p.CreatedAt,
		UpdateAt: p.UpdatedAt,
		DeleteAt: p.DeletedAt,
		UserName: p.Username,
		FullName: p.FullName,
		PassWord: p.PassWord,
		Salt:     p.Salt,
		Email:    p.Email,
		Phone:    p.Phone,
		Role:     p.Role,
		Active:   p.Active,
	}
}

func (User) TableName() string {
	return "arrow_user"
}
