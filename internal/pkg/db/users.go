package db

import (
	"fmt"
	"github.com/jackc/pgx"
	"tech-park-db-hw/internal/pkg/models"
)

func CreateUser(user *models.User) error {
	_, err := Exec(CreateUserQuery, user.Nickname, user.Email, user.Fullname, user.About)
	//defer rows.Close()
	if err != nil {
		if pqError, ok := err.(pgx.PgError); ok {
			if pqError.Code == uniqueIntegrityError {
				return ErrConflict
			}
		}
		return err
	}
	increaseUserCount()
	return err
}

func SelectUsersOnConflict(nickname, email string) ([]models.User, error) {
	rows, err := dbObj.Query(SelectUsersWithNickOrEmail, nickname, email)
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
	if err == pgx.ErrNoRows {
		return nil, ErrNotFound
	}
	return &user, nil
}

func UpdateUser(nickname string, updateData *models.UserUpdate) (us *models.User, err error) {
	us = &models.User{}
	currentUserData, err := GetUser(nickname)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("WTF")
	}
	if updateData.About != "" {
		currentUserData.About = updateData.About
	}
	if updateData.Email != "" {
		currentUserData.Email = updateData.Email
	}
	if updateData.Fullname != "" {
		currentUserData.Fullname = updateData.Fullname
	}

	row, err := QueryRow(UpdateUserQuery, currentUserData.Fullname, currentUserData.Email, currentUserData.About, currentUserData.Nickname)
	if err != nil {
		panic(err)
	}

	scanErr := row.Scan(&us.Nickname, &us.Fullname, &us.About, &us.Email)
	if scanErr != nil {
		if scanErr == pgx.ErrNoRows {
			return nil, ErrNotFound
		}

		if pqError, ok := scanErr.(pgx.PgError); ok {
			switch pqError.Code {
			case uniqueIntegrityError:
				return nil, ErrConflict
			}
		}

		return us, scanErr
	}
	return us, err
}
