package file

import (
	"github.com/RaymondCode/simple-demo/conf"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func Upload(c *gin.Context) {
	// 获取上传的文件
	config := conf.OthersConfig

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := file.Filename
	if filename != config.SystemModel && filename != config.IuModel {
		log.Errorf("fileName %v is not valid", filename)
		c.JSON(http.StatusInternalServerError, gin.H{"filename is not valid": err.Error()})
		return
	}

	// 将文件保存到服务器路径
	err = c.SaveUploadedFile(file, config.UploadDir+"/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})
}

func Download(c *gin.Context) {
	filename := c.Query("filename")
	config := conf.OthersConfig
	fileLocation := config.UploadDir + "/" + filename
	// 检查文件是否存在
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}
	// 设置响应头，告诉浏览器以附件方式下载文件
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(fileLocation)
}
