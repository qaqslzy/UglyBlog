package controller

import (
	"blog/models"
	"blog/utils"
	"fmt"
	"net/http"
	"strconv"
)

//新建文章的页面
//GET
func NewEntry(w http.ResponseWriter, r *http.Request) {
	data := struct {
		User  models.User
		Entry models.Entry
	}{}
	query := r.URL.Query()
	eid := query.Get("id")
	data.Entry, _ = models.ReadEntry(eid)
	sess, err := utils.Session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		data.User, _ = sess.User()
		utils.GenerateHTML(w, data, "layout", "private.navbar", "edit")
	}
}

// 删除文章
func DeleteEntry(w http.ResponseWriter, r *http.Request) {
	sess, err := utils.Session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		user, err := sess.User()
		if err != nil {
			utils.Danger("Session Wrong")
			http.Redirect(w, r, "/", 302)
			return
		}
		query := r.URL.Query()
		eid := query.Get("id")
		err = user.DeleteEntry(eid)
		if err != nil {
			utils.Danger("User or Delete Wrong")
			http.Redirect(w, r, "/", 302)
			return
		}
		http.Redirect(w, r, "/", 302)
	}
}

//TODO 查找某一用户写过的文章
func UserEntry(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := query.Get("page")
	id := query.Get("id")
	if id == "" {
		http.Redirect(w, r, "/", 302)
		return
	}
	data := struct {
		Entries  []models.Entry
		User     models.User
		Page     int
		MaxPage  int
		HasNext  bool
		HasFont  bool
		NextPage int
		FontPage int
		Search   string
	}{Search: id}
	num, _ := models.UserEntriesCount(id)
	if num != 0 {
		if num%models.PageCount == 0 {
			data.MaxPage = num / models.PageCount
		} else {
			data.MaxPage = num/models.PageCount + 1
		}
	} else {
		data.MaxPage = 1
		data.Page = 1
	}

	if page != "" {
		var err error
		data.Page, err = strconv.Atoi(page)
		if err != nil {
			utils.Error_message(w, r, "Cannot get entries")
			return
		}
		if data.Page > data.MaxPage {
			data.Page = data.MaxPage
		}

		if data.Page < 1 {
			data.Page = 1
			data.HasFont = false
		} else if data.Page == 1 {
			data.HasFont = false
		} else {
			data.HasFont = true
			data.FontPage = data.Page - 1
		}
		entries, err := models.UserEntriesPage(id, data.Page)
		data.Entries = entries

	} else {
		entries, err := models.GetUserEntry(id)
		data.Entries = entries
		data.Page = 1
		data.HasFont = false
		if err != nil {
			utils.Error_message(w, r, "Cannot get entries")
			return
		}
	}
	if data.Page >= data.MaxPage {
		data.HasNext = false
	} else {
		data.HasNext = true
		data.NextPage = data.Page + 1
	}
	sess, err := utils.Session(w, r)
	if err != nil {
		utils.GenerateHTML(w, data, "layout", "public.navbar", "userindex")
	} else {
		user, err := sess.User()
		data.User = user
		if err != nil {
			utils.GenerateHTML(w, data, "layout", "public.navbar", "userindex")
		} else {
			utils.GenerateHTML(w, data, "layout", "private.navbar", "userindex")
		}
	}
}

//新建文章 和 编辑
//POST
func CreatEntry(w http.ResponseWriter, r *http.Request) {
	sess, err := utils.Session(w, r)
	if err != nil {
		http.Redirect(w, r, "/", 302)
	} else {
		err = r.ParseForm()
		if err != nil {
			utils.Danger(err, "Cannot parse form")
		}
		user, err := sess.User()
		if err != nil {
			utils.Danger(err, "Cannot get user from session")
		}
		body := r.PostFormValue("body")
		topic := r.PostFormValue("topic")
		abstract := r.PostFormValue("abstract")
		title := r.PostFormValue("title")
		uuid := r.PostFormValue("uuid")
		var url string
		if uuid != "" {
			err := user.UpdateEntry(uuid, title, topic, body, abstract)
			if err != nil {
				http.Redirect(w, r, "/err?msg=你想干嘛", 302)
				return
			}
			url = fmt.Sprintf("/entry/read?id=%s", uuid)
		} else {
			entry, err := user.CreateEntry(title, topic, body, abstract)
			if err != nil {
				utils.Danger(err, "Cannot creat post")
			}
			url = fmt.Sprintf("/entry/read?id=%s", entry.Uuid)
		}

		http.Redirect(w, r, url, 302)
	}
}

func ReadEntry(w http.ResponseWriter, r *http.Request) {
	sess, usererr := utils.Session(w, r)
	data := struct {
		User     models.User
		Entry    models.Entry
		IsMaster bool
	}{IsMaster: false}
	query := r.URL.Query()
	eid := query.Get("id")
	var e error
	data.Entry, e = models.ReadEntry(eid)
	if e != nil {
		http.Redirect(w, r, "/err?msg=不存在该文章", 302)
	} else {
		if usererr != nil {
			utils.GenerateHTML(w, data, "layout", "public.navbar", "entry")
		} else {
			data.User, _ = sess.User()
			if data.Entry.UserId == data.User.Uuid {
				data.IsMaster = true
			}
			utils.GenerateHTML(w, data, "layout", "private.navbar", "entry")
		}
	}

}
