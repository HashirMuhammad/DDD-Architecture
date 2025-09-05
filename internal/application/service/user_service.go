package service

import (
	"context"
	"ddd-user-service/internal/application/dto"
	"ddd-user-service/internal/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	emailExists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, domain.ErrEmailExists
	}

	usernameExists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, domain.ErrUsernameExists
	}

	user, err := domain.NewUser(req.Name, req.Email, req.Username)
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return s.userToResponse(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	userID := domain.UserID(id)
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.userToResponse(user), nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = s.userToResponse(user)
	}

	return responses, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	userID := domain.UserID(id)
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		if err := user.UpdateName(*req.Name); err != nil {
			return nil, err
		}
	}

	if req.Email != nil {
		emailExists, err := s.userRepo.ExistsByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if emailExists && user.Email != *req.Email {
			return nil, domain.ErrEmailExists
		}
		if err := user.UpdateEmail(*req.Email); err != nil {
			return nil, err
		}
	}

	if req.Username != nil {
		usernameExists, err := s.userRepo.ExistsByUsername(ctx, *req.Username)
		if err != nil {
			return nil, err
		}
		if usernameExists && user.Username != *req.Username {
			return nil, domain.ErrUsernameExists
		}
		if err := user.UpdateUsername(*req.Username); err != nil {
			return nil, err
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return s.userToResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	userID := domain.UserID(id)
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(ctx, userID)
}

func (s *UserService) userToResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:       user.ID.String(),
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
	}
}
