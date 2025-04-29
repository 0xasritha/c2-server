package taskrecords

type taskRecordResult struct {
	content   []string `json:"content"`
	isSuccess bool     `json:"success"`
}

type TaskRecord struct {
	ID          uint32             `json:"id"`
	TaskType    string             `json:"type"`
	TaskOptions map[string]string  `json:"options"`
	Results     []taskRecordResult `json:"results"`
}

// TODO: should these parameter names be capitalized?
func NewTaskRecord(ID uint32, TaskType string, TaskOptions map[string]string) *TaskRecord {
	return &TaskRecord{
		ID:          ID,
		TaskType:    TaskType,
		TaskOptions: TaskOptions,
		Results:     make([]taskRecordResult, 0),
	}
}

func (tr *TaskRecord) AddResult(content []string, isSuccess bool) {
	newTaskRecordResult := taskRecordResult{
		content:   content,
		isSuccess: isSuccess,
	}
	tr.Results = append(tr.Results, newTaskRecordResult)
}

// TODO: create add multiple results function
