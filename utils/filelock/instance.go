package filelock

import (
	"file-flow-service/config"
	"file-flow-service/utils/logger"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type Instance struct {
	Config      config.AppConfig
	Logger      logger.Logger
	LockManager *FileLockManager // 新增锁管理器
}

func NewInstance(cfg config.AppConfig, logger logger.Logger) *Instance {
	return &Instance{
		Config:      cfg,
		Logger:      logger,
		LockManager: NewFileLockManager(), // 初始化锁管理器
	}
}

// 文件操作核心方法（新增锁控制）
func (i *Instance) handleFileOperation(c *gin.Context, path string, op func() error) {
	lock := i.LockManager.GetLock(path)
	lock.Lock()
	defer lock.Unlock()

	// 权限检查保持不变
	if !i.hasPermission(path) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	err := op()
	if err != nil {
		i.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

// 权限检查（使用配置模块方法）
func (i *Instance) hasPermission(path string) bool {
	return i.Config.AllowPath(filepath.Dir(path))
}

// Web接口实现
func (i *Instance) Routes() http.Handler {
	r := gin.Default()

	r.POST("/api/file/create", func(c *gin.Context) {
		path := c.PostForm("path")
		content := c.PostForm("content")
		i.handleFileOperation(c, path, func() error {
			return os.WriteFile(path, []byte(content), i.Config.DefaultFileMode)
		})
	})

	r.POST("/api/file/edit", func(c *gin.Context) {
		path := c.PostForm("path")
		content := c.PostForm("content")
		i.handleFileOperation(c, path, func() error {
			return ioutil.WriteFile(path, []byte(content), i.Config.DefaultFileMode)
		})
	})

	r.DELETE("/api/file/delete", func(c *gin.Context) {
		path := c.PostForm("path")
		i.handleFileOperation(c, path, func() error {
			return os.Remove(path)
		})
	})

	r.POST("/api/file/copy", func(c *gin.Context) {
		src := c.PostForm("src")
		dst := c.PostForm("dst")
		i.handleFileOperation(c, src, func() error {
			return copyFile(dst, src)
		})
	})

	r.POST("/api/file/move", func(c *gin.Context) {
		src := c.PostForm("src")
		dst := c.PostForm("dst")
		i.handleFileOperation(c, src, func() error {
			return os.Rename(src, dst)
		})
	})

	r.GET("/api/file/download", func(c *gin.Context) {
		path := c.Query("path")
		i.handleFileOperation(c, path, func() error {
			c.Header("Content-Disposition", "attachment; filename="+filepath.Base(path))
			c.Header("Content-Type", "application/octet-stream")
			http.ServeFile(c.Writer, c.Request, path)
			return nil
		})
	})

	// 新增运行接口
	r.POST("/api/file/run", func(c *gin.Context) {
		path := c.PostForm("path")
		i.handleFileOperation(c, path, func() error {
			return i.runFile(path)
		})
	})

	return r
}

// 新增运行文件实现
func (i *Instance) runFile(path string) error {
	cmd := exec.Command("python", path) // 修改点：添加python
	return cmd.Run()
}

// 辅助函数
func copyFile(dstName, srcName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}