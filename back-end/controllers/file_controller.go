package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Disha-2292/data-drive-system/config"
	"github.com/Disha-2292/data-drive-system/models"
	"github.com/gin-gonic/gin"
)

func CreateFileOrFolder(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Type     string `json:"type"`      // "file" or "folder"
		ParentID *uint  `json:"parent_id"` // nullable
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	file := models.File{
		Name:     input.Name,
		Type:     input.Type,
		ParentID: input.ParentID,
		UserID:   userID,
	}

	if err := config.DB.Create(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	c.JSON(http.StatusCreated, file)
}

func GetUserFiles(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	var files []models.File

	if err := config.DB.Where("user_id = ?", userID).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func GetFileByID(c *gin.Context) {
	id := c.Param("id")
	var file models.File

	if err := config.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func UpdateFile(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Name     string `json:"name"`
		ParentID *uint  `json:"parent_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var file models.File
	if err := config.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if the user has 'write' permission
	if !CheckFilePermission(c, file.ID, "write") {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this file"})
		return
	}

	file.Name = input.Name
	file.ParentID = input.ParentID

	config.DB.Save(&file)

	c.JSON(http.StatusOK, file)
}

func DeleteFile(c *gin.Context) {
	id := c.Param("id")
	var file models.File

	if err := config.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Delete actual file if it's a real file
	if file.Type == "file" && file.Path != "" {
		os.Remove(file.Path)
	}

	config.DB.Delete(&file)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	userID := c.MustGet("user_id").(uint)
	parentIDStr := c.PostForm("parent_id")

	var parentID *uint
	if parentIDStr != "" {
		id, _ := strconv.ParseUint(parentIDStr, 10, 64)
		tmp := uint(id)
		parentID = &tmp
	}

	dstPath := "uploads/" + filepath.Base(fileHeader.Filename)
	if err := c.SaveUploadedFile(fileHeader, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	file := models.File{
		Name:     fileHeader.Filename,
		Type:     "file",
		Path:     dstPath,
		UserID:   userID,
		ParentID: parentID,
	}

	config.DB.Create(&file)
	c.JSON(http.StatusOK, file)
}

func DownloadFile(c *gin.Context) {
	id := c.Param("id")
	var file models.File

	if err := config.DB.First(&file, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if the user has 'read' permission on the file
	if !CheckFilePermission(c, file.ID, "read") {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to download this file"})
		return
	}

	if file.Type != "file" || file.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not a downloadable file"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+file.Name)
	c.File(file.Path)
}

func GetAllFiles(c *gin.Context) {
	var files []models.File
	if err := config.DB.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	c.JSON(http.StatusOK, files)
}

func ShareFile(c *gin.Context) {
	var input struct {
		UserIDs    []uint `json:"user_ids"`
		Permission string `json:"permission"` // "read" or "write"
	}

	// Get the file ID from the URL parameter
	fileID := c.Param("id")

	// Bind the request input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the file exists
	var file models.File
	if err := config.DB.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Ensure the current user owns the file before sharing
	userID := c.MustGet("user_id").(uint)
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this file"})
		return
	}

	// Share the file with the users and set their permissions
	for _, userID := range input.UserIDs {
		filePermission := models.FilePermission{
			FileID:     file.ID,
			UserID:     userID,
			Permission: input.Permission,
		}
		if err := config.DB.Create(&filePermission).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share file"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "File shared successfully"})
}

func CheckFilePermission(c *gin.Context, fileID uint, requiredPermission string) bool {
	userID := c.MustGet("user_id").(uint)
	var filePermission models.FilePermission

	// Check if the user has the required permission on the file
	if err := config.DB.Where("file_id = ? AND user_id = ? AND permission = ?", fileID, userID, requiredPermission).First(&filePermission).Error; err != nil {
		return false
	}

	return true
}

func GetFilePermissions(c *gin.Context) {
	// Get the file ID from the URL parameter
	fileID := c.Param("id")
	userID := c.MustGet("user_id").(uint)

	// Find the file permission for the current user
	var filePermission models.FilePermission
	if err := config.DB.Where("file_id = ? AND user_id = ?", fileID, userID).First(&filePermission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permissions not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"file_id":    filePermission.FileID,
		"user_id":    filePermission.UserID,
		"permission": filePermission.Permission,
	})
}
func UpdateFilePermissions(c *gin.Context) {
	var input struct {
		Permission string `json:"permission"` // "read" or "write"
	}

	// Get the file ID from the URL parameter
	fileID := c.Param("id")

	// Bind the request input to the struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the file exists
	var file models.File
	if err := config.DB.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Ensure the current user is the owner of the file
	userID := c.MustGet("user_id").(uint)
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't own this file"})
		return
	}

	// Update the file permissions
	var filePermission models.FilePermission
	if err := config.DB.Where("file_id = ? AND user_id = ?", fileID, userID).First(&filePermission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File permissions not found"})
		return
	}

	filePermission.Permission = input.Permission
	if err := config.DB.Save(&filePermission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions updated successfully"})
}

func SearchFiles(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	// Query filters
	name := c.Query("name")
	fileType := c.Query("type")
	size := c.Query("size")

	// Sorting & Pagination
	sortBy := c.DefaultQuery("sort", "created_at") // default sort by created_at
	order := c.DefaultQuery("order", "desc")       // default order is desc
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Build query
	var files []models.File
	query := config.DB.Model(&models.File{}).Where("user_id = ?", userID)

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if fileType != "" {
		query = query.Where("type = ?", fileType)
	}
	if size != "" {
		if strings.HasPrefix(size, ">") {
			query = query.Where("size > ?", size[1:])
		} else if strings.HasPrefix(size, "<") {
			query = query.Where("size < ?", size[1:])
		} else if strings.HasPrefix(size, "=") {
			query = query.Where("size = ?", size[1:])
		}
	}

	// Sorting
	query = query.Order(sortBy + " " + order)

	// Pagination
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"limit":   limit,
		"results": len(files),
		"data":    files,
	})
}
