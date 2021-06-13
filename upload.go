package main

import (
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userid")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := filepath.Join("avatars", userId+filepath.Ext(header.Filename))
	var permAllUsersCanAccess = 0777
	err = ioutil.WriteFile(filename, data, fs.FileMode(permAllUsersCanAccess))
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "success")
}
