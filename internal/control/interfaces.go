package control


type Store interface {

	Save(workflow *Workflow) error

	Get(id string) (*Workflow,error)

	Delete(id string) error

	List() []*Workflow
}