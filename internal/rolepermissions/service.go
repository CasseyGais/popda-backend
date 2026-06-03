package rolepermissions

import (
	"errors"
	"fmt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

// Request structs
type CreateRolePermissionRequest struct {
	RoleID       uint `json:"role_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
}

type RolePermissionResponse struct {
	RoleID       uint `json:"role_id"`
	PermissionID uint `json:"permission_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// GetAll role permissions
func (s *Service) GetAll() ([]RolePermissionResponse, error) {
	rolePermissions, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	var responses []RolePermissionResponse
	for _, rp := range rolePermissions {
		responses = append(responses, RolePermissionResponse{
			RoleID:       rp.RoleID,
			PermissionID: rp.PermissionID,
			CreatedAt:    rp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    rp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return responses, nil
}

// GetByRoleID returns all permissions for a specific role
func (s *Service) GetByRoleID(roleID uint) ([]RolePermissionResponse, error) {
	if roleID == 0 {
		return nil, errors.New("role ID is required")
	}

	rolePermissions, err := s.repository.GetByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions by role ID: %w", err)
	}

	var responses []RolePermissionResponse
	for _, rp := range rolePermissions {
		responses = append(responses, RolePermissionResponse{
			RoleID:       rp.RoleID,
			PermissionID: rp.PermissionID,
			CreatedAt:    rp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    rp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return responses, nil
}

// GetByPermissionID returns all roles that have a specific permission
func (s *Service) GetByPermissionID(permissionID uint) ([]RolePermissionResponse, error) {
	if permissionID == 0 {
		return nil, errors.New("permission ID is required")
	}

	rolePermissions, err := s.repository.GetByPermissionID(permissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions by permission ID: %w", err)
	}

	var responses []RolePermissionResponse
	for _, rp := range rolePermissions {
		responses = append(responses, RolePermissionResponse{
			RoleID:       rp.RoleID,
			PermissionID: rp.PermissionID,
			CreatedAt:    rp.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    rp.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return responses, nil
}

// AssignPermissionToRole assigns a permission to a role
func (s *Service) AssignPermissionToRole(request CreateRolePermissionRequest) (*RolePermissionResponse, error) {
	if request.RoleID == 0 {
		return nil, errors.New("role ID is required")
	}
	if request.PermissionID == 0 {
		return nil, errors.New("permission ID is required")
	}

	// Check if the assignment already exists
	existing, err := s.repository.GetByRoleID(request.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role permission: %w", err)
	}

	for _, rp := range existing {
		if rp.PermissionID == request.PermissionID {
			return nil, errors.New("permission already assigned to this role")
		}
	}

	rolePermission := &RolePermission{
		RoleID:       request.RoleID,
		PermissionID: request.PermissionID,
	}

	if err := s.repository.Create(rolePermission); err != nil {
		return nil, fmt.Errorf("failed to assign permission to role: %w", err)
	}

	response := &RolePermissionResponse{
		RoleID:       rolePermission.RoleID,
		PermissionID: rolePermission.PermissionID,
		CreatedAt:    rolePermission.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    rolePermission.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// RemovePermissionFromRole removes a permission from a role
func (s *Service) RemovePermissionFromRole(roleID, permissionID uint) error {
	if roleID == 0 {
		return errors.New("role ID is required")
	}
	if permissionID == 0 {
		return errors.New("permission ID is required")
	}

	if err := s.repository.Delete(roleID, permissionID); err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}

	return nil
}

// DeleteByRoleID removes all permissions from a role
func (s *Service) DeleteByRoleID(roleID uint) error {
	if roleID == 0 {
		return errors.New("role ID is required")
	}

	if err := s.repository.DeleteByRoleID(roleID); err != nil {
		return fmt.Errorf("failed to delete role permissions by role ID: %w", err)
	}

	return nil
}

// DeleteByPermissionID removes a permission from all roles
func (s *Service) DeleteByPermissionID(permissionID uint) error {
	if permissionID == 0 {
		return errors.New("permission ID is required")
	}

	if err := s.repository.DeleteByPermissionID(permissionID); err != nil {
		return fmt.Errorf("failed to delete role permissions by permission ID: %w", err)
	}

	return nil
}
