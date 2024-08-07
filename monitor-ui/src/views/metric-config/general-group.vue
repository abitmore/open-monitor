<template>
  <div ref="maxheight" class="monitor-general-group">
    <Row>
      <Col :span="8">
      <!--对象类型-->
      <span style="font-size: 14px;">
        {{$t('m_basic_type')}}：
      </span>
      <Select filterable v-model="monitorType" @on-change="changeMonitorType" style="width:300px">
        <Option v-for="(i, index) in monitorTypeOptions" :value="i" :key="index">{{ i }}</Option>
      </Select>
      </Col>
      <Col :span="16">
      <div class="btn-group">
        <!--新增-->
        <Button type="success" @click="handleAdd">{{ $t('m_button_add') }}</Button>
      </div>
      </Col>
    </Row>
    <Table size="small" :columns="tableColumns.filter(col=>col.showType.includes(metricType))" :data="tableData" class="general-table"/>
    <Modal
      v-model="deleteVisible"
      :title="$t('m_delConfirm_title')"
      @on-ok="submitDelete"
      @on-cancel="deleteVisible = false"
    >
      <div class="modal-body" style="padding: 10px 20px">
        <p style="color: red">{{ $t('m_metric_deleteTips') }}</p>
      </div>
    </Modal>
    <AddGroupDrawer
      v-if="addVisible && showDrawer === 'originalMetrics'"
      :visible.sync="addVisible"
      :monitorType="monitorType"
      :data="row"
      :operator="type"
      :viewOnly="viewOnly"
      fromPage="generalGroup"
      @fetchList="reloadData('originalMetrics')"
    ></AddGroupDrawer>
    <YearOverYear
      ref="yearOverYearRef"
      v-if="addVisible && showDrawer === 'comparisonMetrics'"
      :visible.sync="addVisible"
      :monitorType="monitorType"
      :originalMetricsId="originalMetricsId"
      serviceGroup=""
      :data="row"
      :operator="type"
      :viewOnly="viewOnly"
      @fetchList="reloadData('comparisonMetrics')"
    ></YearOverYear>
  </div>
</template>

