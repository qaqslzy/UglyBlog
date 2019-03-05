package main

import (
	"blog/route"
	"blog/utils"
	"net/http"
	"time"
)

/**
*
* @author Liu Weiyi
* @date 2019-03-04 10:08
 */

func main() {
	utils.P("uglyBlog", utils.Config.Version, "started at", utils.Config.Address)
	//启动！
	server := &http.Server{
		Addr:           utils.Config.Address,
		Handler:        route.ReturnMux(),
		ReadTimeout:    time.Duration(utils.Config.ReadTimeout * int64(time.Second)),
		WriteTimeout:   time.Duration(utils.Config.WriteTimeout * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}

	_ = server.ListenAndServe()
}
