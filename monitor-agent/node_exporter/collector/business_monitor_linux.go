package collector

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/hpcloud/tail"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	businessCollectorName = "business_monitor"
	businessMonitorFilePath = "data/business_monitor_cache.data"
)

var (
	businessMonitorJobs []*businessMonitorObj
	businessMonitorLock = new(sync.RWMutex)
	businessMonitorMetrics []*businessRuleMetricObj
	businessMonitorMetricLock = new(sync.RWMutex)
	newLogger  log.Logger
)

type businessMonitorCollector struct {
	businessMonitor  *prometheus.Desc
	logger  log.Logger
}

func InitNewLogger(logger  log.Logger)  {
	newLogger = logger
}

func init() {
	registerCollector(businessCollectorName, defaultEnabled, BusinessMonitorCollector)
}

func BusinessMonitorCollector(logger log.Logger) (Collector, error) {
	return &businessMonitorCollector{
		businessMonitor: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, businessCollectorName, "value"),
			"Show business data from log file.",
			[]string{"key", "tags", "path", "agg"}, nil,
		),
		logger: logger,
	}, nil
}

func (c *businessMonitorCollector) Update(ch chan<- prometheus.Metric) error {
	businessMonitorMetricLock.RLock()
	for _,v := range businessMonitorMetrics {
		ch <- prometheus.MustNewConstMetric(c.businessMonitor,
			prometheus.GaugeValue,
			v.Value, v.Metric, v.TagsString, v.Path, v.Agg)
	}
	businessMonitorMetricLock.RUnlock()
	return nil
}

type businessStoreMonitorObj struct {
	Path  string  `json:"path"`
	Rules  []*businessStoreMetricObj  `json:"rules"`
}

type businessStoreMetricObj struct {
	Regular  string  `json:"regular"`
	StringMap  []*businessStringMapObj  `json:"string_map"`
	TagsKey  []string  `json:"tags_key"`
	TagsValue  []string  `json:"tags_value"`
	MetricConfig  []*businessMetricConfigObj  `json:"metric_config"`
}

type businessRuleObj struct {
	Regular  string  `json:"regular"`
	RegExp  *regexp.Regexp  `json:"-"`
	StringMap  []*businessStringMapObj  `json:"string_map"`
	TagsKey  []string  `json:"tags_key"`
	TagsValue  []string  `json:"tags_value"`
	MetricConfig  []*businessMetricConfigObj  `json:"metric_config"`
	DataChannel chan map[string]interface{}  `json:"-"`
}

type businessRuleMetricObj struct {
	Metric  string  `json:"metric"`
	Path  string  `json:"path"`
	Agg   string  `json:"agg"`
	Tags    []string  `json:"tags"`
	TagsString string  `json:"tags_string"`
	Value   float64  `json:"value"`
}

type businessMonitorObj struct {
	Path  string  `json:"path"`
	TailSession  *tail.Tail  `json:"-"`
	Lock  *sync.RWMutex  `json:"-"`
	Rules  []*businessRuleObj  `json:"rules"`
}

type businessStringMapObj struct {
	Key  string  `json:"key"`
	StringValue  string  `json:"string_value"`
	IntValue  float64  `json:"int_value"`
}

type businessMetricConfigObj struct {
	Key  string  `json:"key"`
	Metric  string  `json:"metric"`
	AggType  string  `json:"agg_type"`
}

type businessMonitorCfgObj struct {
	Regular  string  `json:"regular"`
	Tags  string  `json:"tags"`
	StringMap  []*businessStringMapObj  `json:"string_map"`
	MetricConfig  []*businessMetricConfigObj  `json:"metric_config"`
}

type businessAgentDto struct {
	Path  string  `json:"path"`
	Config  []*businessMonitorCfgObj  `json:"config"`
}

