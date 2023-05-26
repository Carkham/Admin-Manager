package service

import (
	"admin/model"
	UserUser "admin/model/model"
	"admin/utils"
	"encoding/json"
	"fmt"
	"log"
)

func StartFunc(funcID int64) (err error) {
	var dbFuncInfos []model.DBFuncInfo
	err = model.Q.Function.Where(
		model.Q.Function.FunctionID.Eq(funcID),
	).LeftJoin(
		model.Q.Trigger,
		model.Q.Function.TriggerID.EqCol(model.Q.Trigger.TriggerID),
	).LeftJoin(
		model.Q.Template,
		model.Q.Function.TemplateID.EqCol(model.Q.Template.TemplateID),
	).Select(
		model.Q.Function.FunctionID,
		model.Q.Function.UserID,
		model.Q.Function.FunctionLabel,
		model.Q.Function.SrcType,
		model.Q.Function.SrcLoc,
		model.Q.Trigger.TriggerType,
		model.Q.Trigger.TriggerConfig,
		model.Q.Function.Replicas,
		model.Q.Template.ImageName,
		model.Q.Function.QuotaInfo,
	).Scan(&dbFuncInfos)

	if err != nil {
		log.Printf("[Start Function] Query Function Info Error : %s", err.Error())
		return
	}

	if len(dbFuncInfos) < 1 {
		log.Printf("[Start Function] Query Function Info Error: No Such Function")
		err = fmt.Errorf("no such function %d", funcID)
		return
	}

	dbFuncInfo := dbFuncInfos[0]
	var cronConf model.CronJobConfig
	var quotaConf model.QuotaInfo
	_ = json.Unmarshal([]byte(dbFuncInfo.TriggerConfig), &cronConf)
	_ = json.Unmarshal([]byte(dbFuncInfo.QuotaInfo), &cronConf)

	var users []UserUser.UserUser
	err = model.Q.UserUser.Where(model.Q.UserUser.ID.Eq(dbFuncInfo.UserID)).Scan(&users)
	if err != nil {
		log.Printf("[Start Function] Query Function Info Error : %s", err.Error())
		return
	}

	if len(users) < 1 {
		log.Printf("[Start Function] Query Function Info Error: No User")
		err = fmt.Errorf("no such user %d", dbFuncInfo.UserID)
		return
	}

	funcInfo := model.FuncInfo{
		FunctionID:     funcID,
		UserID:         dbFuncInfo.UserID,
		UserName:       users[0].Username,
		FuncLabel:      dbFuncInfo.FunctionLabel,
		SourceType:     dbFuncInfo.SrcType,
		SourceLocation: dbFuncInfo.SrcLoc,
		TrigType:       dbFuncInfo.TriggerType,
		TimeStr:        utils.ParseCronExpr(&cronConf),
		PodCount:       dbFuncInfo.Replicas,
		ImageName:      dbFuncInfo.ImageName,
		CPUQuotaM:      [2]int{quotaConf.CpuRequest, quotaConf.CpuLimit},
		MemQuotaMi:     [2]int{quotaConf.MemRequest, quotaConf.MemLimit},
		GPUQuota:       quotaConf.GpuQuota,
	}

	depName, err := utils.StartDeployment(&funcInfo)

	if err != nil {
		log.Printf("[Start Function] Start Deployment Error: %s", err.Error())
	}
	log.Printf("[Start Function] Function %d Successfully Start as %s", funcID, depName)
	return
}

func StopFunc(userID int64, funcID int64, userName string) (err error) {

	dbFuncInfo, err := model.Q.Function.Where(
		model.Q.Function.FunctionID.Eq(funcID),
		model.Q.Function.UserID.Eq(userID),
	).First()

	if err != nil {
		log.Printf("[Stop Function] Query Function Info Error: %s", err.Error())
	}

	dbTrigInfo, err := model.Q.Trigger.Where(
		model.Q.Trigger.TriggerID.Eq(dbFuncInfo.TriggerID),
	).First()

	if err != nil {
		log.Printf("[Delete Function] Get Function Info Error: %s", err.Error())
		return err
	}

	funcInfo := model.FuncInfo{
		UserName:  userName,
		FuncLabel: dbFuncInfo.FunctionLabel,
		TrigType:  dbTrigInfo.TriggerType,
	}

	err = utils.StopDeployment(&funcInfo)

	if err != nil {
		log.Printf("[Stop Function] Stop Function %d Error: %s", funcID, err.Error())
	}

	return
}