<script>
import AddGroupDrawer from './components/add-group.vue'
import YearOverYear from './components/year-over-year.vue'
export default {
  components: {
    AddGroupDrawer,
    YearOverYear
  },
  data() {
    return {
      metricType: 'originalMetrics', // 原始指标originalMetrics、同环比指标comparisonMetrics
      monitorType: '',
      serviceGroup: '',
      monitorTypeOptions: [],
      maxHeight: 500,
      tableData: [],
      tableColumns: [
        {
          title: this.$t('m_year_over_year_metrics'), // 指标
          key: 'metric',
          minWidth: 250,
          showType: ['comparisonMetrics']
        },
        {
          title: this.$t('m_metric'), // 原始指标-指标
          key: 'metric',
          minWidth: 200,
          showType: ['originalMetrics']
        },
        {
          title: this.$t('m_metric'), // 同环比指标-指标
          key: 'origin_metric',
          minWidth: 200,
          showType: ['comparisonMetrics']
        },
        {
          title: this.$t('m_scope'), // 作用域
          key: 'workspace',
          width: 160,
          render: (h, params) => {
            const workspaceMap = {
              all_object: this.$t('m_all_object'), // 全部对象
              any_object: this.$t('m_any_object') // 层级对象
            }
            return <Tag size="medium">{workspaceMap[params.row.workspace] || '-'}</Tag>
          },
          showType: ['originalMetrics', 'comparisonMetrics']
        },
        {
          title: this.$t('m_configuration_page'), // 配置页面
          key: 'metric_type',
          width: 140,
          render: (h, params) => {
            const typeList = [
              {
                label: this.$t('m_basic_type'),
                value: 'common',
                color: '#2d8cf0'
              },
              {
                label: this.$t('m_business_configuration'),
                value: 'business',
                color: '#81b337'
              },
              {
                label: this.$t('m_metric_list'),
                value: 'custom',
                color: '#b886f8'
              }
            ]
            const find = typeList.find(item => item.value === params.row.metric_type) || {}
            return <Tag size="medium" type="border" color={find.color}>{find.label || '-'}</Tag>
          },
          showType: ['originalMetrics', 'comparisonMetrics']
        },
        {
          title: this.$t('m_tableKey_expr'), // 表达式
          key: 'prom_expr',
          minWidth: 400,
          render: (h, params) => (
            <Tooltip max-width="300" content={params.row.prom_expr} transfer>
              <span class="eclipse">{params.row.prom_expr || '-'}</span>
            </Tooltip>
          ),
          showType: ['originalMetrics']
        },
        {
          title: this.$t('m_comparison_types'), // 对比类型
          key: 'comparisonType',
          width: 100,
          render: (h, params) => {
            const comparisonTypeToDisplay = {
              day: this.$t('m_dod_comparison'),
              week: this.$t('m_wow_comparison'),
              month: this.$t('m_mom_comparison'),
              '-': '-'
            }
            return <span>{comparisonTypeToDisplay[params.row.comparisonType || '-']}</span>
          },
          showType: ['comparisonMetrics']
        },
        {
          title: this.$t('m_calc_value'), // 计算数值
          key: 'calcType',
          width: 160,
          render: (h, params) => {
            const calcTypeToDisplay = {
              diff: this.$t('m_difference'),
              diff_percent: this.$t('m_percentage_difference')
            }
            const calcTypeCache = (params.row.calcType || []).map(t => calcTypeToDisplay[t]).join('，')
            return <span>{calcTypeCache || '-'}</span>
          },
          showType: ['comparisonMetrics']
        },
        {
          title: this.$t('m_calc_method'), // 计算方法
          key: 'calcMethod',
          width: 100,
          render: (h, params) => {
            const calcMethodToDisplay = {
              avg: this.$t('m_average'),
              min: this.$t('m_min'),
              max: this.$t('m_max'),
              sum: this.$t('m_sum'),
              '-': '-'
            }
            return <span>{calcMethodToDisplay[params.row.calcMethod || '-']}</span>
          },
          showType: ['comparisonMetrics']
        },
        {
          title: this.$t('m_calculation_period'), // 计算周期
          key: 'calcPeriod',
          width: 100,
          render: (h, params) => <span>{params.row.calcPeriod || '-'}S</span>,
          showType: ['comparisonMetrics']
        },
        {
          title: this.$t('m_update_time'), // 更新时间
          key: 'update_time',
          width: 160,
          render: (h, params) => <span>{params.row.update_time || '-'}</span>,
          showType: ['originalMetrics', 'comparisonMetrics']
        },
        {
          title: this.$t('m_updatedBy'), // 更新人
          key: 'update_user',
          width: 160,
          render: (h, params) => <span>{params.row.update_user || '-'}</span>,
          showType: ['originalMetrics', 'comparisonMetrics']
        },
        {
          title: this.$t('m_table_action'),
          key: 'action',
          width: 170,
          showType: ['originalMetrics', 'comparisonMetrics'],
          fixed: 'right',
          render: (h, params) => (
            <div style="display:flex;justify-content:center;">
              {
                this.metricType === 'originalMetrics'
                  /* 新增同环比指标 */
                  && <Tooltip content={this.$t('m_button_add')+this.$t('m_year_over_year_metrics')} placement="bottom" transfer>
                    <Button
                      size="small"
                      type="success"
                      onClick={() => {
                        this.handleAddYearOverYear(params.row)
                      }}
                      style="margin-right:5px;"
                    >
                      <Icon type="md-add" size="16"></Icon>
                    </Button>
                  </Tooltip>
              }
              {
                /* 查看 */
                <Tooltip content={this.$t('m_button_view')} placement="top" transfer={true}>
                  <Button
                    size="small"
                    type="info"
                    onClick={() => this.handleEdit(params.row, true)}
                    style="margin-right:5px;"
                  >
                    <Icon type="ios-eye" size="16"></Icon>
                  </Button>
                </Tooltip>
              }
              {
                /* 编辑 */
                <Tooltip content={this.$t('m_button_edit')} placement="bottom" transfer>
                  <Button
                    size="small"
                    type="primary"
                    onClick={() => {
                      this.handleEdit(params.row, false)
                    }}
                    style="margin-right:5px;"
                  >
                    <Icon type="md-create" size="16"></Icon>
                  </Button>
                </Tooltip>
              }
              {
                /* 删除 */
                <Tooltip content={this.$t('m_button_remove')} placement="bottom" transfer>
                  <Button
                    size="small"
                    type="error"
                    onClick={() => {
                      this.handleDelete(params.row)
                    }}
                    style="margin-right:5px;"
                  >
                    <Icon type="md-trash" size="16"></Icon>
                  </Button>
                </Tooltip>
              }
            </div>
          )
        }
      ],
      row: {},
      type: '', // add、edit
      addVisible: false,
      deleteVisible: false,
      originalMetricsId: '',
      showDrawer: '', // 控制显示抽屉的类型
      viewOnly: false, // 仅查看
    }
  },
  mounted() {
    this.getMonitorType()
    const clientHeight = document.documentElement.clientHeight
    this.maxHeight = clientHeight - this.$refs.maxheight.getBoundingClientRect().top - 100
  },
  methods: {
    reloadData(metricType) {
      this.metricType = metricType
      this.getList()
    },
    getList() {
      const params = {
        monitorType: this.monitorType,
        onlyService: 'Y',
        serviceGroup: this.serviceGroup
      }
      const api = this.metricType === 'originalMetrics' ? '/monitor/api/v2/monitor/metric/list' : '/monitor/api/v2/monitor/metric_comparison/list'
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', api, params, responseData => {
        this.tableData = responseData
      }, {isNeedloading: true})
      this.getMetricTotalNumber()
    },
    getMetricTotalNumber() {
      const params = {
        monitorType: this.monitorType,
        onlyService: 'Y',
        serviceGroup: this.serviceGroup
      }
      const api = '/monitor/api/v2/monitor/metric/list/count'
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', api, params, response => {
        this.$emit('totalCount', response, this.metricType)
      }, {isNeedloading: false})
    },
    getMonitorType() {
      this.$root.$httpRequestEntrance.httpRequestEntrance('GET', this.$root.apiCenter.getEndpointType, '', responseData => {
        this.monitorTypeOptions = responseData || []
        this.monitorType = this.monitorTypeOptions[0]
        this.getList()
      }, {isNeedloading: false})
    },
    changeMonitorType() {
      this.getList()
    },
    handleAdd() {
      this.type = 'add'
      this.viewOnly = false
      this.showDrawer = this.metricType
      this.addVisible = true
      this.originalMetricsId = ''
    },
    handleEdit(row, viewOnly) {
      this.showDrawer = this.metricType
      this.type = 'edit'
      this.viewOnly = viewOnly
      this.row = row
      this.addVisible = true
    },
    handleDelete(row) {
      this.row = row
      this.deleteVisible = true
    },
    submitDelete() {
      const api = this.metricType === 'originalMetrics' ? `${this.$root.apiCenter.metricManagement}?id=${this.row.guid}` : `/monitor/api/v1/dashboard/new/comparison_metric/${this.row.guid}`
      this.$root.$httpRequestEntrance.httpRequestEntrance(
        'DELETE',
        api,
        '',
        () => {
          this.$Message.success(this.$t('m_tips_success'))
          this.getList()
        }
      )
    },
    // 同环比入口
    handleAddYearOverYear(row) {
      this.showDrawer = 'comparisonMetrics'
      this.addVisible = true
      this.type = 'add'
      this.viewOnly = false
      this.originalMetricsId = row.guid
    },
  }
}
</script>

<style lang="less">
.monitor-general-group {
  .eclipse {
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }
}
</style>
<style lang="less" scoped>
.monitor-general-group {
  padding-bottom: 20px;
  .btn-group {
    display: flex;
    justify-content: flex-end;
  }
  .general-table {
    margin-top: 12px;
  }
}
</style>