func (c *businessMonitorObj) start()  {
	var err error
	c.TailSession,err = tail.TailFile(c.Path, tail.Config{Follow:true, ReOpen:true})
	if err != nil {
		level.Error(newLogger).Log("msg",fmt.Sprintf("start business collector fail, path: %s, error: %v", c.Path, err))
		return
	}
	for line := range c.TailSession.Lines {
		c.Lock.RLock()
		for _,rule := range c.Rules {
			fetchList := rule.RegExp.FindStringSubmatch(line.Text)
			if len(fetchList) > 1 {
				fetchKeyMap := make(map[string]interface{})
				for i,v := range fetchList {
					if i == 0 {
						continue
					}
					tmpKeyMap := make(map[string]interface{})
					tmpErr := json.Unmarshal([]byte(v), &tmpKeyMap)
					if tmpErr != nil {
						level.Error(newLogger).Log("line fetch regexp fail", fmt.Sprintf("line:%s error:%s", v, tmpErr.Error()))
					}else{
						for tmpKeyMapKey,tmpKeyMapValue := range tmpKeyMap {
							fetchKeyMap[tmpKeyMapKey] = tmpKeyMapValue
						}
					}
				}
				if len(fetchKeyMap) > 0 {
					rule.DataChannel <- fetchKeyMap
				}
			}
		}
		c.Lock.RUnlock()
	}
}

func (c *businessMonitorObj) destroy()  {
	c.TailSession.Stop()
	c.Rules = []*businessRuleObj{}
}

func BusinessMonitorHttpHandle(w http.ResponseWriter, r *http.Request) {
	buff, err := ioutil.ReadAll(r.Body)
	var errorMsg string
	if err != nil {
		errorMsg = fmt.Sprintf("Handel business monitor http request fail,read body error: %v \n", err)
		level.Error(newLogger).Log("msg", errorMsg)
		w.Write([]byte(errorMsg))
		return
	}
	var param []*businessAgentDto
	err = json.Unmarshal(buff, &param)
	if err != nil {
		errorMsg = fmt.Sprintf("Handel business monitor http request fail,json unmarshal error: %v \n", err)
		level.Error(newLogger).Log("msg", errorMsg)
		w.Write([]byte(errorMsg))
		return
	}
	businessMonitorLock.Lock()
	var newBmj []*businessMonitorObj
	for _,v := range businessMonitorJobs {
		delFlag := true
		for _,vv := range param {
			if vv.Path == v.Path {
				delFlag = false
				v.Lock.Lock()
				updateBusinessRules(v, vv)
				v.Lock.Unlock()
				break
			}
		}
		if delFlag {
			v.destroy()
		}else{
			newBmj = append(newBmj, v)
		}
	}
	businessMonitorJobs = newBmj
	for _,v := range param {
		addFlag := true
		for _,vv := range businessMonitorJobs {
			if vv.Path == v.Path {
				addFlag = false
				break
			}
		}
		if addFlag {
			newBmo := businessMonitorObj{}
			newBmo.Path = v.Path
			newBmo.Lock = new(sync.RWMutex)
			for _,vv := range v.Config {
				tmpRuleObj := businessRuleObj{}
				tmpRuleObj.StringMap = vv.StringMap
				tmpRuleObj.MetricConfig = vv.MetricConfig
				tmpRuleObj.Regular = vv.Regular
				tmpRuleObj.RegExp = transBusinessRegular(vv.Regular)
				var tmpTagsKey,tmpTagsValue []string
				for _,tmpKey := range strings.Split(vv.Tags, ",") {
					tmpTagsKey = append(tmpTagsKey, tmpKey)
					tmpTagsValue = append(tmpTagsValue, "")
				}
				tmpRuleObj.TagsKey = tmpTagsKey
				tmpRuleObj.TagsValue = tmpTagsValue
				tmpRuleObj.DataChannel = make(chan map[string]interface{}, 10000)
				newBmo.Rules = append(newBmo.Rules, &tmpRuleObj)
			}
			go newBmo.start()
			businessMonitorJobs = append(businessMonitorJobs, &newBmo)
		}
	}
	businessMonitorLock.Unlock()
	level.Info(newLogger).Log("msg","success")
	w.Write([]byte("success"))
}

