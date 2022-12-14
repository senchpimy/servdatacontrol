package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"strconv"
	"net/http"
	"strings"

	"os"
	"github.com/gorilla/context"
	"github.com/shirou/gopsutil/cpu"
)

var mainpage = "/"
var discospage = "/discosDisponibles"
var discosmontadospage = "/discos"
var sambapage = "/SambaConfi"
var UserConfig = "/UserConfig"
var Login = "/login"
var Logout = "/logout"
var System = "/System"
var Buttons = "/buttons"
var sambaGlobal = "/smbGlobal"
var nfspage = "/Nfs"
var dashboard = "/dashboard"

var tpl *template.Template

type dashboardInfo struct{
	CpuUsage []float64
	DisksUsage []int
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func errorHandler(w http.ResponseWriter, r *http.Request, PageName string) (foo bool){
	status:=http.StatusNotFound
	if r.URL.Path != PageName {
		w.WriteHeader(status)
		if status == http.StatusNotFound {
		tpl.ExecuteTemplate(w,"404.html",nil)
	        return true
		}
	}

	username:=getUserName(r)
	fmt.Println(username+"end")
	if len(username)==0{
	http.Redirect(w, r, Login, 302)
        return true
    }
	fmt.Println(PageName)
    return false

}

func readHtmlFromFile(fileName string) ([]byte) {
    bs, _ := ioutil.ReadFile(fileName)
    return bs
}

func INIT()  {
	if os.Geteuid() != 0 {
		fmt.Println("The program Needs to be run by root")
		os.Exit(0)
	}
	CreateParentDir()
	for _,filename := range []string{"/usr/local/gault/disks","/usr/local/gault/errorlog","/usr/local/gault/passwords"}{ //Agregar contraseña
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			if filename=="/usr/local/gault/passwords"{
				fmt.Println("No users/passwords file check the git for the default")
				os.Exit(0)
			}
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("Error creating file:", err)
				fmt.Printf("FIle: %s",filename )
				os.Exit(0)
			}
			defer file.Close()
	}}
	MountByFile()
	fmt.Println("INIT pasado")
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if errorHandler(w,r,mainpage) {
		return
	}
	tpl.ExecuteTemplate(w,"index.html",nil)
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func DiscosMontados(w http.ResponseWriter ,r *http.Request)  {
	if errorHandler(w,r,discosmontadospage){
		return
	}

	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "discos.html", FormaterDiskInfo(GetInfoSystem()))
	case "POST":
		fmt.Println("POST")
		if err := r.ParseForm(); err !=nil{
			fmt.Fprintf(w,"ParseForm() err: v%",err)
			return
		}

		diskUuid:=r.FormValue("diskselected")
		Umount(diskUuid)
		Data:=FormaterDiskInfo(GetInfoSystem())
		tpl.ExecuteTemplate(w, "discos.html", Data)

	default: fmt.Fprintf(w,"Error")
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func DiscosDisponibles(w http.ResponseWriter, r *http.Request)  {
	if errorHandler(w, r, discospage){
		return
    }
	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "discosDisponibles.html", GetDisks())
	case "POST":
		fmt.Println("POST")
		if err := r.ParseForm(); err !=nil{
			fmt.Fprintf(w,"ParseForm() err: v%",err)
			return
		}

		diskUuid:=r.FormValue("diskselected")
		VerifyDisk(diskUuid)
		tpl.ExecuteTemplate(w, "discosDisponibles.html", GetDisks())

	default: fmt.Fprintf(w,"Error")
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func SambaGlobal(w http.ResponseWriter, r *http.Request)  {
    if errorHandler(w,r,sambaGlobal){
	return
	}
	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "sambaGlobal.html", GetAllConfigurations())
	case "POST":
		fmt.Println("POST")
	default: fmt.Fprintf(w,"Error")
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func SambaConfiguration(w http.ResponseWriter, r *http.Request)  {
    if errorHandler(w,r,sambapage){
	return
    }
	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "samba.html", GetAllConfigurations())
	case "POST":
		if err := r.ParseForm(); err !=nil{
			fmt.Fprintf(w,"ParseForm() err: v%",err)
			return
		}

		if len(r.Form)!=1{
			Correct:=true
			var NewShare Share
			ConfigurationsReceived:=make([]Configuration,12)
			i:=0
			for key, element:= range r.Form{
				if key=="Titulo"{
					NewShare.Title=r.FormValue("Titulo")
					continue
				}
				if key=="Delete"{continue}
				if key=="valid users"{
					userslist, valid:=UsersExist(strings.Join(element," "))
					fmt.Println(valid)
					if valid==true{
						ConfigurationsReceived[i].Variable=key
						UsersFormated:=strings.Join(element," ")
						ConfigurationsReceived[i].Value=strings.Replace(UsersFormated,","," ",-1)
						i++
						continue
					}else{
						CreateError("The following users don exist:")
						for _, item := range userslist{
							CreateError(item)

						}
						CreateError("Imposible de Crear El Share")
						Correct=false
						break
					}
				}
				ConfigurationsReceived[i].Variable=key
				ConfigurationsReceived[i].Value=strings.Join(element," ")
				i++
			}
			NewShare.Contents=ConfigurationsReceived
			if Correct==true{VerifyShare(NewShare)}else{CreateError("Un Error Sucedio Imposible de Crear Share")}
		}else{
		Share:=r.FormValue("Delete")
		DeleteShare(Share)
		}
		tpl.ExecuteTemplate(w, "samba.html", GetAllConfigurations())

	default: fmt.Fprintf(w,"Error")
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func SystemOutput(w http.ResponseWriter, r *http.Request)  {
    if errorHandler(w,r,System){
	return
    }
		tpl.ExecuteTemplate(w, "status.html", SystemStatus())
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func Users(w http.ResponseWriter, r *http.Request)  {
    if errorHandler(w,r,UserConfig){
	return
    }
	switch r.Method {
	case "GET":
		tpl.ExecuteTemplate(w, "users.html", GetUsers())
	case "POST":
		fmt.Println("POST")
		if err := r.ParseForm(); err !=nil{
			fmt.Fprintf(w,"ParseForm() err: v%",err)
			return
		}

		User:=r.FormValue("User")
		Passw1:=r.FormValue("Passw1")
		Passw2:=r.FormValue("Passw2")
		TypeOfUser:=r.FormValue("Type")
		AddUser(User,Passw1,Passw2,TypeOfUser)
		tpl.ExecuteTemplate(w, "users.html", GetUsers())

	default: fmt.Fprintf(w,"Error")
	}
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func NfsPage(w http.ResponseWriter, r *http.Request)  {
    if errorHandler(w,r,nfspage){return}
    switch r.Method {
    case "GET":
    	tpl.ExecuteTemplate(w, "nfspage.html",ListExports())
    case "POST":
		nfs_path:=r.FormValue("Path")
		nfs_host:=r.FormValue("Host")
		nfs_permissions:=r.FormValue("Permissions")
		nfs_options:=r.FormValue("Options")
		nfs_delete:=r.FormValue("Delete")
		fmt.Println("nfs delete="+nfs_delete+"end")
		if nfs_path!=""{CreateExport(nfs_path, nfs_permissions, nfs_host,nfs_options);fmt.Println("Creado")}
		if nfs_delete!=""{DeleteNfs(nfs_delete);fmt.Println(nfs_delete);fmt.Println("Borrado")}
    		tpl.ExecuteTemplate(w, "nfspage.html",ListExports())
    }
}
///////////////////////////////////////////////////////////////////////////////////////////////////
func Dashboard(w http.ResponseWriter, r *http.Request)  {
	if errorHandler(w,r,dashboard){return}


	DasboardData:=dashboardInfo{}
	cpuUsage, err := cpu.Percent(0, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	DasboardData.CpuUsage=cpuUsage
	tmp:=FormaterDiskInfo(GetInfoSystem()).Todos
	tmp2:=make([]int,0)
	for _,i := range tmp{
		if len(i.UsePercent)!=0{
			value,_:=strconv.Atoi(i.UsePercent[:len(i.UsePercent)-1])
			tmp2=append(tmp2,value)
		}
	}
	DasboardData.DisksUsage=tmp2

	tpl.ExecuteTemplate( w,"dashboard.html",DasboardData)
	}
///////////////////////////////////////////////////////////////////////////////////////////////////

func main() {
	INIT()
	port := ":3000"
	fmt.Println("Startting on " + port)
	tpl, _ = template.ParseGlob("templates/*.html")
	http.HandleFunc(mainpage, indexHandler)
	http.HandleFunc(dashboard, Dashboard)
	http.HandleFunc(discosmontadospage, DiscosMontados)
	http.HandleFunc(discospage, DiscosDisponibles)
	http.HandleFunc(sambapage, SambaConfiguration)
	http.HandleFunc(sambaGlobal, SambaGlobal)
	http.HandleFunc(System, SystemOutput)
	//http.HandleFunc(ftpPage, FTPConfiguration)
	http.HandleFunc(UserConfig, Users)
	http.HandleFunc(nfspage, NfsPage)
	http.HandleFunc(Login, login)
	http.HandleFunc("/styles/style.css", func(response http.ResponseWriter, request *http.Request) {
	http.ServeFile(response, request, "templates/styles/style.css")})
	http.HandleFunc(Logout, logout)
	http.HandleFunc(Buttons, HandleButtons)
	http.ListenAndServe(port, context.ClearHandler(http.DefaultServeMux))

}
