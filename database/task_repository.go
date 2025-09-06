package database

import (
	"time"
	"file-flow-service/utils/logger"
	"go.uber.org/zap"
)

// Task represents a task entity
type Task struct {
	ID          string
	Name        string
	Status      string
	Creator     string
	CreatedAt   string
	AssignedTo  string
	Description string
	ResultPath  string
	Progress    int64
	Duration    int64
	FinishedAt  int64
	StartedAt   int64 // Renamed from StartTime to match service.Task's StartedAt
}

// TaskInterface methods implementation
func (t *Task) GetID() string {
	return t.ID
}

func (t *Task) GetName() string {
	return t.Name
}

func (t *Task) GetStatus() string {
	return t.Status
}

func (t *Task) GetCreator() string {
	return t.Creator
}

func (t *Task) GetCreatedAt() string {
	return t.CreatedAt
}

func (t *Task) GetAssignedTo() string {
	return t.AssignedTo
}

func (t *Task) GetDescription() string {
	return t.Description
}

func (t *Task) GetResultPath() string {
	return t.ResultPath
}

func (t *Task) GetProgress() int64 {
	return t.Progress
}

func (t *Task) GetDuration() int64 {
	return t.Duration
}

func (t *Task) GetFinishedAt() int64 {
	return t.FinishedAt
}

func (t *Task) GetStartTime() int64 {
	return t.StartedAt // Now correctly refers to StartedAt
}

func (t *Task) SetStatus(status string) {
	t.Status = status
}

func (t *Task) SetStartTime(startTime int64) {
	t.StartedAt = startTime // Now sets StartedAt
}

func (t *Task) SetDuration(duration int64) {
	t.Duration = duration
}

func (t *Task) SetFinishedAt(finishTime int64) {
	t.FinishedAt = finishTime
}

// Execute implementation for TaskInterface
func (t *Task) Execute() error {
	// Placeholder for task execution logic
	return nil
}

// CreateTask inserts a new task into the database
func CreateTask(task *Task) error {
	if task.CreatedAt == "" {
		task.CreatedAt = time.Now().Format(time.RFC3339)
	}
	
	_, err := db.Exec(`
		INSERT INTO tasks (id, name, status, creator, createdAt, assignedTo, description, resultPath, progress, duration, finishedAt, startedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.Name, task.Status, task.Creator, task.CreatedAt, task.AssignedTo, task.Description, task.ResultPath, task.Progress, task.Duration, task.FinishedAt, task.StartedAt,
	)
	if err != nil {
		logger.GetLogger().Error("创建任务失败", zap.Error(err))
		return err
	}
	logger.GetLogger().Error("成功创建任务", zap.String("id", task.ID))
	return nil
}

// GetTaskByID retrieves a task by ID
func GetTaskByID(id string) (*Task, error) {
	row := db.QueryRow("SELECT id, name, status, creator, createdAt, assignedTo, description, resultPath, progress, duration, finishedAt, startedAt FROM tasks WHERE id = ?", id)
	
	var task Task
	err := row.Scan(&task.ID, &task.Name, &task.Status, &task.Creator, &task.CreatedAt, &task.AssignedTo, &task.Description, &task.ResultPath, &task.Progress, &task.Duration, &task.FinishedAt, &task.StartedAt)
	if err != nil {
		logger.GetLogger().Error("查询任务失败", zap.Error(err))
		return nil, err
	}
	return &task, nil
}

// UpdateTask updates an existing task
func UpdateTask(task *Task) error {
	_, err := db.Exec(`
		UPDATE tasks 
		SET name = ?, status = ?, creator = ?, assignedTo = ?, description = ?, resultPath = ?, progress = ?, duration = ?, finishedAt = ?, startedAt = ?
		WHERE id = ?`,
		task.Name, task.Status, task.Creator, task.AssignedTo, task.Description, task.ResultPath, task.Progress, task.Duration, task.FinishedAt, task.StartedAt, task.ID,
	)
	if err != nil {
		logger.GetLogger().Error("更新任务失败", zap.Error(err))
		return err
	}
	return nil
}

// DeleteTask removes a task by ID
func DeleteTask(id string) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		logger.GetLogger().Error("删除任务失败", zap.Error(err))
		return err
	}
	return nil
}

// GetTasks retrieves all tasks
func GetTasks() ([]Task, error) {
	rows, err := db.Query("SELECT id, name, status, creator, createdAt, assignedTo, description, resultPath, progress, duration, finishedAt, startedAt FROM tasks")
	if err != nil {
		logger.GetLogger().Error("获取任务列表失败", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Status, &task.Creator, &task.CreatedAt, &task.AssignedTo, &task.Description, &task.ResultPath, &task.Progress, &task.Duration, &task.FinishedAt, &task.StartedAt); err != nil {
			logger.GetLogger().Error("任务扫描失败", zap.Error(err))
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}