package main

import (
	"Practica/Web"
	"Practica/api"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fs := http.FileServer(http.Dir("assets"))

	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", Web.HomeHandler) // каждый запрос вызывает функцию handler
	//mux.HandleFunc("/answer", AnswerHandler)

	go func() {
		var handlers api.Handler
		srv := new(api.Server)
		err := srv.Run("3000", handlers.InitRoutes())
		if err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	err := http.ListenAndServe(":"+port, mux) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*
func HomeHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("---------------------------")
	r.ParseForm() //анализ аргументов,
	fmt.Println(r.Form)  // ввод информации о форме на стороне сервера
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	//fmt.Fprintf(w, "Hello!") // отправляем данные на клиентскую сторону
	tpl.Execute(w, new(res))
}
func AnswerHandler(w http.ResponseWriter, r *http.Request){
	var mes res
	fmt.Println("---------------------------")
	r.ParseForm() //анализ аргументов

	mes.Password=r.Form["Password"][0]
	mes.Text=r.Form["Text"][0]

	fmt.Println(r.Form)  // ввод информации о форме на стороне сервера
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	if r.Form["Button"][0]=="Encrypt"{
		mes.Answer="Зашифрованный текст"
	}
	if r.Form["Button"][0]=="Decrypt"{
		mes.Answer="Расшифрованный текст"
	}
	tpl.Execute(w, mes)
	//fmt.Fprintf(w, "Hello!") // отправляем данные на клиентскую сторону
	//tpl.Execute(w, nil)
}
*/
