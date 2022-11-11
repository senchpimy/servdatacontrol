package main


import (
//	"crypto/sha512"
//	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"github.com/gorilla/sessions"
	// "golang.org/x/crypto/bcrypt"
)

var storeOfsessions = sessions.NewCookieStore([]byte("mysession"))
type User struct{
User_name string
Password string
}

func login(w http.ResponseWriter, r *http.Request) {
	status:=http.StatusNotFound
	if r.URL.Path != Login {
	//if r.URL.Path != Login || r.URL.Path != Logout{
		w.WriteHeader(status)
		if status == http.StatusNotFound {
			w.Write(readHtmlFromFile("./404.html"))
	        return 
		}
	}
	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "login.html", nil)
	case "POST":
		if err := r.ParseForm(); err !=nil{
			fmt.Fprintf(w,"ParseForm() err: v%",err)
			return
		}

		//fmt.Fprintf(w,"Post form website r.postfrom =%v \n",r.PostForm)
		name:=r.FormValue("name")
		pass:=r.FormValue("pass")
		if GetPasswordConfirmation(name,pass){
			session, _ := storeOfsessions.Get(r, "session")
			session.Values["authenticated"] = true
			session.Save(r, w)
			fmt.Println("correct")
			//http.Redirect(w, r, "/", http.StatusFound)
			tpl.ExecuteTemplate(w, "index.html", nil)
		}else{
			fmt.Println("No correct")
		}
//		sha_512:=sha512.New()
//		sha_512.Write([]byte(name))
//		fmt.Fprintf(w,"value = %s\n",name)
//		fmt.Fprintf(w,"value in sha_512 = \t%s",base64.StdEncoding.EncodeToString(sha_512.Sum(nil)))
//		if GetPasswordConfirmation(name,pass){
//			session, _ := store.Get(r, "session")
//			session.Values["userID"] = name
//			session.Save(r, w)
//			http.Redirect(w, r, "/", http.StatusFound)
//		}
//		}else{
//			http.Redirect(w, r, "/failed", http.StatusFound)
//		}


	default: fmt.Fprintf(w,"Error")
	}
}

func GetPasswordConfirmation(User string,Password string)(foo bool){
	UsersInGault:= GetGaultUsers()
	exist_user:=false
	 for _,item:= range UsersInGault{
		 if item.User_name==User{
			exist_user=true
			//password:=sha512.New()
			//password.Write([]byte(Password))
			//if item.Password==base64.StdEncoding.EncodeToString(password.Sum(nil)){
			// err := bcrypt.CompareHashAndPassword([]byte(item.Password), []byte(Password))
			// if err==nil{
			if item.Password==Password{
				fmt.Println("LoggedSuccedful")
				return true
			}else{
				fmt.Println("Passwords Dont match");
				CreateError("Passwords Dont match");
	 			return false
			}
		}
	}
	if exist_user==false{
	 fmt.Println("User dont found")
	 CreateError("User dont found")
	 return false
 	}
	return
}

func GetGaultUsers()(foo []User){  //completo
    file, err := os.Open("./passwords")
    if err != nil {
        log.Fatal(err)
    }
    data, err := ioutil.ReadAll(file)
    if err != nil {
        log.Fatal(err)
    }
    data_string:=string(data)

    data_split:=strings.Split(data_string,"\n")
    Users:=make([]User,len(data_split)-1)
    for index,line:= range data_split[:len(data_split)-1]{
	    data_from_line:=strings.Split(line,":")
	    Users[index].Password=data_from_line[1]
	    Users[index].User_name=data_from_line[0]

    }
    return Users
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****logoutHandler running*****")
	session, _ := storeOfsessions.Get(r, "session")
	session.Values["authenticated"] = false
	session.Save(r, w)
	fmt.Println("Redirigiendo plantilla")
}

func logout(w http.ResponseWriter, r *http.Request) {
	errorHandler(w,r,Logout)
	tpl.ExecuteTemplate(w, "logout.html", nil)
}
