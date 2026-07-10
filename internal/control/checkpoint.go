package control

import "time"


func (c *Controller) CreateCheckpoint(
	id string,
	step uint64,
	payload []byte,
) error{


	c.mu.Lock()
	defer c.mu.Unlock()


	w,ok:=c.workflows[id]


	if !ok{
		return ErrWorkflowNotFound
	}


	w.Checkpoint=&Checkpoint{

		ID:id,

		WorkflowID:id,

		Step:step,

		Payload:payload,

		CreatedAt:time.Now(),
	}


	w.CurrentStep=step

	w.UpdatedAt=time.Now()


	return nil

}



func (c *Controller) RestoreCheckpoint(id string)(*Checkpoint,error){


	c.mu.RLock()
	defer c.mu.RUnlock()


	w,ok:=c.workflows[id]


	if !ok{
		return nil,ErrWorkflowNotFound
	}


	if w.Checkpoint==nil{
		return nil,ErrCheckpointNotFound
	}


	return w.Checkpoint,nil

}