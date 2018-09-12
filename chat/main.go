package main

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "chat.html")
	//フラグを解釈する
	flag.Parse()
	gomniauth.SetSecurityKey("aqwsedrftgyhujikol")
	gomniauth.WithProviders(
		facebook.New("592533754731-q9fv2ie1l2et9dqkk9d59htoq8fhjkrm.apps.googleusercontent.com", "OUx2g-PSg5j0bxNBA4DlTZfn", "http://localhost:8080/auth/callback/facebook"),
		github.New("592533754731-q9fv2ie1l2et9dqkk9d59htoq8fhjkrm.apps.googleusercontent.com", "OUx2g-PSg5j0bxNBA4DlTZfn", "http://localhost:8080/auth/callback/github"),
		google.New("592533754731-q9fv2ie1l2et9dqkk9d59htoq8fhjkrm.apps.googleusercontent.com", "OUx2g-PSg5j0bxNBA4DlTZfn", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	//なんでnewRoom内でtracerを作成しないの？デフォルトのtracerを設定するため
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	//チャットルームの開始
	go r.run()
	//webサーバを起動
	//addrはポインタなので、デリファレンスして実際の値を取得する必要がある
	log.Println("Webサーバーを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	//if err := http.ListenAndServe(":8080", nil); err != nil {
	//	log.Fatal("ListenAndServe:", err)
	//}
}
