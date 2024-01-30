package monitor

import (
	"os/user"
)

type User struct {
	Id             string   `json:"id"`
	UserId         string   `json:"user_id"`
	Name           string   `json:"name"`
	Username       string   `json:"username"`
	PrimaryGroupId string   `json:"primary_group_id"`
	GroupIds       []string `json:"group_ids"`
	HomeDir        string   `json:"home_dir"`
}

func GetUserByUsername(username string) (*User, error) {
	o, err := user.Lookup(username)
	if err != nil {
		return nil, err
	}
	u := getUser(*o)
	return u, nil
}

func getUser(u user.User) *User {
	gids, _ := u.GroupIds()
	return &User{
		Id:             calculateUserId(u.Uid),
		UserId:         u.Uid,
		Name:           u.Name,
		Username:       u.Username,
		HomeDir:        u.HomeDir,
		PrimaryGroupId: u.Gid,
		GroupIds:       gids,
	}
}

func calculateUserId(uid string) string {
	return NewUUID5(GetHostId(), []byte(uid))
}
