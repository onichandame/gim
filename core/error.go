package core

type GimError struct {
	Message string
	Body    interface{}
	Status  int
}

func NewGimError(status int, body interface{},props ...string)*GimError{
	e:=GimError{Status: status,Body: body}
	if len(props)>0{e.Message=props[0]}
	return &e
}

func (e *GimError) Error() string {
	return e.Message
}
