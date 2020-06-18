package main

import (
	"fmt"
	"html/template"
	"os"

	"net/http"

	"path/filepath"
	"strings"

	//"html/template"

	"encoding/json"
	"io/ioutil"

	"github.com/julienschmidt/httprouter"
)

var itemPageTemplate *template.Template

func main() {
	//For somereason itempageTemplate was becoming reinstatiated inside main, so not using the var we set up
	var err error
	itemPageTemplate, err = template.New("").ParseFiles("./static/itemPage.html", "./static/listPage.html")
	if err != nil{
		panic(err)
	}
	items = make([]*item, 0)
	loadItemsFromBackup()
	fmt.Println(items)
	count = len(items)

	//Create the images folder if it does not exist
	_, err = os.Stat("./static/images")
	if os.IsNotExist(err) {
		errDir := os.MkdirAll("./static/images", 0755)
		if errDir != nil {
			panic(err)
		}
	}

	//it := item{Title: "Temp"}
	//items = append(items, it)
	activeRouter := httprouter.New()
	staticRouter := httprouter.New()

	activeRouter.GET("/", getList)
	activeRouter.GET("/newItem", newItem)
	activeRouter.POST("/upload", uploadFile)
	activeRouter.GET("/item/:id", getItem)
	activeRouter.POST("/item/:id", updateItem)
	

	staticRouter.ServeFiles("/*filepath", neuteredFileSystem{http.Dir("./static")})
	activeRouter.NotFound = staticRouter

	server := &http.Server{
		Addr:    ":17181",
		Handler: activeRouter,
	}

	//In production, to be switched with the tlc service
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func saveItemsToBackup(){
	file, _ := json.Marshal(items)
	_ = ioutil.WriteFile("backup.json", file, 0644)
}

func loadItemsFromBackup(){
	b, err := ioutil.ReadFile("backup.json")

	//Change, but we should check if teh file exists, if not, then we create it
	if err != nil{
		fmt.Println(err)
		ioutil.WriteFile("backup.json", []byte{}, 0644)
		return
	}

	err = json.Unmarshal(b, &items)
	fmt.Println(err)
}

func pError(err error){
	if err != nil{
		panic(err)
	}
}

//From https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
//There is a copy of this also at https://medium.com/@hau12a1/golang-http-serve-static-files-correctly-5feb98ae9da1
type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	//Modified code, no need to have .html at the end of our html static filepaths
	//If our filepath is not a directory and is looking for a file
	if filepath.Base(path) != "\\" && filepath.Base(path) != "" {
		if filepath.Ext(path) == "" {
			//Append the file path with a .html
			path = fmt.Sprintf("%s.html", path)
		}
	}

	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
