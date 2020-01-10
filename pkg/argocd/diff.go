package argocd

import (
	structtype "github.com/rbxorkt12/applink/pkg/type"
)
func Appdiff(before []structtype.Item,after []structtype.Item)(create []structtype.Item,delete []structtype.Item,update []structtype.Item){
	var flag int
	for _,app := range before{
		flag = 1
		name := app.Meta.Name
		for _,com := range after{
			if name == com.Meta.Name {
				if app.Spec.Source.Revision == com.Spec.Source.Revision{
					flag =0
					break
				} else {
					update = append(update,com)
				}
			}
		}
		if flag ==1 {
			delete = append(delete,app)
		}
	}

	for _,app := range after{
		flag = 1
		name := app.Meta.Name
		for _,com := range before{
			if name == com.Meta.Name {
				if app.Spec.Source.Revision == com.Spec.Source.Revision{
					flag =0
					break
				}
			}
		}
		if flag ==1 {
			create = append(create,app)
		}
	}
	return create,delete,update

}