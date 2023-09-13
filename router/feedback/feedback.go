package feedback

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/database"
	"github.com/RaymondCode/simple-demo/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type getFeedbackDeviceParam struct {
	StartTime int64 `form:"startTime" binding:"required"`
	EndTime   int64 `form:"endTime" binding:"required"`
}

func GetFeedback(ctx *gin.Context) {
	g := common.GetGin(ctx)
	param := getFeedbackDeviceParam{}
	err := ctx.ShouldBindQuery(&param)
	if err != nil {
		log.Errorf("get param fail, err:%v", err)
		g.ResponseFail()
		return
	}
	feedback := models.Feedback{}
	db := database.GetInstanceConnection().GetPrimaryDB()
	feedbacks, err := feedback.GetFeedbackOfTime(db, param.StartTime, param.EndTime)
	if err != nil {
		log.Errorf("get indoor device info fail, err:%v", err)
		g.ResponseFail()
		return
	}
	g.ResponseSuccess(feedbacks)
}

func GetLatestFeedBack(ctx *gin.Context) {
	g := common.GetGin(ctx)
	db := database.GetInstanceConnection().GetPrimaryDB()
	feedback := models.Feedback{}
	resultType := ctx.Query("resultType")
	page := 0
	limit := 1000
	ans := make([]models.Feedback, 0)

	for {
		feedbacks, err := feedback.GetLatestFeedback(db, resultType, "Normal", page, limit)
		page++
		if err != nil {
			log.Errorf("get indoor device info fail, err:%v", err)
			g.ResponseFail()
			return
		}

		if feedbacks == nil || len(feedbacks) == 0 {
			break
		}

		for _, fb := range feedbacks {
			if len(ans) == 0 {
				ans = append(ans, fb)
				continue
			}
			tail := &ans[len(ans)-1]
			if fb.Pk == tail.Pk-1 {
				tail.Pk = fb.Pk
				tail.StartTime = fb.StartTime
			} else {
				ans = append(ans, fb)
				if len(ans) == 10 {
					log.Infof("ans: %v", ans)
					g.ResponseNormal(ans)
					return
				}
			}
		}
	}
	g.ResponseNormal(ans)
}
