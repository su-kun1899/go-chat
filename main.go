package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/su-kun1899/go-chat/trace"
)

var avatars Avatar = UseFileSystemAvatar

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	var port = flag.String("port", "8080", "アプリケーションのアドレス")
	var securityKey = flag.String("security_key", "指定必須", "セキュリティキー")
	var googleClientID = flag.String("google_client_id", "指定必須", "GoogleのクライアントID")
	var googleSecret = flag.String("google_secret", "指定必須", "Googleのクライアント鍵")
	flag.Parse()
	// Gomniauthのセットアップ
	// TODO securityKeyも外から渡せるように
	gomniauth.SetSecurityKey(*securityKey)
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の値", "http://localhost:"+*port+"/auth/callback/facebook"),
		github.New("クライアントID", "秘密の値", "http://localhost:"+*port+"/auth/callback/github"),
		google.New(*googleClientID, *googleSecret, "http://localhost:"+*port+"/auth/callback/google"),
	)
	//r := newRoom(UseAuthAvatar)
	//r := newRoom(UseGravatar)
	r := newRoom(UseFileSystemAvatar)
	r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)

	//チャットルームを開始
	go r.run()

	// Starting web server
	log.Println("Webサーバを開始します。ポート: ", *port)
	if err := http.ListenAndServe((":" + *port), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
