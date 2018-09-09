package main

import (
	"flag"
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
	r := newRoom()
	//なんでnewRoom内でtracerを作成しないの？デフォルトのtracerを設定するため
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
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
