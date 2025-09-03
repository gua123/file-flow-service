// filemanager.go
// 文件管理模块，处理文件上传、下载
// 管理执行环境文件夹

package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"file-flow-service/utils/logger"
	"mime/multipart"

	"go.uber.org/zap"
)

type FileService struct {
	StoragePath string
	Logger      logger.Logger
}

func NewFileService(storagePath string, logger logger.Logger) *FileService {
	return &FileService{
		StoragePath: storagePath,
		Logger:      logger,
	}
}

// Upload 上传文件
// 参数：file 文件头
// 返回：错误信息
func (f *FileService) Upload(file *multipart.FileHeader) error {
	// 使用zap的字段构造方式
	f.Logger.Info("文件上传", zap.String("filename", file.Filename))
	
	// 确保存储目录存在
	err := os.MkdirAll(f.StoragePath, 0755)
	if err != nil {
		f.Logger.Error("创建存储目录失败", zap.Error(err))
		return fmt.Errorf("创建存储目录失败: %v", err)
	}
	
	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		f.Logger.Error("打开上传文件失败", zap.Error(err))
		return fmt.Errorf("打开上传文件失败: %v", err)
	}
	defer src.Close()
	
	// 创建目标文件
	dstPath := filepath.Join(f.StoragePath, file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		f.Logger.Error("创建目标文件失败", zap.Error(err))
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()
	
	// 复制文件内容
	_, err = io.Copy(dst, src)
	if err != nil {
		f.Logger.Error("复制文件内容失败", zap.Error(err))
		return fmt.Errorf("复制文件内容失败: %v", err)
	}
	
	f.Logger.Info("文件上传成功", zap.String("filepath", dstPath))
	return nil
}

// Download 下载文件
// 参数：fileID 文件ID
// 返回：文件路径，错误信息
func (f *FileService) Download(fileID string) (string, error) {
	// 构建文件路径
	filePath := filepath.Join(f.StoragePath, fileID)
	
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		f.Logger.Error("文件不存在", zap.String("filepath", filePath), zap.Error(err))
		return "", fmt.Errorf("文件不存在: %v", err)
	}
	
	return filePath, nil
}
