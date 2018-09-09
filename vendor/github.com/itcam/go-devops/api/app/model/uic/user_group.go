package uic

import (
	"github.com/jinzhu/gorm"
	"time"
)

type UserGroup struct {
	gorm.Model
	UserId  uint `gorm:"type: int(4); not null; column:UserId"`
	GroupId uint `gorm:"type: int(4); not null; column:GroupId"`
}

type UserGroupSerializer struct {
	ID       uint       `json:"id"`
	CreateAt time.Time  `json:"create_at"`
	UpdateAt time.Time  `json:"update_at"`
	DeleteAt *time.Time `json:"delete_at"`
	UserId   uint       `json:"userId"`
	GroupId  uint       `json:"groupId"`
}

func (g *UserGroup) Serializer() UserGroupSerializer {
	return UserGroupSerializer{
		ID:       g.ID,
		CreateAt: g.CreatedAt,
		UpdateAt: g.UpdatedAt,
		DeleteAt: g.DeletedAt,
		UserId:   g.UserId,
		GroupId:  g.GroupId,
	}
}
func (UserGroup) TableName() string {
	return "user_group"
}
