package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"server/config"
	"server/logic"
	"server/model/system"
	"server/plugin/common/util"
)

// SingleUpload 单文件上传, 暂定为图片上传
func SingleUpload(c *gin.Context) {
	// 获取执行操作的用户信息
	v, ok := c.Get(config.AuthUserClaims)
	if !ok {
		system.Failed("上传失败, 当前用户信息异常", c)
		return
	}
	// 结合搜文件内容
	file, err := c.FormFile("file")
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}
	// 创建文件保存路径, 如果不存在则创建
	//if _, err = os.Stat(config.ImageDir); os.IsNotExist(err) {
	//	err = os.MkdirAll(config.ImageDir, os.ModePerm)
	//	if err != nil {
	//		return
	//	}
	//}

	// 生成文件名, 保存文件到服务器
	fileName := fmt.Sprintf("%s/%s%s", config.FilmPictureUploadDir, util.RandomString(8), filepath.Ext(file.Filename))
	err = c.SaveUploadedFile(file, fileName)
	if err != nil {
		system.Failed(err.Error(), c)
		return
	}

	uc := v.(*system.UserClaims)
	// 记录图片信息到系统表中, 并获取返回的图片访问路径
	link := logic.FileL.SingleFileUpload(fileName, int(uc.UserID))
	// 返回图片访问地址以及成功的响应
	system.Success(link, "上传成功", c)

}

// PhotoWall 照片墙数据
func PhotoWall(c *gin.Context) {

	// 获取系统保存的文件的图片分页数据
	page := system.Page{PageSize: 10, Current: 1}
	// 获取分页数据
	pl := logic.FileL.GetPhotoPage(&page)
	system.Success(pl, "图片分页数据获取成功", c)
}
