package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"strings"
	"time"

	"github.com/PlayerLog/playerlog/internal/database"
	"github.com/PlayerLog/playerlog/internal/models"
	"github.com/PlayerLog/playerlog/internal/types"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/lucsky/cuid"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"golang.org/x/crypto/bcrypt"
)

var dbTimeout = time.Second * 3

type Service struct {
	db database.Service
}

func NewService(db database.Service) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateAccount(values types.RegisterTeamValues) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	id, err := cuid.NewCrypto(rand.Reader)

	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(values.Password), bcrypt.DefaultCost)

	user := models.Users.Insert(&models.UserSetter{
		ID:           omit.From(id),
		FirstName:    omitnull.From(values.FirstName),
		LastName:     omitnull.From(values.LastName),
		Email:        omit.From(strings.ReplaceAll(strings.ToLower(values.Email), " ", "")),
		PasswordHash: omit.From(string(password)),
		Role:         omit.From("User"),
		Phone:        omitnull.From(""),
		IsActive:     omitnull.From(true),
		LastLogin:    omitnull.From(time.Now()),
	})

	_, err = user.Exec(ctx, s.db.GetDB())

	if err != nil {
		if models.ErrUniqueConstraint.Is(err) || models.UserErrors.ErrUniqueEmail.Is(err) {
			return ErrUserAlreadyExist
		}
		return err
	}

	return nil
}

func (s *Service) GetUsers() (models.UserSlice, error) {
	users := models.Users.Query()

	usrs, err := users.All(context.Background(), s.db.GetDB())

	if err != nil {
		return nil, err
	}

	return usrs, nil
}

func (s *Service) ValidateUser(email, password string) (types.ValidateUserParams, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	user, err := models.Users.Query(
		sm.Columns(models.UserColumns.ID, models.UserColumns.Email, models.UserColumns.PasswordHash),
		models.SelectWhere.Users.Email.EQ(
			strings.ReplaceAll(strings.ToLower(email), " ", ""))).
		One(ctx, s.db.GetDB())

	if err != nil {
		if models.UserErrors.ErrUniqueEmail.Is(err) {
			return types.ValidateUserParams{}, ErrNoUserFound
		}
		if err == sql.ErrNoRows {
			return types.ValidateUserParams{}, ErrNoUserFound
		}
		return types.ValidateUserParams{}, err
	}

	if user == nil {
		return types.ValidateUserParams{}, ErrNoUserFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return types.ValidateUserParams{}, ErrWrongEmailOrPassword
	}

	return types.ValidateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}, nil
}
