package services

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/services/auth/repositories"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/envutil"
	"github.com/IrusHunter/duckademic/shared/events"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	platform.BaseService[entities.User]
	ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error
	Login(ctx context.Context, login, password string) (ResponseUser, error)
	Refresh(ctx context.Context, accessTokenStr, refreshTokenStr string) (newAT, newRT string, err error)
	ResetPassword(context.Context, uuid.UUID) error
}

func NewUserService(
	ur repositories.UserRepository,
	rr repositories.RoleRepository,
	rpr repositories.RolePermissionsRepository,
	eb events.EventBus,
	defaultPassword string,
	jwtSecret []byte,
	superAdminLogin, superAdminPassword, superAdminRole string,
) UserService {
	sc := platform.NewServiceConfig(
		"UserService",
		filepath.Join("data", "users.json"),
		"user",
	)

	res := &userService{
		repository:               ur,
		roleRepository:           rr,
		rolePermissionRepository: rpr,
		defaultPassword:          defaultPassword,
		jwtSecret:                jwtSecret,
	}
	res.BaseService = platform.NewBaseService(sc, ur,
		map[platform.ServiceExternalFuncType]platform.ServiceExternalFunc[entities.User]{
			platform.OnAddPrepare: res.onAddPrepare,
		},
	)

	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.StudentRT),
		res.studentEventHandler,
	)
	eb.Subscribe(
		contextutil.SetTraceID(context.Background()),
		string(events.TeacherRT),
		res.teacherEventHandler,
	)

	res.Seed(contextutil.SetTraceID(context.Background()))
	res.createSuperAdmin(contextutil.SetTraceID(context.Background()), superAdminLogin, superAdminPassword, superAdminRole)

	return res
}

type userService struct {
	platform.BaseService[entities.User]
	repository               repositories.UserRepository
	roleRepository           repositories.RoleRepository
	rolePermissionRepository repositories.RolePermissionsRepository
	defaultPassword          string
	jwtSecret                []byte
	studentRoleID            uuid.UUID
	teacherRoleID            uuid.UUID
}

func (s *userService) onAddPrepare(ctx context.Context, u *entities.User) error {
	var err error
	u.HashedPassword, err = s.hashPassword(s.defaultPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.IsDefaultPassword = true

	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	return nil
}

func (s *userService) studentEventHandler(ctx context.Context, b []byte) {
	studentEvent, err := events.FromByteConvertor[events.StudentRE](b)
	if err != nil {
		s.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx), "StudentRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.GetLogger().Log(contextutil.GetTraceID(ctx), "StudentRTHandler",
		fmt.Sprintf("received %s", studentEvent), logger.EventDataReceived,
	)

	u := entities.User{
		ID:     studentEvent.ID,
		Login:  studentEvent.Email,
		RoleID: s.studentRoleID,
	}

	switch studentEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, u)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, studentEvent.ID, u)
	case events.EntityDeleted:
		s.Delete(ctx, studentEvent.ID)
	}
}
func (s *userService) teacherEventHandler(ctx context.Context, b []byte) {
	teacherEvent, err := events.FromByteConvertor[events.TeacherRE](b)
	if err != nil {
		s.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx), "TeacherRTHandler",
			err, logger.EventDataReadFailed)
		return
	}

	s.GetLogger().Log(contextutil.GetTraceID(ctx), "TeacherRTHandler",
		fmt.Sprintf("received %s", teacherEvent), logger.EventDataReceived,
	)

	u := entities.User{
		ID:     teacherEvent.ID,
		Login:  teacherEvent.Email,
		RoleID: s.teacherRoleID,
	}

	switch teacherEvent.Event {
	case events.EntityCreated:
		s.Add(ctx, u)
	case events.EntityUpdated:
		s.ExternalUpdate(ctx, teacherEvent.ID, u)
	case events.EntityDeleted:
		s.Delete(ctx, teacherEvent.ID)
	}
}

