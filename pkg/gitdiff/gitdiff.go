package gitdiff

import(
	"github.com/rbxorkt12/applink/pkg/config"
	"log"
	"os"
	"reflect"
)

func Ifchangefilemake(directory string) error {
	diff,old,new:=IsdifferenceConfig(directory)
	if diff==false{
		return nil
	}
	err:=difffilemake(old,new)
	if (err!= nil){
		log.Fatal(err)
		os.Exit(1)
	}
	return nil
}
func difffilemake(old *config.Appoconfig,new *config.Appoconfig) error{

	return nil
}

func IsdifferenceConfig(directory string) (bool,*config.Appoconfig,*config.Appoconfig){
	oldconfig,err:=config.GetConfig(directory)
	if(err!= nil) {
		log.Fatal("Can't marshal old Appoconfig")
		os.Exit(1)
	}
	newconfig,err:=config.GetConfig("new")
	if(err!= nil) {
		log.Fatal("Can't marshal new Appoconfig")
		os.Exit(1)
	}
	if(reflect.DeepEqual(oldconfig,newconfig)) {
		log.Println("You changed remote, but no change in Appoconfig file")
		os.Remove(directory)
		os.Rename("new",directory)
		return false,nil,nil
	}
	return true,oldconfig,newconfig
}
