package sql

import (
	"context"
	"github.com/challenge/pkg/helpers"
	"github.com/challenge/pkg/models"
	"github.com/challenge/pkg/modules/storage"
)

type sqlUserRepository struct {
	dbCnn storage.IDBConnection
}

func NewSqlUserRepository (sessionName string) (*sqlUserRepository, error) {
	//get connection string
	dbConn, err := storage.DBManagerSingleton().CreateConnection(sessionName)

	if err != nil {
		return nil, err
	}
	return &sqlUserRepository{
		dbCnn: dbConn,
	}, nil
}


func (s *sqlUserRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	query :=
	`INSERT INTO Users (
		Username,
		Password
	)
	VALUES (
		@Username,
		@Password
	);

	SELECT last_insert_rowid();`

	res, err := s.dbCnn.Exec(ctx, query, user.Username, user.Password)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	user.Id = int (id)
	return int(id),  nil
}

func (s *sqlUserRepository) ExistUsername(ctx context.Context, username string) (bool, error) {
	user, err := s.doGetUser(ctx, nil, &username)
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

func (s *sqlUserRepository) GetPassword(ctx context.Context, userName string) (string, error) {
	user, err := s.doGetUser(ctx, nil, &userName)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", nil
	}

	return user.Password, nil
}

func (s *sqlUserRepository) GetProfileById(ctx context.Context, userId int) (*models.UserProfile, error) {
	user, err := s.doGetUser(ctx, &userId, nil)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &models.UserProfile{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}

func (s sqlUserRepository) GetProfileByUsername(ctx context.Context, username string) (*models.UserProfile, error) {
	user, err := s.doGetUser(ctx, nil, &username)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &models.UserProfile{
		Id:       user.Id,
		Username: user.Username,
	}, nil
}

func (s sqlUserRepository) doGetUser(ctx context.Context, userId *int, username *string) (*models.User, error) {
	query :=
	`SELECT Id,
		Username,
		Password
	FROM Users
	WHERE (@userId is null OR @userId = Id) AND (@userName is null OR @userName = Username)`

	rows, err := s.dbCnn.Query(ctx, query, helpers.GetIntOrNullInt(userId), helpers.GetStringOrNullString(username))
	defer func() {
		if rows != nil {
			rows.Close();
		}
	}()

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		row, err := helpers.GetMapFromReader(rows)
		if err != nil {
			return nil, err
		}

		user := &models.User{
			Id:       helpers.GetInt("Id", row),
			Username: helpers.GetString("Username", row),
			Password: helpers.GetString("Password", row),
		}

		return user, nil
	}

	return nil, nil
}

