package builder

import (
	"github.com/okutsen/PasswordManager/model/api"
	"github.com/okutsen/PasswordManager/model/controller"
	"github.com/okutsen/PasswordManager/model/db"
)

func BuildControllerUserFromDBUser(user *db.User) controller.User {
	return controller.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Login:     user.Login,
		Password:  user.Password,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func BuildControllerUsersFromDBUsers(users []db.User) []controller.User {
	usersController := make([]controller.User, len(users))
	for i, v := range users {
		usersController[i] = BuildControllerUserFromDBUser(&v)
	}

	return usersController
}

func BuildAPIUserFromControllerUser(user *controller.User) api.User {
	return api.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Login:     user.Login,
		Password:  user.Password,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func BuildAPIUsersFromControllerUsers(users []controller.User) []api.User {
	usersController := make([]api.User, len(users))
	for i, v := range users {
		usersController[i] = BuildAPIUserFromControllerUser(&v)
	}

	return usersController
}

func BuildControllerUserFromAPIUser(user *api.User) controller.User {
	return controller.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Login:     user.Login,
		Password:  user.Password,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func BuildDBUserFromControllerUser(user *controller.User) db.User {
	return db.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Login:     user.Login,
		Password:  user.Password,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
