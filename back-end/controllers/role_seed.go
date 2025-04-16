// controllers/role_seed.go

package controllers

import (
	"log"

	"github.com/Disha-2292/data-drive-system/config"
	"github.com/Disha-2292/data-drive-system/models"
)

func SeedRoles() {
	adminRole := models.Role{Name: "admin", Description: "Admin Role with full permissions"}
	userRole := models.Role{Name: "user", Description: "Regular User with limited permissions"}

	if err := config.DB.Create(&adminRole).Error; err != nil {
		log.Fatal("Error seeding admin role", err)
	}

	if err := config.DB.Create(&userRole).Error; err != nil {
		log.Fatal("Error seeding user role", err)
	}
}
