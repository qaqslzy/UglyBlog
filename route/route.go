package route

import (
	"blog/controller"
	"blog/utils"
	"net/http"
)

var mux *http.ServeMux

func init() {

	mux = http.NewServeMux()
	files := http.FileServer(http.Dir(utils.Config.Static))
	//静态资源
	mux.Handle("/static/", http.StripPrefix("/static/", files))
	//index
	mux.HandleFunc("/", controller.Index)
	// err
	mux.HandleFunc("/err", controller.Err)
	// 登录
	mux.HandleFunc("/login", controller.Authenticate)
	// 登出
	mux.HandleFunc("/logout", controller.Logout)
	// 新建文章
	mux.HandleFunc("/newentry", controller.NewEntry)
	// 提交文章
	mux.HandleFunc("/postentry", controller.CreatEntry)
	// 注册
	mux.HandleFunc("/signup", controller.SignupAccount)
	// 文章内容
	mux.HandleFunc("/entry/read", controller.ReadEntry)
	// 删除文章
	mux.HandleFunc("/delete/entry", controller.DeleteEntry)
	// 友情链接
	mux.HandleFunc("/friend/", controller.FriendLink)
	// 关于本人
	mux.HandleFunc("/aboutme/", controller.AboutMe)
	// 查看某一个人的
	mux.HandleFunc("/user", controller.UserEntry)
}

func ReturnMux() *http.ServeMux {
	return mux
}
