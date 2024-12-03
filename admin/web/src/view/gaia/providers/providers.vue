<template>
  <div>
    <div v-if="false" class="gva-search-box">
      <el-form
        ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule"
        @keyup.enter="onSubmit"
      >
        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>
        <el-form-item label="工作区">
          <el-select
            v-model="searchInfo.tenant_id"
            clearable
            placeholder="请选择"
            disabled
          >
            <el-option
              v-for="item in tenantOptions"
              :key="item.value"
              :label="`${item.label}`"
              :value="item.value"
            >
              <span style="float: left">{{ item.label }}</span>
              <span
                style="
          float: right;
          color: var(--el-text-color-secondary);
          font-size: 13px;
        "
              >
                {{ item.value }}
              </span>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <!--      <div class="gva-btn-list">-->
      <!--        <el-button type="primary" icon="plus" @click="openDialog">-->
      <!--          新增-->
      <!--        </el-button>-->
      <!--        <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="onDelete">-->
      <!--          删除-->
      <!--        </el-button>-->
      <!--      </div>-->
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="id"
      >
        <el-table-column align="left" label="id" prop="id" width="300" show-overflow-tooltip />
        <el-table-column align="left" label="工作区名称" prop="tenant_id" width="120">
          <template #default="scope">
            {{ filterDict(scope.row.tenant_id, tenantOptions) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="模型供应商" prop="provider_name" />
        <el-table-column align="left" label="模型名称" prop="provider_model_name" />
        <el-table-column align="left" label="模型类型" prop="model_type" width="120" />
        <el-table-column align="left" label="同步工作区数" width="120">
          <template #default="scope">
            {{ scope.row.tenant_ids ? scope.row.tenant_ids.length : 0 }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="是否同步所有" prop="is_all" width="120">
          <template #default="scope">
            {{ formatBoolean(scope.row.is_all) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px">
                <InfoFilled />
              </el-icon>
              查看
            </el-button>
            <el-button type="primary" link icon="Compass" class="table-button" @click="syncProvidersFunc(scope.row)">
              同步
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
    <el-drawer v-model="dialogFormVisible" destroy-on-close size="600" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">模型同步工作空间</span>
          <div>
            <el-button type="primary" @click="enterDialog">
              同 步
            </el-button>
            <el-button @click="closeDialog">
              取 消
            </el-button>
          </div>
        </div>
      </template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="id">
          {{ formData.id }}
        </el-descriptions-item>
        <el-descriptions-item label="所在工作区">
          {{ filterDict(formData.tenant_id, tenantOptions) + "（" + formData.tenant_id + "）" }}
        </el-descriptions-item>
        <el-descriptions-item label="模型供应商">
          {{ formData.provider_name }}
        </el-descriptions-item>
        <el-descriptions-item label="模型名称">
          {{ formData.provider_model_name }}
        </el-descriptions-item>
        <el-descriptions-item label="模型类型">
          {{ formData.model_type }}
        </el-descriptions-item>
      </el-descriptions>
      <el-form ref="elFormRef" :model="formData" label-position="left" :rules="rule" label-width="120px" style="margin-top: 20px">
        <el-form-item label="是否同步所有:" prop="is_all">
          <el-switch
            v-model="formData.is_all" active-color="#13ce66" inactive-color="#ff4949" active-text="是"
            inactive-text="否" clearable
          />
        </el-form-item>

        <el-form-item v-show="!formData.is_all" label="已同步工作空间:" prop="tenant_ids">
          <el-select
            v-model="formData.tenant_ids"
            clearable
            placeholder="请选择"
            multiple
            filterable
            collapse-tags
            collapse-tags-tooltip
          >
            <el-option
              v-for="item in tenantOptions"
              :key="item.value"
              :label="`${item.label}`"
              :value="item.value"
              :disabled="item.value === formData.tenant_id"
            >
              <span style="float: left">{{ item.label }}</span>
              <span
                style="
          float: right;
          color: var(--el-text-color-secondary);
          font-size: 13px;
        "
              >
                {{ item.value }}
              </span>
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer
      v-model="detailShow" destroy-on-close size="800" :show-close="true" :before-close="closeDetailShow"
      title="查看"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="id">
          {{ detailFrom.id }}
        </el-descriptions-item>
        <el-descriptions-item label="所在工作区">
          {{ filterDict(detailFrom.tenant_id, tenantOptions) + "（" + detailFrom.tenant_id + "）" }}
        </el-descriptions-item>
        <el-descriptions-item label="模型供应商">
          {{ detailFrom.provider_name }}
        </el-descriptions-item>
        <el-descriptions-item label="模型名称">
          {{ detailFrom.provider_model_name }}
        </el-descriptions-item>
        <el-descriptions-item label="模型类型">
          {{ detailFrom.model_type }}
        </el-descriptions-item>
        <el-descriptions-item label="是否同步所有">
          {{ detailFrom.is_all ? '是' : '否' }}
        </el-descriptions-item>
        <el-descriptions-item label="已同步工作区数">
          {{ detailFrom.tenant_ids ? detailFrom.tenant_ids.length :'0' }}
        </el-descriptions-item>
        <el-descriptions-item label="已同步工作空间" width="10">
          <!--  列表展示, 如果没有显示暂无-->
          <span v-if="!detailFrom.tenant_ids">暂无</span>
          <template v-else>
            <el-tag v-for="item in detailFrom.tenant_ids" :key="item" type="info" class="mr-2">
              {{ filterDict(item, tenantOptions) }}
            </el-tag>
          </template>
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  syncProviders,
  findProviders,
  getProvidersList
} from '@/api/gaia/providers'

// 全量引入格式化工具 请按需保留
import {
  formatBoolean,
  filterDict,
} from '@/utils/format'
import {ElMessage} from 'element-plus'
import {ref, reactive} from 'vue'
import {getAllTenants} from "@/api/gaia/tenants";


defineOptions({
  name: "ProviderManage"
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

// 自动化生成的字典（可能为空）以及
const formData = ref({
  id: '',
  tenant_id: '',
  provider_name: '',
  provider_model_name: '',
  model_type: '',
  tenant_ids: [],
  is_all: false,
})


// 验证规则
const rule = reactive({})

const searchRule = reactive({
  createdAt: [
    {
      validator: (rule, value, callback) => {
        if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
          callback(new Error('请填写结束日期'))
        } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
          callback(new Error('请填写开始日期'))
        } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
          callback(new Error('开始日期应当早于结束日期'))
        } else {
          callback()
        }
      }, trigger: 'change'
    }
  ],
})

const elFormRef = ref()
const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({
  // tenant_id: "f8043947-fba6-48b4-9a2c-040aad6b41d9", // TODO：默认工作区
})

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async (valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
    if (searchInfo.value.is_all === "") {
      searchInfo.value.is_all = null
    }
    getTableData()
  })
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async () => {
  const table = await getProvidersList({page: page.value, pageSize: pageSize.value, ...searchInfo.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () => {
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 更新行
const syncProvidersFunc = async (row) => {
  const res = await findProviders({id: row.id})
  type.value = 'update'
  if (res.code === 0) {
    formData.value = res.data
    dialogFormVisible.value = true
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  formData.value = {
    id: '',
    tenant_id: '',
    provider_name: '',
    provider_model_name: '',
    model_type: '',
    tenant_ids: [],
    is_all: false,
  }
}
// 弹窗确定
const enterDialog = async () => {
  elFormRef.value?.validate(async (valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'update':
        res = await syncProviders(formData.value)
        break
      default:
        res = await syncProviders(formData.value)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '同步成功'
      })
      closeDialog()
      getTableData()
    }
  })
}


const detailFrom = ref({})

// 查看详情控制标记
const detailShow = ref(false)


// 打开详情弹窗
const openDetailShow = () => {
  detailShow.value = true
}


// 打开详情
const getDetails = async (row) => {
  // 打开弹窗
  const res = await findProviders({id: row.id})
  if (res.code === 0) {
    detailFrom.value = res.data
    openDetailShow()
  }
}


// 关闭详情弹窗
const closeDetailShow = () => {
  detailShow.value = false
  detailFrom.value = {}
}

const tenantOptions = ref([
  {
    value: 'f8043947-fba6-48b4-9a2c-040aad6b41d9',
    label: '默认空间'
  },
])

const getAllTenantData = async () => {
  // 打开弹窗
  const res = await getAllTenants()
  if (res.code === 0) {
    tenantOptions.value = []
    for (const item of res.data){
      console.log(item.id)
      console.log(item.name)
      tenantOptions.value.push({
        value: item.id,
        label: item.name
      })
    }
  }
}
getAllTenantData()

</script>

<style>

</style>
