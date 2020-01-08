package gitchecker

import (
	git "gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"os"
)


func Gitupdate(url string,directory string)error{
	flag,err:=Isrepotobeupdate(url,directory)
	if(err!=nil){return err}
	if(flag==false){
		log.Println("Repo is not to be update")
		return nil
	}
	log.Println("Repo is has to be update")
	_, err = gitclone(url, "new")
	if (err!=nil){
		return err
	}
	log.Println("Download new appoconfig commit")
	return nil
}

func Isrepotobeupdate(url string,directory string)(bool,error){
	if!(fileExists(directory)){
		log.Println("There is no file in directory")
		return true,nil
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

func errorfunc(err error) error{
	if (err!=nil){
		log.Fatalln(err)
		return err
	}
	return nil
}

func gitclone(url string, directory string) (*object.Commit,error){
	log.Printf("git clone %s %s --recursive", url, directory)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	CheckIfError(err)
	ref, err := r.Head()
	CheckIfError(err)
	// ... retrieving the commit object
	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)
	return commit,err
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}