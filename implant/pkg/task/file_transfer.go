package task

// TODO
type FileTransfer struct {
	ID uint `json:"id"`
}

func (f *FileTransfer) Do() Result {
	return Result{
		Content: "FileTransfer tasks not supported yet",
	}
}

func (f *FileTransfer) GetID() uint {
	return f.ID
}
