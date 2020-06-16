package main

import (
	"fmt"

	"net/http"

	"path/filepath"
	"strings"

	//"html/template"

	"github.com/julienschmidt/httprouter"
)

func main() {
	items = make([]item, 0)
	activeRouter := httprouter.New()
	staticRouter := httprouter.New()

	activeRouter.GET("/", getList)
	activeRouter.GET("/:id", getItem)
	activeRouter.POST("/:id", updateItem)
	activeRouter.GET("/newItem", newItem)

	staticRouter.ServeFiles("/*filepath", neuteredFileSystem{http.Dir("./frontend")})
	activeRouter.NotFound = staticRouter

	server := &http.Server{
		Addr:    ":17181",
		Handler: activeRouter,
	}

	//In production, to be switched with the tlc service
	err := server.ListenAndServe()
	if err != nil {
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
