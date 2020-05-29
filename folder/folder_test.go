package folder

import (
	"testing"
)

func TestGetFolder(t *testing.T) {
	if fileInfoList, err := GetFolder("../../test_rsync/"); err != nil {
		t.Log(err)
	} else {
		for i := range fileInfoList {
			t.Log(fileInfoList[i]) //打印当前文件或目录下的文件或目录名
		}
	}
}
