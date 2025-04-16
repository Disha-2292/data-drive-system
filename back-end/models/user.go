package models

type User struct {
	BaseModel
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
	Files    []File `gorm:"foreignKey:UserID"`
	RoleID   uint   `json:"role_id"` // Reference to the Role
	Role     Role   `json:"role"`    // Role struct included for easier access
}
