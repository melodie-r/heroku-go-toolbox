package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
	// w.Write([]byte("こんにちは！"))
	t := template.Must(template.ParseFiles("template/index.gohtml", "template/_menu.gohtml"))
	if err := t.Execute(w, struct {
		UserName string
		Time     time.Time
	}{
		"ゲスト",
		time.Now(),
	}); err != nil {
		log.Printf("テンプレート %s の実行に失敗！: %v", t.Name(), err)
		http.Error(w, "内部エラーです", http.StatusInternalServerError)
	}
}

func handleSecret(w http.ResponseWriter, r *http.Request) {
	user, password, _ := r.BasicAuth()
	correct_user := os.Getenv("USER")
	correct_pwd := os.Getenv("PASSWORD")
	if user != correct_user || password != correct_pwd {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "認証に失敗しました", http.StatusUnauthorized)
		return
	}
	log.Printf("%s %s", r.Method, r.RequestURI)
	w.Write([]byte("秘密のページです！"))
}

func handleHolidays(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
	t := template.Must(template.ParseFiles("template/holidays.gohtml", "template/_menu.gohtml"))
	if err := t.Execute(w, nil); err != nil {
		log.Printf("テンプレート %s の実行に失敗！: %v", t.Name(), err)
		http.Error(w, "内部エラーです", http.StatusInternalServerError)
	}
}

func main() {
	port := os.Getenv("PORT") // 実行時に Heroku が指定するポート番号を取得
	if len(port) == 0 {
		port = "8080" // ローカルで実行するときのポート番号を指定
	}
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/secret", handleSecret)
	http.HandleFunc("/holidays", handleHolidays)
	log.Printf("ポート %s で待ち受けを開始します...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("サーバーが異常終了しました: %v", err)
	}
}
