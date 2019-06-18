package db

import "tech-park-db-hw/internal/pkg/models"

func CreateUser(user *models.User) error {
	_, err := Exec(CreateUserQuery, user.Nickname, user.Email, user.Fullname, user.About)
	return err
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
