package folder

import (
	"github.com/PointStoneTeam/cugsync/setting"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type File struct {
	Name           string    `json:"name"`
	Path           string    `json:"path"`
	IsDir          bool      `json:"is_dir"`
	LastModifyTime time.Time `json:"last_modify_time"`
}

func GetFolder(path string) ([]*File, error) {
	// 默认放在 /data1/
	pathPreix := setting.GetStorePath()
	totalPath := pathPreix + path
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(totalPath)

	if err != nil {
		log.WithField("err", err).Info("GetFolder 打开文件夹出错")
		return nil, err
	}
	return ConvertFolder(fileInfoList, path), nil
}

func ConvertFolder(list []os.FileInfo, p string) []*File {
	var fileList []*File
	for _, fileInfo := range list {
		fileList = append(fileList, &File{
			Name:           fileInfo.Name(),
			Path:           path.Join(p, fileInfo.Name()),
			IsDir:          fileInfo.IsDir(),
			LastModifyTime: fileInfo.ModTime(),
		})
	}
	return fileList
}
