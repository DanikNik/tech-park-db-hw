package db

import (
	"github.com/jackc/pgx"
	"tech-park-db-hw/internal/pkg/models"
)

func CreateUser(user *models.User) error {
	_, err := Exec(CreateUserQuery, user.Nickname, user.Email, user.Fullname, user.About)
	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			if pqError.Code == uniqueIntegrityError {
				return ErrConflict
			}
		}
		return err
	}
	return nil
}

func SelectUsersOnConflict(nickname, email string) ([]models.User, error) {
	rows, err := Query(SelectUsersWithNickOrEmail, nickname, email)
	if err != nil {
		return nil, err
	}

	alikeUsers := []models.User{}
	defer rows.Close()
	for rows.Next() {
		user := models.User{}
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
		if err != nil {
			return nil, err
		}
		alikeUsers = append(alikeUsers, user)
	}
	return alikeUsers, nil
}

func GetUser(nickname string) (*models.User, error) {
	userObj, err := QueryRow(GetUserQuery, nickname)
	if err != nil {
		return nil, err
	}
	var user models.User
	err = userObj.Scan(&user.Nickname, &user.Email, &user.Fullname, &user.About)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
