package listener

import (
	"asritha.dev/c2server/pkg/model"
	"gorm.io/gorm"
)

// getPendingTasks gets all incomplete tasks for the given implant
func getPendingTasksForImplant(db *gorm.DB, implantID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := db.
		Where("implant_id = ? AND completed = ?", implantID, false).
		Find(&tasks).
		Error
	return tasks, err

}

func storeNewResult(db *gorm.DB, implantID uint, res model.Result) error {
	// TODO: make it so don't repeat this -> could move creating Result to storeNewResult, that would make way more sense
	return db.Transaction(func(tx *gorm.DB) error {
		var task model.Task
		if err := tx.
			Where("id = ? AND implant_id = ?", res.TaskID, implantID).
			First(&task).Error; err != nil {
			return err
		}

		if err := tx.Create(&res).Error; err != nil {
			return err
		}

		return tx.
			Model(&task).
			Update("completed", true).
			Error
	})
}