func updateBusinessRules(bmo *businessMonitorObj,config  *businessAgentDto)  {
	var newRules []*businessRuleObj
	for _,v := range bmo.Rules {
		delFlag := true
		for _,vv := range config.Config {
			if vv.Regular == v.Regular {
				delFlag = false
				v.StringMap = vv.StringMap
				v.MetricConfig = vv.MetricConfig
				var newTagsKey,newTagsValue []string
				for _,cfgKey := range strings.Split(vv.Tags, ",") {
					newTagsKey = append(newTagsKey, cfgKey)
					keyExistFlag := false
					for existKeyIndex,existKey := range v.TagsKey {
						if existKey == cfgKey {
							keyExistFlag = true
							newTagsValue = append(newTagsValue, v.TagsValue[existKeyIndex])
							break
						}
					}
					if !keyExistFlag {
						newTagsValue = append(newTagsValue, "")
					}
				}
				v.TagsKey = newTagsKey
				v.TagsValue = newTagsValue
				break
			}
		}
		if !delFlag {
			newRules = append(newRules, v)
		}
	}
	for _,v := range config.Config {
		addFlag := true
		for _,vv := range newRules {
			if v.Regular == vv.Regular {
				addFlag = false
				break
			}
		}
		if addFlag {
			tmpRuleObj := businessRuleObj{}
			tmpRuleObj.StringMap = v.StringMap
			tmpRuleObj.MetricConfig = v.MetricConfig
			tmpRuleObj.Regular = v.Regular
			tmpRuleObj.RegExp = transBusinessRegular(v.Regular)
			var tmpTagsKey,tmpTagsValue []string
			for _,tmpKey := range strings.Split(v.Tags, ",") {
				tmpTagsKey = append(tmpTagsKey, tmpKey)
				tmpTagsValue = append(tmpTagsValue, "")
			}
			tmpRuleObj.TagsKey = tmpTagsKey
			tmpRuleObj.TagsValue = tmpTagsValue
			tmpRuleObj.DataChannel = make(chan map[string]interface{}, 10000)
			newRules = append(newRules, &tmpRuleObj)
		}
	}
	printByte,_ := json.Marshal(newRules)
	level.Info(newLogger).Log("updateBusinessRules",string(printByte))
	bmo.Rules = newRules
}

type businessCollectorStore struct {
	Data  []*businessStoreMonitorObj  `json:"data"`
}

var BusinessCollectorStore businessCollectorStore

func (c *businessCollectorStore) Save()  {
	for _,v := range businessMonitorJobs {
		var newStoreRules []*businessStoreMetricObj
		for _,vv := range v.Rules {
			newStoreRules = append(newStoreRules, &businessStoreMetricObj{Regular: vv.Regular, StringMap: vv.StringMap, MetricConfig: vv.MetricConfig, TagsKey: vv.TagsKey, TagsValue: vv.TagsValue})
		}
		c.Data = append(c.Data, &businessStoreMonitorObj{Path: v.Path,Rules: newStoreRules})
	}
	var tmpBuffer bytes.Buffer
	enc := gob.NewEncoder(&tmpBuffer)
	err := enc.Encode(c.Data)
	if err != nil {
		level.Error(newLogger).Log("msg",fmt.Sprintf("gob encode business monitor error : %v ", err))
	}else{
		ioutil.WriteFile(businessMonitorFilePath, tmpBuffer.Bytes(), 0644)
		level.Info(newLogger).Log("msg",fmt.Sprintf("write %s succeed ", businessMonitorFilePath))
	}
}

func (c *businessCollectorStore) Load()  {
	file,err := os.Open(businessMonitorFilePath)
	if err != nil {
		level.Info(newLogger).Log("msg",fmt.Sprintf("read %s file error %v ", businessMonitorFilePath, err))
	}else{
		dec := gob.NewDecoder(file)
		err = dec.Decode(&c.Data)
		if err != nil {
			level.Error(newLogger).Log("msg",fmt.Sprintf("gob decode %s error %v ", businessMonitorFilePath, err))
		}else{
			level.Info(newLogger).Log("msg",fmt.Sprintf("load %s file succeed ", businessMonitorFilePath))
		}
	}
	businessMonitorLock.Lock()
	businessMonitorJobs = []*businessMonitorObj{}
	for _,v := range c.Data {
		if v.Path != "" {
			newBusinessMonitorObj := businessMonitorObj{Path: v.Path}
			newBusinessMonitorObj.Lock = new(sync.RWMutex)
			for _,vv := range v.Rules {
				tmpRuleObj := businessRuleObj{Regular: vv.Regular, MetricConfig: vv.MetricConfig, StringMap: vv.StringMap, TagsKey: vv.TagsKey, TagsValue: vv.TagsValue}
				tmpRuleObj.RegExp = transBusinessRegular(vv.Regular)
				tmpRuleObj.DataChannel = make(chan map[string]interface{}, 10000)
				newBusinessMonitorObj.Rules = append(newBusinessMonitorObj.Rules, &tmpRuleObj)
			}
			businessMonitorJobs = append(businessMonitorJobs, &newBusinessMonitorObj)
		}
	}
	for _,v := range businessMonitorJobs {
		go v.start()
	}
	businessMonitorLock.Unlock()
}

