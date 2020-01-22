package structtype

type Reciver struct{
	Items []Item `json:"items"`
}

type Item struct {
	Meta Metadata `json:"metadata"`
	Spec Spec	`json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
	Annotations Annotation `json:"annotation"`
}

type Annotation struct {
	AppCDoption string `json:"appcdoption"`
}

type Spec struct {
	Source Source `json:"source"`
	Dest Destination `json:"destination"`
	Project string `json:"project"`
	Sync *Syncpolicy `json:"syncPolicy,omitempty" protobuf:"bytes,4,name=syncPolicy"`
}

type Syncpolicy struct{
	Automated *SyncPolicyAutomated `json:"automated,omitempty"`
}

type SyncPolicyAutomated struct{
	// Prune will prune resources automatically as part of automated sync (default: false)
	Prune bool `json:"prune,omitempty" protobuf:"bytes,1,opt,name=prune"`
	// SelfHeal enables auto-syncing if  (default: false)
	SelfHeal bool `json:"selfHeal,omitempty" protobuf:"bytes,2,opt,name=selfHeal"`
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

