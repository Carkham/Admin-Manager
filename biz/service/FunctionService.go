package service

import (
	"admin/conf"
	"admin/model"
	"admin/utils"
	"fmt"
	"log"
	"strconv"
)

func GetFunctionList(userID int64, userName string) (list []model.GetFuncListItem, err error) {
	var funcList []model.DBFuncOverview
	var funcMap map[int64]*model.GetFuncListItem
	err = model.Q.Function.Where(model.Q.Function.UserID.Eq(userID)).
		LeftJoin(model.Q.Template, model.Q.Function.TemplateID.EqCol(model.Q.Template.TemplateID)).
		LeftJoin(model.Q.Trigger, model.Q.Function.TriggerID.EqCol(model.Q.Trigger.TriggerID)).
		Select(
			model.Q.Function.FunctionID,
			model.Q.Function.FunctionLabel,
			model.Q.Template.TemplateLabel,
			model.Q.Trigger.TriggerType,
		).
		Scan(&funcList)

	for _, v := range funcList {
		accessUrl := fmt.Sprintf(
			"http://%s/%s/%s/",
			conf.Config.Service.ExposeIP,
			userName,
			v.FunctionLabel,
		)

		funcMap[v.FunctionID] = &model.GetFuncListItem{
			FunctionID:   v.FunctionID,
			FunctionName: v.FunctionLabel,
			AccessURL:    accessUrl,
			TemplateName: v.TemplateLabel,
			State:        "Stop",
		}
	}

	if err != nil {
		log.Printf("[Get Func List] Get Function List Error: %s", err.Error())
		return nil, err
	}

	depList, err := utils.GetDeploymentList(userID)

	for idStr, v := range depList {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		if item, exist := funcMap[id]; exist {
			item.State = v.Status
		}
	}

	if err != nil {
		log.Printf("[Get Func List] Get Deployment List Info of User %d Error: %s", userID, err.Error())
		return nil, err
	}

	return
}

func DeleteFunc(functionID int64) (err error) {

	dbFuncInfo, err := model.Q.Function.Where(
		model.Q.Function.FunctionID.Eq(functionID),
	).First()

	if err != nil {
		log.Printf("[Delete Function] Get Function Info Error: %s", err.Error())
		return err
	}

	dbTrigInfo, err := model.Q.Trigger.Where(
		model.Q.Trigger.TriggerID.Eq(dbFuncInfo.TriggerID),
	).First()

	if err != nil {
		log.Printf("[Delete Function] Get Function Info Error: %s", err.Error())
		return err
	}

	user, err := model.Q.UserUser.Where(model.Q.UserUser.ID.Eq(dbFuncInfo.UserID)).First()
	if err != nil {
		log.Printf("[Delete Function] Get User Error: %s", err.Error())
		return err
	}

	funcInfo := model.FuncInfo{
		UserName:  user.Username,
		FuncLabel: dbFuncInfo.FunctionLabel,
		TrigType:  dbTrigInfo.TriggerType,
	}

	_ = utils.StopDeployment(&funcInfo)

	if err != nil {
		log.Printf("[Delete Function] Delete Function Error: %s", err.Error())
		return err
	}

	deleteResp, err := model.Q.Function.Where(
		model.Q.Function.FunctionID.Eq(functionID),
	).Delete()

	if err != nil {
		log.Printf("[Delete Function] Delete Function Error: %s", err.Error())
		return err
	}

	if deleteResp.Error != nil {
		log.Printf("[Delete Function] Delete Function Error: %s", deleteResp.Error.Error())
		return deleteResp.Error
	}

	if deleteResp.RowsAffected < 1 {
		log.Printf(
			"[Delete Function] Delete Function Error: no function %d",
			functionID,
		)
		return fmt.Errorf("no function %d", functionID)
	}
	return nil
}

func GetFuncInfo(funcID, userID int64) (data model.GetFuncInfoRespData, err error) {
	var dbFuncInfos []model.DBFuncInfo
	err = model.Q.Function.Where(
		model.Q.Function.FunctionID.Eq(funcID),
		model.Q.Function.UserID.Eq(userID),
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
		err = fmt.Errorf("no such function %d of user %d", funcID, userID)
		return
	}

	dbFuncInfo := dbFuncInfos[0]

	data.FunctionName = dbFuncInfo.FunctionLabel
	data.TriggerID = dbFuncInfo.TriggerID
	data.Replicas = dbFuncInfo.Replicas
	data.SourceLoc = dbFuncInfo.SrcLoc
	data.SourceType = dbFuncInfo.SrcType

	podInfos, err := utils.GetMetricsList(funcID)
	for _, pod := range podInfos {
		data.ReplicasInfo = append(data.ReplicasInfo, model.FuncReplicasInfo{
			Name:     pod.PodName,
			NodeName: pod.NodeName,
			CPUUsage: pod.CpuUsage,
			MemUsage: pod.MemUsage,
			GpuUsage: pod.GpuUsage,
			State:    pod.State,
		})
	}
	return
}
