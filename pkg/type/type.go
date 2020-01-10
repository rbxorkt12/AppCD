package structtype

type Reciver struct{
	Items []Item `json:"items"`
}

type Item struct {
	Meta Metadata `json:"metadata"`
	Spec Spec	`json"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

type Spec struct {
	Source Source `json:"source"`
	Dest Destination `json:"destination"`
}

type Source struct {
	Url string `json:"repoURL"`
	Path string `json:"path"`
	Revision string `json:"targetRevision"`
}

type Destination struct {
	Server	string `json:"server"`
	Namespace string `json:"namespace"`
}

func NewReciver() *Reciver{
	return &Reciver{}
}

func NeWItemlist() []Item{
	var itemlist []Item
	return itemlist
}