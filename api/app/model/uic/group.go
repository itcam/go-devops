package uic

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Group struct {
	gorm.Model
	GroupName string `gorm:"type:varchar(64);not null;column:group_name"`
	Active    string `gorm:"type:varchar(64);not null;column:Active"`
}

type GroupSerializer struct {
	ID        uint       `json:"id"`
	CreateAt  time.Time  `json:"create_at"`
	UpdateAt  time.Time  `json:"update_at"`
	DeleteAt  *time.Time `json:"delete_at"`
	GroupName string     `json:"group_name"`
	Active    string     `json:"Active"`
}

func (g *Group) Serializer() GroupSerializer {
	return GroupSerializer{
		ID:        g.ID,
		CreateAt:  g.CreatedAt,
		UpdateAt:  g.UpdatedAt,
		DeleteAt:  g.DeletedAt,
		GroupName: g.GroupName,
		Active:    g.Active,
	}
}

func (Group) TableName() string {
	return "group"
}
