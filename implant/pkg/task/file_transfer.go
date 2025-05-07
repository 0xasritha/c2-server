package task

// TODO
type FileTransfer struct {
	ID uint `json:"id"`
}

func (f *FileTransfer) Do() Result {
	return *new(Result)
}

func (f *FileTransfer) GetID() uint {
	return f.ID
}
