package service

import (
	"admin/model"
	"admin/utils"
	"fmt"
	"log"
	"strconv"
)

func GetFunctionList() (list []model.GetFuncList, err error) {
	var funcList []model.DBFuncOverview
	var funcMap map[int64]*model.GetFuncList
	err = model.Q.Function.
		LeftJoin(model.Q.Template, model.Q.Function.TemplateID.EqCol(model.Q.Template.TemplateID)).
		LeftJoin(model.Q.UserUser, model.Q.Function.UserID.EqCol(model.Q.UserUser.ID)).
		Select(
			model.Q.Function.FunctionID,
			model.Q.Function.FunctionLabel,
			model.Q.Template.TemplateLabel,
			model.Q.UserUser.Username,
		).
		Scan(&funcList)

	for _, v := range funcList {
		replicasList, _ := utils.GetMetricsList(v.FunctionID)
		funcMap[v.FunctionID] = &model.GetFuncList{
			UserName:     v.Username,
			FunctionId:   int(v.FunctionID),
			FunctionName: v.FunctionLabel,
			TemplateName: v.TemplateLabel,
			ReplicasInfo: replicasList,
		}
	}

	if err != nil {
		log.Printf("[Get Func List] Get Function List Error: %s", err.Error())
		return nil, err
	}

	depList, err := utils.GetDeploymentList()

	for idStr, v := range depList {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		if item, exist := funcMap[id]; exist {
			item.State = v.Status
		}
	}

	if err != nil {
		log.Printf("[Get Func List] Get Deployment List Info Error: %s", err.Error())
		return nil, err
	}

	for _, v := range funcMap {
		list = append(list, *v)
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
