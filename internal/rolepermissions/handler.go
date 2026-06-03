package rolepermissions

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAll godoc
// @Summary Get all role permissions
// @Description Get all role permissions
// @Tags role-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/role-permissions [get]
func (h *Handler) GetAll(c *gin.Context) {
	rolePermissions, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get role permissions",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions retrieved successfully",
		"data":    rolePermissions,
	})
}

// GetByRoleID godoc
// @Summary Get permissions by role ID
// @Description Get all permissions assigned to a specific role
// @Tags role-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/role-permissions/role/{id} [get]
func (h *Handler) GetByRoleID(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role ID",
			"error":   err.Error(),
		})
		return
	}

	rolePermissions, err := h.service.GetByRoleID(uint(roleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get role permissions by role ID",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions retrieved successfully",
		"data":    rolePermissions,
	})
}

// GetByPermissionID godoc
// @Summary Get roles by permission ID
// @Description Get all roles that have a specific permission
// @Tags role-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/role-permissions/permission/{id} [get]
func (h *Handler) GetByPermissionID(c *gin.Context) {
	permissionIDStr := c.Param("id")
	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid permission ID",
			"error":   err.Error(),
		})
		return
	}

	rolePermissions, err := h.service.GetByPermissionID(uint(permissionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get role permissions by permission ID",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions retrieved successfully",
		"data":    rolePermissions,
	})
}

// AssignPermissionToRole godoc
// @Summary Assign permission to role
// @Description Assign a permission to a role
// @Tags role-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRolePermissionRequest true "Assign permission to role request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/role-permissions [post]
func (h *Handler) AssignPermissionToRole(c *gin.Context) {
	var request CreateRolePermissionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	rolePermission, err := h.service.AssignPermissionToRole(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to assign permission to role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Permission assigned to role successfully",
		"data":    rolePermission,
	})
}

// RemovePermissionFromRole godoc
// @Summary Remove permission from role
// @Description Remove a permission from a role
// @Tags role-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param permissionId path int true "Permission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/role-permissions/role/{id}/permission/{permissionId} [delete]
func (h *Handler) RemovePermissionFromRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	permissionIDStr := c.Param("permissionId")

	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role ID",
			"error":   err.Error(),
		})
		return
	}

	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid permission ID",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.RemovePermissionFromRole(uint(roleID), uint(permissionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to remove permission from role",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Permission removed from role successfully",
	})
}

// DeleteByRoleID godoc
// @Summary Delete all permissions from role
// @Description Remove all permissions from a specific role
// @Tags role-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/role-permissions/role/{id} [delete]
func (h *Handler) DeleteByRoleID(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid role ID",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.DeleteByRoleID(uint(roleID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete role permissions by role ID",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions deleted successfully",
	})
}

// DeleteByPermissionID godoc
// @Summary Remove permission from all roles
// @Description Remove a specific permission from all roles
// @Tags role-permissions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Permission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /admin/role-permissions/permission/{id} [delete]
func (h *Handler) DeleteByPermissionID(c *gin.Context) {
	permissionIDStr := c.Param("id")
	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid permission ID",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.DeleteByPermissionID(uint(permissionID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete role permissions by permission ID",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Permission removed from all roles successfully",
	})
}
