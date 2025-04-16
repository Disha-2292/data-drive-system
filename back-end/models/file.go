package models

import "gorm.io/gorm"

type File struct {
	BaseModel
	Name        string `json:"name"`
	Type        string `json:"type"` // "file" or "folder"
	Path        string `json:"path"` // only for files
	UserID      uint   `json:"user_id"`
	ParentID    *uint  `json:"parent_id"`                                   // nullable for root folders
	Parent      *File  `gorm:"foreignKey:ParentID" json:"parent,omitempty"` // Parent folder (for nested folders)
	Children    []File `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Version     *uint  `json:"version,omitempty"`     // Versioning support (for files)
	Size        int64  `json:"size,omitempty"`        // Size of the file (optional)
	Permissions string `json:"permissions,omitempty"` // Permissions (read, write, etc.)
	SharedUsers []User `gorm:"many2many:file_permissions;" json:"shared_users,omitempty"`
}

type FilePermission struct {
	gorm.Model
	FileID     uint   `json:"file_id"`
	UserID     uint   `json:"user_id"`
	Permission string `json:"permission"` // read, write, etc.
}
