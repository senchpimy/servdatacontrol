package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/go-ini/ini"
)

//	GLOBAL
//host allow: Solo computadoras en la red local
//workgourp:
//server string:
//guest account: cuanta de invitado
//max log sizie: tamaño en kb de los archivos de log
//passdb backend:
//interfaces: configuracion de interfaces
//wins support: soporte para windows

//	SHARE DEFINITIONS
// #######todo##########
//hosts allow = 192.168.0.0/16
//hosts deny= 0.0.0.0/0

type Configuration struct{
Variable string
Value string
}

type Share struct{
Title string
Contents []Configuration
}

type ConfigurationsStruct struct{
Sections []Share
Mounted []Disk_DF
}

func WriteToFile(Texto string, File string, location int) {
	// Read Write Mode
	file, err := os.OpenFile(File, os.O_RDWR, 0644)
	
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()
	
	len, err := file.WriteString(Texto)
	if err != nil {
		log.Fatalf("failed writing to file: %s %s", err,len)
	}
}

func ExistSambaConf()  { //completado
	if _, err := os.Stat("/etc/samba/smb.conf"); errors.Is(err, os.ErrNotExist) {
		log.Fatal("Samba Config File 'smb.conf' does not exist")
		}
}

func CreateConfiguration(Configuration Share)(foo []string){ //Completado
	title:="\n[" +Configuration.Title +"]\n"
	elementsLen:=len(Configuration.Contents)+1
	s:=make([]string,elementsLen)
	s[0]=title
	for i:=0;i<elementsLen-1;i++{
	s[i+1]=Configuration.Contents[i].Variable+" = "+Configuration.Contents[i].Value+"\n"
	}

	return s

}

func WriteShareConf(bar []string){ //Completado
	f, err := os.OpenFile("./smb.conf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
	    log.Fatal(err)
	}
	defer f.Close()
	    for i:=0;i<len(bar);i++{
	if _, err := f.Write([]byte(bar[i])); err != nil {
	    log.Fatal(err)
	    }
	}
	if err := f.Close(); err != nil {
	    log.Fatal(err)
	}
}

func DeleteShare(share string)(error error){
	var start int = -1
	var end int = -1

	 //data, err := ioutil.ReadFile("/etc/samba/smb.conf")
	 data, err := ioutil.ReadFile("./smb.conf")
	if err != nil {
	    log.Fatal(err)
	}
	
	file := string(data)
	temp := strings.Split(file, "\n")
	
	 for index, item := range temp {
	         if strings.Contains(item, share) && strings.Contains(item,"#")==false{start=index;break}
	}
	if start==-1{return errors.New("No Share Named "+share+" Found")}
	for index, item := range temp {
	        if strings.Contains(item, "[",) && index>start{end=index;break}
	   }
        //newFile, err := os.Create("/etc/samba/smb.conf.after")
        newFile, err := os.Create("./smb.conf.after")
	if err != nil {
        log.Fatal(err)
    	}
    	defer newFile.Close()
	
	if end==-1{
		for index, item:=range temp{
			if index<start {
	        	   _, err := newFile.Write([]byte(item+"\n"))
			    if err != nil {
			        log.Fatal(err)
			    	}
			}
		}

	}else{
		for index, item:=range temp{
			if index<start || index>end-1{
	        	   _, err := newFile.Write([]byte(item+"\n"))
			    if err != nil {
			        log.Fatal(err)
			    	}
			}
		}
	}


	//cmd,err:=exec.Command( "mv","/etc/samba/smb.conf","/etc/samba/smb.conf.bak").Output()
	cmd,err:=exec.Command( "mv","./smb.conf","./smb.conf.bak").Output()
	if err != nil {log.Fatal(err)}
	fmt.Println(cmd)
	//cmd2,err:=exec.Command( "mv","/etc/samba/smb.conf.after","/etc/samba/smb.conf").Output()
	cmd2,err:=exec.Command( "mv","./smb.conf.after","./smb.conf").Output()
	if err != nil {log.Fatal(err)}
	fmt.Println(cmd2)
return nil
}

func GetAllConfigurations()(foo ConfigurationsStruct){
	file,err:=ini.Load("./smb.conf")
	 if err != nil {
        fmt.Printf("Fail to read file: %v", err)
        os.Exit(1)
    }
    Configurations:=make([]Share,len(file.SectionStrings()))

    for i,section:=range file.SectionStrings(){
	Configurations[i].Title=section
    	Variables:=make([]Configuration,len(file.Section(section).KeyStrings()))
	for j,key:=range file.Section(section).KeyStrings(){
		Variables[j].Variable=key
		Variables[j].Value=file.Section(section).Key(key).String()
}
	Configurations[i].Contents=Variables
}
var Configs ConfigurationsStruct
Configs.Sections=Configurations
Data:=FormaterDiskInfo(GetInfoSystem())
Configs.Mounted=Data.Todos
    return Configs
}

 func StartSamba(){
	
    cmd := exec.Command("systemctl", "start", "smb")
    err := cmd.Run()
    if err != nil {log.Fatal(err)}
    cmd2 := exec.Command("systemctl", "start", "nmb")
    err2 := cmd2.Run()
    if err2 != nil {log.Fatal(err2)}
 }

 func EnableSamba(){
	
    cmd := exec.Command("systemctl", "enable", "smb")
    err := cmd.Run()
    if err != nil {log.Fatal(err)}
    cmd2 := exec.Command("systemctl", "enable", "smb")
    err2 := cmd2.Run()
    if err2 != nil {log.Fatal(err2)}
 }

 func VerifyShare(ReceivedShare Share)  {
	NewContents:=ReceivedShare.Contents
	//NewContents2:=NewContents
	elementes:=[]string{"comment", "guest only", "writable", "valid users", "guest ok", "browseable", "hosts deny", "read only", "wins support", "hosts allow", "LocalsOnly","public"}
	elementes1:=make([]string,len(elementes))
		for index,item:= range ReceivedShare.Contents{
			if item.Variable==""{continue}
			if item.Value=="on"{NewContents[index].Value="yes"}
			elementes1[index]=item.Variable
		}
	ReceivedShare.Contents=NewContents
	noselected:=difference(elementes,elementes1)
	for index,item:=range noselected{
		ReceivedShare.Contents[len(ReceivedShare.Contents)-index-1].Variable=item
		ReceivedShare.Contents[len(ReceivedShare.Contents)-index-1].Value="no"
	}

	WriteShareConf(CreateConfiguration(ReceivedShare))
 }

func AddSambaUser(user string, password string)  {
	err:=exec.Command("sh","./CreateUser.sh",user, password).Run()
    	if err != nil {CreateError(err.Error())}
}

func difference(a, b []string) []string {
    mb := make(map[string]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []string
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
}
