<template>
  <div>
    <Row>
      <Col span="8" :style="{height: showTemplate ? '' : '510px'}" style="overflow: auto;">
      <Form :label-width="120" style="margin-top:12px">
        <FormItem :label="$t('m_template_name')">
          <Tooltip :content="configInfo.name" transfer :disabled="configInfo.name === ''" style="width: 100%;" max-width="200">
            <Input
              v-model="configInfo.name"
              maxlength="30"
              show-word-limit
              style="width: 96%"
              disabled
            />
            <span style="color: red">*</span>
          </Tooltip>
        </FormItem>
        <!-- <FormItem>
            <Button type="primary" @click="showTemplate = !showTemplate" ghost size="small" >{{showTemplate ? $t('m_hide_template'):$t('m_expand_template')}}</Button>
          </FormItem> -->
        <template v-if="showTemplate === true" >
          <FormItem :label="$t('m_updatedBy')">
            {{ configInfo.update_user }}
          </FormItem>
          <FormItem :label="$t('m_title_updateTime')">
            {{ configInfo.update_time }}
          </FormItem>
        </template>
        <FormItem  v-if="showTemplate === true" :label="$t('m_json_regular')" style="margin-bottom: 12px;">
          <Input
            v-model="configInfo.json_regular"
            maxlength="200"
            show-word-limit
            type="textarea"
            style="width: 96%"
            disabled
          />
          <div v-if="isParmasChanged && configInfo.json_regular.length > 200" style="color: red">
            {{ $t('m_json_regular') }}{{ $t('tw_limit_200') }}
          </div>
        </FormItem>
        <FormItem  v-if="showTemplate === true" :label="$t('m_log_example')">
          <Input
            v-model="configInfo.demo_log"
            type="textarea"
            :rows="12"
            style="width: 96%"
            disabled
          />
          <div v-if="isParmasChanged && configInfo.demo_log.length === 0" style="color: red">
            {{ $t('m_log_example') }} {{ $t('m_tips_required') }}
          </div>
        </FormItem>
        <FormItem>
          <!-- <Button type="primary" @click="confirmGenerateBackstageTrial" ghost size="small" style="float:right;margin:12px" :disabled="configInfo.demo_log===''||configInfo.json_regular===''">{{ $t('m_match') }}</Button> -->
        </FormItem>
        <FormItem  v-if="showTemplate === true" :label="$t('m_matching_result')" style="margin-top: 12px;">
          <Input
            disabled
            v-model="configInfo.calc_result.match_text"
            type="textarea"
            :rows="6"
            style="width: 96%"
          />
        </FormItem>
      </Form>
      </Col>
      <Col span="16" style="border-left: 2px solid rgb(232 234 236)">
      <div style="margin-left: 8px;">
        <!-- 采集参数 -->
        <div v-if="showTemplate === true">
          <Divider orientation="left" size="small">{{ $t('m_parameter_collection') }}</Divider>
          <Table
            style="position: inherit;"
            size="small"
            :columns="columnsForParameterCollection"
            :data="configInfo.param_list"
            width="100%"
          ></Table>
        </div>
        <!-- 计算指标 -->
        <div>
          <Divider orientation="left" size="small">{{ $t('m_compute_metrics') }}</Divider>
          <Table
            size="small"
            :columns="columnsForComputeMetrics"
            :data="configInfo.metric_list"
            width="100%"
          ></Table>

        </div>
      </div>
      </Col>
    </Row>
  </div>
</template>

<script>
export default {
  name: 'json-regex',
  data() {
    return {
      isParmasChanged: false,
      isAdd: false,
      showTemplate: false,
      columnsForParameterCollection: [
        {
          title: this.$t('m_field_displayName'),
          key: 'display_name',
          width: 120
        },
        {
          title: this.$t('m_parameter_key'),
          key: 'name',
          width: 140
        },
        {
          title: this.$t('m_json_key'),
          key: 'json_key',
          renderHeader: () => (
            <span>
              <span style="color:red">*</span>
              <span>{this.$t('json_key')}</span>
            </span>
          ),
          render: (h, params) => {
            const selectOptions = this.configInfo.calc_result.json_key_list
            return (
              <Select
                style="'z-index: 1000'"
                disabled
                value={params.row.json_key}
              >
                {selectOptions.map(option => (
                  <Option key={option} value={option}>
                    {option}
                  </Option>
                ))}
              </Select>
            )
          }
        },
        {
          title: this.$t('m_matching_result'),
          ellipsis: true,
          tooltip: true,
          renderHeader: () => (
            <span>
              <span style="color:red">*</span>
              <span>{this.$t('m_matching_result')}</span>
            </span>
          ),
          key: 'demo_match_value',
          render: (h, params) => {
            const demo_match_value = params.row.demo_match_value
            return (
              <Tooltip content={demo_match_value} max-width="300" >
                <span style={demo_match_value?'':'color:#c5c8ce'}>{demo_match_value || this.$t('m_no_matching')}</span>
              </Tooltip>
            )
          }
        },
      ],
      columnsForComputeMetrics: [
        {
          title: this.$t('m_field_displayName'),
          key: 'display_name',
          width: 120
        },
        {
          title: this.$t('m_metric_key'),
          key: 'metric',
          width: 140,
          render: (h, params) => (
            <span>
              <span>{this.prefixCode}_{params.row.metric}</span>
            </span>
          )
        },
        {
          title: this.$t('m_statistical_parameters'),
          key: 'log_param_name',
        },
        {
          title: this.$t('m_filter_label'),
          key: 'tag_config',
          render: (h, params) => (
            <span>
              {params.row.tag_config.join(',')}
            </span>
          )
        },
        {
          title: this.$t('m_computed_type'),
          key: 'agg_type',
          render: (h, params) => {
            const agg_type = params.row.agg_type
            return (
              <Tooltip content={agg_type} max-width="300" >
                <span>{agg_type}</span>
              </Tooltip>
            )
          }
        }
      ],
      generateBackstageTrialWarning: false,
      fullScreenHeight: '100px'
    }
  },
  props: {
    configInfo: Object,
    prefixCode: String
  },
  mounted() {},
  methods: {
    hideTemplate() {
      this.showTemplate = false
      this.isfullscreen = false
    },
    changeTemplateStatus() {
      this.showTemplate = !this.showTemplate
    },
    returnCurrentStatus() {
      return this.showTemplate
    }
  }
}
</script>

<style lang="less" scoped>
.custom-modal-header {
  line-height: 20px;
  font-size: 16px;
  color: #17233d;
  font-weight: 500;
  .fullscreen-icon {
    float: right;
    margin-right: 28px;
    font-size: 18px;
    cursor: pointer;
  }
}
.ivu-form-item {
  margin-bottom: 0px;
}
</style>