func (s *userService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func (s *userService) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func (s *userService) createSuperAdmin(ctx context.Context, login, password, role string) {
	adminRole := s.roleRepository.FindByName(ctx, role)
	if adminRole == nil {
		adminRole = &entities.Role{}
	}
	s.Add(ctx, entities.User{
		ID:     uuid.New(),
		Login:  login,
		RoleID: adminRole.ID,
	})
	user := s.repository.FindByLogin(ctx, login)
	if user != nil {
		user.Password = &s.defaultPassword
		s.ChangePassword(contextutil.SetTraceID(context.Background()), user.ID, s.defaultPassword, password)
	}
}

func (s *userService) ChangePassword(ctx context.Context, id uuid.UUID, currentPassword, newPassword string) error {
	traceID := contextutil.GetTraceID(ctx)

	user := s.repository.FindByID(ctx, id)
	if user == nil {
		return s.GetLogger().LogAndReturnError(
			traceID,
			"ChangePassword",
			fmt.Errorf("user not found"),
			logger.ServiceDataFetchFailed,
		)
	}

	if !s.checkPassword(currentPassword, user.HashedPassword) {
		s.GetLogger().Log(
			contextutil.GetTraceID(ctx),
			"ChangePassword",
			"old password is incorrect",
			logger.ServiceValidationFailed,
		)
		return fmt.Errorf("invalid old password")
	}

	hashed, err := s.hashPassword(newPassword)
	if err != nil {
		return s.GetLogger().LogAndReturnError(
			traceID,
			"ChangePassword",
			err,
			logger.ServiceDataFetchFailed,
		)
	}

	user.HashedPassword = hashed
	user.IsDefaultPassword = false

	_, err = s.repository.Update(ctx, user.ID, *user)
	if err != nil {
		return s.GetLogger().LogAndReturnError(
			traceID,
			"ChangePassword",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	s.GetLogger().Log(
		traceID,
		"ChangePassword",
		"user password successfully changed",
		logger.ServiceOperationSuccess,
	)

	return nil
}

func (s *userService) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	user entities.User,
) (entities.User, error) {
	updatedU, err := s.repository.ExternalUpdate(ctx, id, user)
	if err != nil {
		return entities.User{}, s.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"ExternalUpdate",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	s.GetLogger().Log(contextutil.GetTraceID(ctx), "ExternalUpdate",
		fmt.Sprintf("%s successfully updated", updatedU),
		logger.ServiceOperationSuccess,
	)

	return updatedU, nil
}

func (s *userService) Seed(ctx context.Context) error {
	studentRoleStr := "student"
	teacherRoleStr := "teacher"

	studentRole := s.roleRepository.FindByName(ctx, studentRoleStr)
	if studentRole == nil {
		return s.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to find student role: %s", studentRoleStr), logger.ServiceDataFetchFailed)
	}

	teacherRole := s.roleRepository.FindByName(ctx, teacherRoleStr)
	if teacherRole == nil {
		return s.GetLogger().LogAndReturnError(contextutil.GetTraceID(ctx), "Seed",
			fmt.Errorf("failed to find teacher role: %s", teacherRoleStr), logger.ServiceDataFetchFailed)
	}

	s.studentRoleID = studentRole.ID
	s.teacherRoleID = teacherRole.ID

	assignments := []struct {
		Login    string `json:"login"`
		RoleName string `json:"role_name"`
	}{}

	if err := jsonutil.ReadFileTo(filepath.Join("data", "users.json"), &assignments); err != nil {
		return s.GetLogger().LogAndReturnError(
			contextutil.GetTraceID(ctx),
			"Seed",
			fmt.Errorf("failed to load users seed data: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	var lastError error

	for _, item := range assignments {
		role := s.roleRepository.FindByName(ctx, item.RoleName)
		if role == nil {
			lastError = s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("role not found for user %s: %s", item.Login, item.RoleName),
				logger.ServiceValidationFailed,
			)
			continue
		}

		user := entities.User{
			ID:     uuid.New(),
			Login:  item.Login,
			RoleID: role.ID,
		}

		_, err := s.Add(ctx, user)
		if err != nil {
			lastError = s.GetLogger().LogAndReturnError(
				contextutil.GetTraceID(ctx),
				"Seed",
				fmt.Errorf("failed to add user %s: %w", user.Login, err),
				logger.ServiceValidationFailed,
			)
			continue
		}
	}

	s.createSuperAdmin(contextutil.SetTraceID(context.Background()),
		envutil.GetStringFromENV("SUPER_ADMIN_LOGIN"),
		envutil.GetStringFromENV("SUPER_ADMIN_PASSWORD"),
		envutil.GetStringFromENV("SUPER_ADMIN_ROLE"),
	)

	s.GetLogger().Log(
		contextutil.GetTraceID(ctx),
		"Seed",
		"users seed processed",
		logger.ServiceOperationSuccess,
	)

	return lastError
}

type ResponseUser struct {
	ID                uuid.UUID `json:"id"`
	Login             string    `json:"login"`
	Role              string    `json:"role"`
	IsDefaultPassword bool      `json:"is_default_password"`
	AccessToken       string    `json:"access_token"`
	RefreshToken      string    `json:"refresh_token"`
}

func (s *userService) Login(ctx context.Context, login, password string) (ResponseUser, error) {
	traceID := contextutil.GetTraceID(ctx)

	user := s.repository.FindByLogin(ctx, login)
	if user == nil {
		return ResponseUser{}, s.GetLogger().LogAndReturnError(
			traceID,
			"Login",
			fmt.Errorf("user not found"),
			logger.ServiceDataFetchFailed,
		)
	}

	user = s.repository.Fill(ctx, user.ID)
	if user == nil {
		return ResponseUser{}, s.GetLogger().LogAndReturnError(
			traceID,
			"Login",
			fmt.Errorf("user not found"),
			logger.ServiceDataFetchFailed,
		)
	}

	if !s.checkPassword(password, user.HashedPassword) {
		return ResponseUser{}, s.GetLogger().LogAndReturnError(
			traceID,
			"Login",
			fmt.Errorf("invalid credentials"),
			logger.ServiceValidationFailed,
		)
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return ResponseUser{}, s.GetLogger().LogAndReturnError(
			traceID,
			"Login",
			fmt.Errorf("failed to generate tokens: %w", err),
			logger.ServiceValidationFailed,
		)
	}

	s.GetLogger().Log(traceID, "Login",
		fmt.Sprintf("user %s logged in", user.Login),
		logger.ServiceOperationSuccess,
	)

	return ResponseUser{
		ID:                user.ID,
		Login:             user.Login,
		Role:              *user.RoleName,
		IsDefaultPassword: user.IsDefaultPassword,
		AccessToken:       accessToken,
		RefreshToken:      refreshToken,
	}, nil
}
func (s *userService) Refresh(ctx context.Context, accessTokenStr, refreshTokenStr string,
) (newAT, newRT string, err error) {
	refreshClaims := &jwt.RegisteredClaims{}

	refreshToken, err := jwt.ParseWithClaims(
		refreshTokenStr,
		refreshClaims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return s.jwtSecret, nil
		},
	)

	if err != nil || !refreshToken.Valid {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	userID := refreshClaims.Subject
	if userID == "" {
		return "", "", fmt.Errorf("invalid refresh token claims")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", "", fmt.Errorf("invalid user id in token")
	}

	user := s.repository.FindByID(ctx, userUUID)
	if user == nil {
		return "", "", fmt.Errorf("user not found")
	}

	claims, err := jsonutil.ParseAccessToken(accessTokenStr, s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	if claims.UserID != userID {
		return "", "", fmt.Errorf("invalid token")
	}

	return s.generateTokens(ctx, user)
}
func (s *userService) ResetPassword(
	ctx context.Context,
	userID uuid.UUID,
) error {
	traceID := contextutil.GetTraceID(ctx)

	user := s.repository.FindByID(ctx, userID)
	if user == nil {
		return fmt.Errorf("user not found")
	}

	hashed, err := s.hashPassword(s.defaultPassword)
	if err != nil {
		return s.GetLogger().LogAndReturnError(
			traceID,
			"ResetPassword",
			err,
			logger.ServiceDataFetchFailed,
		)
	}

	user.HashedPassword = hashed
	user.IsDefaultPassword = true

	_, err = s.repository.Update(ctx, user.ID, *user)
	if err != nil {
		return s.GetLogger().LogAndReturnError(
			traceID,
			"ResetPassword",
			err,
			logger.ServiceRepositoryFailed,
		)
	}

	s.GetLogger().Log(traceID, "ResetPassword",
		fmt.Sprintf("password reset to default for user %s", user.Login),
		logger.ServiceOperationSuccess,
	)

	return nil
}

func (s *userService) generateTokens(ctx context.Context, user *entities.User) (accessT, refreshT string, err error) {
	permissions, err := s.rolePermissionRepository.GetPermissionsByRoleID(ctx, user.RoleID)
	if err != nil {
		return "", "", err
	}

	now := time.Now()

	accessClaims := jsonutil.AccessClaims{
		UserID:      user.ID.String(),
		Login:       user.Login,
		Role:        *user.RoleName,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
