package models

import "time"

type LogMonitorTemplate struct {
	Guid             string    `json:"guid" xorm:"guid"`
	Name             string    `json:"name" xorm:"name"`
	LogType          string    `json:"log_type" xorm:"log_type"`
	DemoLog          string    `json:"demo_log" xorm:"demo_log"`
	JsonRegular      string    `json:"json_regular" xorm:"json_regular"`
	CalcResult       string    `json:"-" xorm:"calc_result"`
	CreateUser       string    `json:"create_user" xorm:"create_user"`
	UpdateUser       string    `json:"update_user" xorm:"update_user"`
	CreateTime       time.Time `json:"-" xorm:"create_time"`
	CreateTimeString string    `json:"create_time"`
	UpdateTime       time.Time `json:"-" xorm:"update_time"`
	UpdateTimeString string    `json:"update_time"`
}

type LogMonitorTemplateRole struct {
	Guid               string `json:"guid" xorm:"guid"`
	LogMonitorTemplate string `json:"log_monitor_template" xorm:"log_monitor_template"`
	Role               string `json:"role" xorm:"role"`
	Permission         string `json:"permission" xorm:"permission"`
}

type LogParamTemplate struct {
	Guid               string    `json:"guid" xorm:"guid"`
	LogMonitorTemplate string    `json:"log_monitor_template" xorm:"log_monitor_template"`
	Name               string    `json:"name" xorm:"name"`
	DisplayName        string    `json:"display_name" xorm:"display_name"`
	JsonKey            string    `json:"json_key" xorm:"json_key"`
	Regular            string    `json:"regular" xorm:"regular"`
	DemoMatchValue     string    `json:"demo_match_value" xorm:"demo_match_value"`
	CreateUser         string    `json:"create_user" xorm:"create_user"`
	UpdateUser         string    `json:"update_user" xorm:"update_user"`
	CreateTime         time.Time `json:"create_time" xorm:"create_time"`
	UpdateTime         time.Time `json:"update_time" xorm:"update_time"`
}

func (l *LogParamTemplate) TransToLogParam() (output *LogMetricParamObj) {
	output = &LogMetricParamObj{}
	output.Guid = l.Guid
	output.Name = l.Name
	output.DisplayName = l.DisplayName
	output.Regular = l.Regular
	output.JsonKey = l.JsonKey
	return
}

type LogParamTemplateObj struct {
	LogParamTemplate
	StringMap []*LogMetricStringMapTable `json:"string_map"`
}

type LogMetricTemplate struct {
	Guid               string    `json:"guid" xorm:"guid"`
	LogMonitorTemplate string    `json:"log_monitor_template" xorm:"log_monitor_template"`
	LogParamName       string    `json:"log_param_name" xorm:"log_param_name"`
	Metric             string    `json:"metric" xorm:"metric"`
	DisplayName        string    `json:"display_name" xorm:"display_name"`
	Step               int       `json:"step" xorm:"step"`
	AggType            string    `json:"agg_type" xorm:"agg_type"`
	TagConfig          string    `json:"-" xorm:"tag_config"`
	TagConfigList      []string  `json:"tag_config" xorm:"-"`
	CreateUser         string    `json:"create_user" xorm:"create_user"`
	UpdateUser         string    `json:"update_user" xorm:"update_user"`
	CreateTime         time.Time `json:"create_time" xorm:"create_time"`
	UpdateTime         time.Time `json:"update_time" xorm:"update_time"`
}

func (l *LogMetricTemplate) TransToLogMetric() (output *LogMetricConfigTable) {
	output = &LogMetricConfigTable{
		Guid:          l.Guid,
		Metric:        l.Metric,
		DisplayName:   l.DisplayName,
		AggType:       l.AggType,
		Step:          int64(l.Step),
		TagConfig:     l.TagConfig,
		TagConfigList: l.TagConfigList,
		LogParamName:  l.LogParamName,
	}
	return
}

