package utils

import (
	"blog/models"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Configuration struct {
	Address      string
	ReadTimeout  int64
	WriteTimeout int64
	Static       string
	Version      string
}

var Config Configuration
var Logger *log.Logger

func P(a ...interface{}) {
	fmt.Println(a)
}

//加载配置文件
func loadConfig() {
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		log.Fatalln("Connot open file", err)
	}
	decoder := json.NewDecoder(file)
	Config = Configuration{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
	fmt.Println()
}

//初始化
func init() {
	loadConfig()
	file, err := os.OpenFile("blog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	writers := []io.Writer{
		file,
		os.Stdout,
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	Logger = log.New(fileAndStdoutWriter, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error_message(w http.ResponseWriter, r *http.Request, msg string) {
	url := []string{"/err?msg=", msg}
	http.Redirect(w, r, strings.Join(url, ""), 302)
}

func Session(w http.ResponseWriter, r *http.Request) (tsess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil {
		sess := models.Session{Uuid: cookie.Value}
		var ok bool
		tsess, ok, _ = sess.Check()
		if !ok {
			err = errors.New("Invalid session")
			return
		}
	}
	return
}

//模板
func ParseTemplateFiles(filenames ...string) (t *template.Template) {
	var files []string
	t = template.New("layout")
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	t = template.Must(t.ParseFiles(files...))
	return
}

func GenerateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	_ = templates.ExecuteTemplate(w, "layout", data)
}

//日志
func Info(args ...interface{}) {
	Logger.SetPrefix("INFO")
	Logger.Println(args...)
}

func Danger(args ...interface{}) {
	Logger.SetPrefix("ERROR")
	Logger.Println(args...)
}

func Warning(args ...interface{}) {
	Logger.SetPrefix("WARNING")
	Logger.Println(args...)
}

//版本
func Version() string {
	return Config.Version
}
