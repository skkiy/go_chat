package main

import (
	"flag"
	"github.com/joeshaw/envdecode"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

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
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	var ts struct {
		ClientId    string `env:"GOOGLE_PROVIDER_CLIENT_ID, required"`
		ClientSecret string `env:"GOOGLE_PROVIDER_CLIENT_SECRET, required"`
	}
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}

	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		google.New(ts.ClientId, ts.ClientSecret, "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/logout", logOutHandler)
	http.HandleFunc("/auth/", loginHandler)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/room", r)

	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars"))))

	go r.run()
	log.Println("Webサーバーを開始します。ポート： ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func logOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}