func transBusinessRegular(regRuleString string) *regexp.Regexp {
	regRuleString = strings.ReplaceAll(regRuleString, "[", "\\[")
	regRuleString = strings.ReplaceAll(regRuleString, "]", "\\]")
	regRuleString = strings.ReplaceAll(regRuleString, "${json_content}", "(.*)")
	return regexp.MustCompile(regRuleString)
}

func StartBusinessAggCron()  {
	t := time.NewTicker(10*time.Second).C
	for {
		<- t
		go calcBusinessAggData()
	}
}

func calcBusinessAggData()  {
	var newRuleData []*businessRuleMetricObj
	businessMonitorLock.RLock()
	for _,v := range businessMonitorJobs {
		for _,rule := range v.Rules {
			dataLength := len(rule.DataChannel)
			if dataLength == 0 {
				break
			}
			var sum,avg,count []float64
			for i:=0;i<len(rule.MetricConfig);i++ {
				sum = append(sum, 0)
				count = append(count, 0)
			}
			for i:=0;i<dataLength;i++ {
				tmpMapData := <- rule.DataChannel
				if i == 0 {
					for tagIndex,tagKey := range rule.TagsKey {
						if tmpTagValue,b:=tmpMapData[tagKey];b {
							rule.TagsValue[tagIndex] = fmt.Sprintf("%s", tmpTagValue)
						}
					}
				}
				for metricIndex,metricConfig := range rule.MetricConfig {
					if metricValue,b:=tmpMapData[metricConfig.Key];b {
						metricValueString := fmt.Sprintf("%s", metricValue)
						metricValueFloat,parseError := strconv.ParseFloat(metricValueString,64)
						if parseError != nil {
							for _,tmpStringMapObj := range rule.StringMap {
								if tmpStringMapObj.Key == metricConfig.Key && tmpStringMapObj.StringValue == metricValueString {
									metricValueFloat = tmpStringMapObj.IntValue
									break
								}
							}
						}
						sum[metricIndex] += metricValueFloat
						count[metricIndex]++
					}
				}
			}
			for i,tmpSum := range sum {
				avg = append(avg, tmpSum/count[i])
			}
			var tagStringContent string
			for tmpTagIndex,tmpTags := range rule.TagsKey {
				tagStringContent += fmt.Sprintf("%s=%s", tmpTags, rule.TagsValue[tmpTagIndex])
				if tmpTagIndex < len(rule.TagsKey)-1 {
					tagStringContent += ","
				}
			}
			for metricIndex,metricConfig := range rule.MetricConfig {
				tmpMetricObj := businessRuleMetricObj{Path: v.Path, Agg: metricConfig.AggType, TagsString: tagStringContent}
				tmpMetricObj.Metric = metricConfig.Metric
				if metricConfig.AggType == "sum" {
					tmpMetricObj.Value = sum[metricIndex]
				}else if metricConfig.AggType == "avg" {
					tmpMetricObj.Value = avg[metricIndex]
				}else if metricConfig.AggType == "count" {
					tmpMetricObj.Value = count[metricIndex]
				}
				newRuleData = append(newRuleData, &tmpMetricObj)
			}
		}
	}
	businessMonitorLock.RUnlock()
	businessMonitorMetricLock.Lock()
	businessMonitorMetrics = newRuleData
	businessMonitorMetricLock.Unlock()
}