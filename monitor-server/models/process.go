package models

type ProcessMonitorTable struct {
	Id  int  `json:"id"`
	EndpointId  int  `json:"endpoint_id"`
	Name  string  `json:"name"`
	DisplayName  string  `json:"display_name"`
}

type ProcessUpdateDto struct {
	EndpointId  int  `json:"endpoint_id" binding:"required"`
	ProcessList  []ProcessMonitorTable  `json:"process_list"`
	Check       bool  `json:"check"`
}

type AliveCheckQueueTable struct {
	Id  int  `json:"id"`
	Message  string  `json:"message"`
}

type ProcessUpdateDtoNew struct {
	EndpointId  int  `json:"endpoint_id" binding:"required"`
	ProcessList  []ProcessMonitorTable  `json:"process_list"`
	Check       bool  `json:"check"`
}