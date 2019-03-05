package controller

import (
	"blog/models"
	"blog/utils"
	"net/http"
)

//注册账户
func SignupAccount(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.Danger(err, "Cannot parse form")

	}
	user := models.User{
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	if err := user.Creat(); err != nil {
		utils.Danger(err, "Cannot create user")
		utils.Error_message(w, r, "邮箱已被注册")
		return
	}
	sess, err := user.CreateSession()
	if err != nil {
		utils.Danger(err, "Can't Create Session")
	}
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    sess.Uuid,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 302)
}

//登录
func Authenticate(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	user, err := models.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		utils.Danger(err, "Cannot find user")
	}
	password := models.Encrypt(r.PostFormValue("password"))
	if user.Password == password {
		session, err := user.CreateSession()
		if err != nil {
			utils.Danger(err, "Cannot create session")
		}
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect(w, r, "/err?msg=密码错误", 302)
	}

}

//登出
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != http.ErrNoCookie {
		session := models.Session{Uuid: cookie.Value}
		session.DeleteByUUID()
	} else {
		utils.Warning(err, "Failed to get cookie")
	}
	http.Redirect(w, r, "/", 302)

}