type LogMonitorTemplateDto struct {
	LogMonitorTemplate
	CalcResultObj *CheckRegExpResult            `json:"calc_result"`
	ParamList     []*LogParamTemplate           `json:"param_list"`
	MetricList    []*LogMetricTemplate          `json:"metric_list"`
	Permission    *LogMonitorTemplatePermission `json:"permission"`
}

//func (l *LogMonitorTemplateDto) GetMetrics() (metricConfigList []*LogMetricConfigTable) {
//	for _, metric := range l.MetricList {
//		metricObj := LogMetricConfigTable{Metric: metric.Metric, DisplayName: metric.DisplayName, AggType: metric.AggType, TagConfigList: metric.TagConfigList}
//		for _, param := range l.ParamList {
//			if param.Name == metric.LogParamName {
//				metricObj.JsonKey = param.JsonKey
//				metricObj.Regular = param.Regular
//			}
//		}
//	}
//	return
//}

type LogMonitorTemplatePermission struct {
	MgmtRoles []string `json:"mgmt_roles"`
	UseRoles  []string `json:"use_roles"`
}

type LogMonitorTemplateListParam struct {
	Name       string `json:"name"`
	UpdateUser string `json:"update_user"`
}

type LogMonitorTemplateListResp struct {
	JsonList    []*LogMonitorTemplate `json:"json_list"`
	RegularList []*LogMonitorTemplate `json:"regular_list"`
}

type LogMonitorRegMatchParam struct {
	DemoLog   string                 `json:"demo_log"`
	ParamList []*LogParamTemplateObj `json:"param_list"`
}

type LogMetricGroup struct {
	Guid               string    `json:"guid" xorm:"guid"`
	Name               string    `json:"name" xorm:"name"`
	MetricPrefixCode   string    `json:"metric_prefix_code" xorm:"metric_prefix_code"`
	LogType            string    `json:"log_type" xorm:"log_type"`
	LogMetricMonitor   string    `json:"log_metric_monitor" xorm:"log_metric_monitor"`
	LogMonitorTemplate string    `json:"log_monitor_template" xorm:"log_monitor_template"`
	DemoLog            string    `json:"demo_log" xorm:"demo_log"`
	CalcResult         string    `json:"calc_result" xorm:"calc_result"`
	CreateUser         string    `json:"create_user" xorm:"create_user"`
	UpdateUser         string    `json:"update_user" xorm:"update_user"`
	CreateTime         time.Time `json:"-" xorm:"create_time"`
	CreateTimeString   string    `json:"create_time"`
	UpdateTime         time.Time `json:"-" xorm:"update_time"`
	UpdateTimeString   string    `json:"update_time"`
}

type LogMetricParam struct {
	Guid           string    `json:"guid" xorm:"guid"`
	Name           string    `json:"name" xorm:"name"`
	DisplayName    string    `json:"display_name" xorm:"display_name"`
	LogMetricGroup string    `json:"log_metric_group" xorm:"log_metric_group"`
	Regular        string    `json:"regular" xorm:"regular"`
	DemoMatchValue string    `json:"demo_match_value" xorm:"demo_match_value"`
	CreateUser     string    `json:"create_user" xorm:"create_user"`
	UpdateUser     string    `json:"update_user" xorm:"update_user"`
	CreateTime     time.Time `json:"create_time" xorm:"create_time"`
	UpdateTime     time.Time `json:"update_time" xorm:"update_time"`
}

type LogMetricParamObj struct {
	LogMetricParam
	JsonKey   string                     `json:"json_key"`
	StringMap []*LogMetricStringMapTable `json:"string_map"`
}

type LogTemplateAffectGroupObj struct {
	LogMetricGroupGuid   string
	LogMetricGroupName   string
	LogMetricGroupPrefix string
	LogMetricMonitorGuid string
	ServiceGroup         string
	SucRetCode           string
}

type LogTemplateExportParam struct {
	GuidList []string `json:"guidList"`
}
