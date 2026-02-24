package handler

import (
	"net/http"
	"fmt"
	"github.com/Kaikai20040827/graduation/internal/pkg"
	"github.com/Kaikai20040827/graduation/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSrv *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	fmt.Println("✓ Creating a new user handler done")
	return &UserHandler{
		userSrv: us,
	}
}

func (uh *UserHandler) GetProfile(context *gin.Context) {
	uidv, ok := context.Get("user_id")
	if !ok {
		pkg.JSONError(context, 401, "unauthorized")
		context.Abort()
		return
	}
	uid := uidv.(uint)

	user, err := uh.userSrv.GetByID(uid)
	if err != nil {
		pkg.JSONError(context, 404, "cannot find user")
		context.Abort()
		return
	}
	pkg.JSONOK(context, user)
}

type UpdateProfileReq struct {
	Username string `json:"username" binding:"required"`
}

func (uh *UserHandler) UpdateProfile(context *gin.Context) {
	var req UpdateProfileReq
	if err := context.ShouldBindBodyWithJSON(&req); err != nil {
		pkg.JSONError(context, 40001, "invalid params")
		return
	}

	uidv, _ := context.Get("user_id")
	uid := uidv.(uint)
	u, err := uh.userSrv.UpdateProfile(uid, req.Username)
	if err != nil {
		pkg.JSONError(context, 50001, err.Error())
		return
	}
	pkg.JSONOK(context, u)
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required, min=6"`
}

func (uh *UserHandler) ChangePassword(context *gin.Context) {
	var req ChangePasswordReq
	if err := context.ShouldBindBodyWithJSON(&req); err != nil {
		pkg.JSONError(context, 40001, "invalid params")
		return
	}

	emailv, _ := context.Get("email")
	email := emailv.(string)

	if err := uh.userSrv.ChangePassword(email, req.OldPassword, req.NewPassword); err != nil {
		pkg.JSONError(context, 40002, err.Error())
		return
	}
	context.Status(http.StatusNoContent)
}

type DeleteUserReq struct {
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func (uh *UserHandler) DeleteUser(context *gin.Context) {
	var req DeleteUserReq
	if err := context.ShouldBindBodyWithJSON(&req); err != nil {
		pkg.JSONError(context, 40001, "invalid params")
		return
	}

	passwordv, _ := context.Get("password")
	password := passwordv.(string)

	emailv, _ := context.Get("email")
	email := emailv.(string)

	if err := uh.userSrv.DeleteUser(email, password); err != nil {
		pkg.JSONError(context, 40002, err.Error())
		return
	}
	context.Status(http.StatusNoContent)//后续开发跳转网页
}
