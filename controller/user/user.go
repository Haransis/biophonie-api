package user

import "github.com/golang-jwt/jwt/v4"

type AddUser struct {
	Name string `json:"name" example:"bob" binding:"required,min=3,max=20"`
}

type AccessToken struct {
	Token string `json:"token"`
}

type AuthUser struct {
	Name     string `json:"name" example:"bob" binding:"required,min=3,max=20"`
	Password string `json:"password" example:"9b768967-d491-4baa-a812-24ea8a9c274d" binding:"required,uuid"`
}

type User struct {
	Id        int    `db:"id" json:"userId" example:"123" validate:"gte=0"`
	Name      string `db:"name" json:"name" example:"bob" binding:"required"`
	Password  string `db:"password" json:"password" example:"9b768967-d491-4baa-a812-24ea8a9c274d"`
	CreatedOn string `db:"created_on" json:"createdOn" example:"2022-05-26T11:17:35.079344Z"`
	Admin     bool   `db:"admin" json:"admin" example:"false"`
}

// Define some custom types were going to use within our tokens
type UserInfo struct {
	Name  string
	Admin bool
}

type CustomClaims struct {
	*jwt.RegisteredClaims
	UserInfo
}
