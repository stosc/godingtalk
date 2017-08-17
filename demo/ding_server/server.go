package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/hugozhu/godingtalk"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	CorpId     string `yaml:"corpId"`
	CorpSecret string `yaml:"corpSecret"`
	AgentID    string `yaml:"agentID"`
}

func (c *conf) getConf() *conf {

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func GetTemplate(tpl string) *template.Template {
	t, _ := template.ParseFiles(tpl)
	return t
}

var client *godingtalk.DingTalkClient

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	client.RefreshAccessToken()
	info, _ := client.UserInfoByCode(code)

	json.NewEncoder(w).Encode(info)
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	chatid := r.FormValue("cid")
	sender := r.FormValue("sender")
	content := r.FormValue("content")

	var resp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	client.RefreshAccessToken()
	err := client.SendTextMessage(sender, chatid, content)
	if err != nil {
		resp.ErrCode = -1
		resp.ErrMsg = err.Error()
	}

	json.NewEncoder(w).Encode(resp)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := path.Join("templates", "layout.html")
	fp := path.Join("templates", "index.html")

	url := "http://" + r.Host + r.RequestURI
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	client.RefreshAccessToken()
	configString := client.GetConfig("abcdabc", timestamp, url)

	data := make(map[string]interface{})
	data["config"] = template.JS(configString)

	tmpl, err := template.ParseFiles(lp, fp)
	if err == nil {
		err = tmpl.ExecuteTemplate(w, "layout", data)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var c conf
	c.getConf()
	corpId := c.CorpId
	corpSecret := c.CorpSecret
	client = godingtalk.NewDingTalkClient(corpId, corpSecret)
	client.AgentID = c.AgentID

	fs := http.FileServer(http.Dir(path.Join(os.Getenv("root"), "public")))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/get_user_info", getUserInfo)
	http.HandleFunc("/send_message", sendMessage)

	http.ListenAndServe(":9091", nil)
}
