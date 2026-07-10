package control

import "time"


func (c *Controller) Start(
	id string,
	name string,
) (*Workflow,error){


	c.mu.Lock()
	defer c.mu.Unlock()


	if _,exists:=c.workflows[id]; exists{
		return nil,ErrWorkflowExists
	}


	w:=&Workflow{

		ID:id,

		Name:name,

		Status:StatusRunning,

		CurrentStep:0,

		CreatedAt:time.Now(),

		UpdatedAt:time.Now(),

	}


	c.workflows[id]=w


	return w,nil
}



func (c *Controller) Fail(id string) error{


	c.mu.Lock()
	defer c.mu.Unlock()


	w,ok:=c.workflows[id]

	if !ok{
		return ErrWorkflowNotFound
	}


	w.Status=StatusFailed

	w.UpdatedAt=time.Now()


	return nil

}



func (c *Controller) Resume(id string) error{


	c.mu.Lock()
	defer c.mu.Unlock()


	w,ok:=c.workflows[id]

	if !ok{
		return ErrWorkflowNotFound
	}


	w.Status=StatusRunning

	w.UpdatedAt=time.Now()


	return nil

}