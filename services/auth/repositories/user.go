package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/services/auth/entities"
	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/logger"
	"github.com/IrusHunter/duckademic/shared/platform"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	platform.BaseRepository[entities.User]
	FindByLogin(context.Context, string) *entities.User
	ExternalUpdate(context.Context, uuid.UUID, entities.User) (entities.User, error)
	UpdatePassword(context.Context, uuid.UUID, entities.User) error
	Fill(context.Context, uuid.UUID) *entities.User
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	config := platform.NewRepositoryConfig(
		"UserRepository",
		entities.User{}.TableName(),
		entities.User{}.EntityName(),
		[]string{"id", "login", "password", "is_default_password", "role_id", "last_login"},
		[]string{"password", "is_default_password", "role_id"},
		[]string{"created_at", "updated_at"},
	)

	ur := &userRepository{
		BaseRepository: platform.NewBaseRepository[entities.User](config, db),
		db:             db,
	}
	ur.logger = ur.GetLogger()

	return ur
}

type userRepository struct {
	platform.BaseRepository[entities.User]
	db     *sqlx.DB
	logger logger.Logger
}

func (r *userRepository) FindByLogin(ctx context.Context, login string) *entities.User {
	return r.FindFirstBy(ctx, "login", login)
}
func (r *userRepository) ExternalUpdate(
	ctx context.Context,
	id uuid.UUID,
	user entities.User,
) (entities.User, error) {
	return r.UpdateFields(ctx, id, []string{
		"login",
	}, user)
}
func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, user entities.User) error {
	_, err := r.UpdateFields(ctx, id, []string{"password", "is_default_password"}, user)
	return err
}

type UserFlat struct {
	ID                uuid.UUID `db:"id"`
	Login             string    `db:"login"`
	HashedPassword    string    `db:"password"`
	IsDefaultPassword bool      `db:"is_default_password"`
	RoleID            uuid.UUID `db:"role_id"`
	LastLogin         time.Time `db:"last_login"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`

	RoleName string `db:"role_name"`
}

func ConvertToUser(flat UserFlat) *entities.User {
	user := &entities.User{
		ID:                flat.ID,
		Login:             flat.Login,
		HashedPassword:    flat.HashedPassword,
		IsDefaultPassword: flat.IsDefaultPassword,
		RoleID:            flat.RoleID,
		LastLogin:         flat.LastLogin,
		CreatedAt:         flat.CreatedAt,
		UpdatedAt:         flat.UpdatedAt,
		RoleName:          &flat.RoleName,
	}

	return user
}

func (r *userRepository) Fill(ctx context.Context, id uuid.UUID) *entities.User {
	query := fmt.Sprintf(`
		SELECT 
			u.id,
			u.login,
			u.password,
			u.is_default_password,
			u.role_id,
			u.last_login,
			u.created_at,
			u.updated_at,

			r.name AS role_name

		FROM %s u
		LEFT JOIN %s r ON u.role_id = r.id
		WHERE u.id = $1;
	`, entities.User{}.TableName(), entities.Role{}.TableName())

	var flat UserFlat

	if err := r.db.GetContext(ctx, &flat, query, id); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			r.logger.Log(contextutil.GetTraceID(ctx), "Fill",
				fmt.Sprintf("user with id %q not found", id),
				logger.RepositoryOperationSuccess,
			)
			return nil
		}

		r.logger.LogAndReturnError(contextutil.GetTraceID(ctx), "Fill",
			fmt.Errorf("failed to scan user with id %q: %w", id, err),
			logger.RepositoryScanFailed,
		)
		return nil
	}

	user := ConvertToUser(flat)

	r.logger.Log(contextutil.GetTraceID(ctx), "Fill",
		fmt.Sprintf("user %s found with role and permissions", user.Login),
		logger.RepositoryOperationSuccess,
	)

	return user
}
