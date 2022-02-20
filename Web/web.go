package Web

import (
	"Pracric/AnsibleVault"
	"fmt"
	"net/http"
	"strings"
	"text/template"
)

var tpl = template.Must(template.ParseFiles("Web/index.html")) // указывается шаблон

type res struct {
	Password string
	Text string
	Answer string
} // структура для вставки в шаблон

func HomeHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("---------------------------")
	var mes res
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Ошибка анализа аргументов")
	} //анализ аргументов
	fmt.Println(r.Form)  // ввод информации о форме на стороне сервера
	//fmt.Println("path", r.URL.Path) // Показывает текущий адрес
	//fmt.Println("scheme", r.URL.Scheme) // Показывает схему ???
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	if len(r.Form["Button"])==0{
	}else{

		mes.Password=r.Form["Password"][0]
		mes.Text=r.Form["Text"][0]

		if r.Form["Button"][0] == "Encrypt" {

			mes.Answer,err = AnsibleVault.Encrypt(mes.Text,mes.Password)
			fmt.Println(mes.Answer)
			if err!=nil{
				fmt.Printf("error: %v", err)
			}
		}
		if r.Form["Button"][0] == "Decrypt" {

			mes.Answer,err = AnsibleVault.Decrypt(mes.Text,mes.Password)
			if err!=nil{
				//fmt.Printf("error: %v", err)
				mes.Answer = "Error of decryption"
			}
		}
	}
	tpl.Execute(w, mes)
}