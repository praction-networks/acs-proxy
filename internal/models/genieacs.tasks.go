package models

type TaskType string

const (
	TaskGetParameterValues TaskType = "getParameterValues"
	TaskSetParameterValues TaskType = "setParameterValues"
	TaskRefreshObject      TaskType = "refreshObject"
	TaskAddObject          TaskType = "addObject"
	TaskDeleteObject       TaskType = "deleteObject"
	TaskReboot             TaskType = "reboot"
	TaskFactoryReset       TaskType = "factoryReset"
)

// ---------------------- Task Requests ----------------------

// ---------------------- Task Requests ----------------------

type GetParameterValuesRequest struct {
	ParameterNames []string `json:"parameterNames" validate:"required,dive,required"`
}

type SetParameterValuesRequest struct {
	ParameterValues [][]interface{} `json:"parameterValues" validate:"required,dive,min=1"`
}

type RefreshObjectRequest struct {
	ObjectName string `json:"objectName" validate:"required"`
}

type AddObjectRequest struct {
	ObjectName string `json:"objectName" validate:"required"`
}

type DeleteObjectRequest struct {
	ObjectName string `json:"objectName" validate:"required"`
}
