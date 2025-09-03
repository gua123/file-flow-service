package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"file-flow-service/internal/test/testutils"
)

func TestTaskManagerBasic(t *testing.T) {
	// 创建测试任务
	task := testutils.CreateTestTask()
	
	// 验证任务结构
	assert.NotEmpty(t, task.ID)
	assert.NotEmpty(t, task.Name)
	assert.Equal(t, "pending", task.Status)
	
	// 验证任务字段
	assert.NotEmpty(t, task.Creator)
	assert.NotEmpty(t, task.AssignedTo)
	assert.NotEmpty(t, task.Cmd)
	assert.NotNil(t, task.Args)
}

func TestTaskCreation(t *testing.T) {
	// 创建测试任务
	task := testutils.CreateTestTask()
	
	// 验证任务创建的基本属性
	assert.NotEmpty(t, task.ID)
	assert.Equal(t, "Test Task", task.Name)
	assert.Equal(t, "A test task for testing purposes", task.Description)
	assert.Equal(t, "pending", task.Status)
	assert.Equal(t, 0.0, task.Progress)
	assert.NotEmpty(t, task.Creator)
	assert.NotEmpty(t, task.AssignedTo)
	assert.NotEmpty(t, task.ResultPath)
	assert.Equal(t, "echo", task.Cmd)
	assert.Equal(t, []string{"hello"}, task.Args)
	assert.Equal(t, "/tmp", task.Dir)
}
