package gitchecker

import (
	"errors"
	"fmt"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"os"
)

//현재는 레포 업데이트 알고리즘이 지우고 다시 받는 알고리즘임.
func Gitupdate(url string,directory string,id string,password string)(error){
	err :=os.RemoveAll(directory)
	if (err!=nil){
		return err
	}
	_, err = AuthGitclone(url, directory,id,password)
	if (err!=nil){
		return err
	}
	log.Println("Successfully change commit")
	return nil
}

func Isrepotobeupdate(url string,directory string)(bool,error){
	if!(fileExists(directory)){
		log.Println("There is no file in directory")
		return true,errors.New("NOFILEEXIST")
	} else{
		r,err:=git.PlainOpen(directory)
		if (err!=nil){return false,err}
		err =r.Fetch(&git.FetchOptions{})
		if(err!=nil){
			if(err == git.NoErrAlreadyUpToDate) {
				log.Println("There is no update in remote git repository")
				return false, nil
			} else {
				return false,err
			}
		}
		log.Println("There is some change in remote repository")
		return true,nil

	}
}

func AuthGitclone(url string, directory string,id string,password string) (*object.Commit,error){
	log.Printf("git clone %s %s --recursive", url, directory)
	authurl:= fmt.Sprintf("%s:%s@%s",id,password,url)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               authurl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,

	})
	ref, err := r.Head()
	if err!=nil {
		return nil,err
	}
	commit, err := r.CommitObject(ref.Hash())
	if err!=nil {
		return nil,err
	}
	return commit,err
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}