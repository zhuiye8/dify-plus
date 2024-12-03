
<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" :rules="searchRule" @keyup.enter="onSubmit">
        <template v-if="showAllQuery">
          <!-- 将需要控制显示状态的查询条件添加到此范围内 -->
        </template>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">
            查询
          </el-button>
          <el-button icon="refresh" @click="onReset">
            重置
          </el-button>
          <el-button v-if="!showAllQuery" link type="primary" icon="arrow-down" @click="showAllQuery=true">
            展开
          </el-button>
          <el-button v-else link type="primary" icon="arrow-up" @click="showAllQuery=false">
            收起
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        
        <el-table-column align="left" label="id字段" prop="id" width="120" />
        <el-table-column align="left" label="name字段" prop="name" width="120" />
        <el-table-column align="left" label="createdAt字段" prop="createdAt" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="updatedAt字段" prop="updatedAt" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="240">
          <template #default="scope">
            <el-button type="primary" link class="table-button" @click="getDetails(scope.row)">
              <el-icon style="margin-right: 5px">
                <InfoFilled />
              </el-icon>查看
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

    <el-drawer v-model="detailShow" destroy-on-close size="800" :show-close="true" :before-close="closeDetailShow" title="查看">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="id字段">
          {{ detailFrom.id }}
        </el-descriptions-item>
        <el-descriptions-item label="name字段">
          {{ detailFrom.name }}
        </el-descriptions-item>
        <el-descriptions-item label="encryptPublicKey字段">
          {{ detailFrom.encryptPublicKey }}
        </el-descriptions-item>
        <el-descriptions-item label="plan字段">
          {{ detailFrom.plan }}
        </el-descriptions-item>
        <el-descriptions-item label="status字段">
          {{ detailFrom.status }}
        </el-descriptions-item>
        <el-descriptions-item label="createdAt字段">
          {{ detailFrom.createdAt }}
        </el-descriptions-item>
        <el-descriptions-item label="updatedAt字段">
          {{ detailFrom.updatedAt }}
        </el-descriptions-item>
        <el-descriptions-item label="customConfig字段">
          {{ detailFrom.customConfig }}
        </el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  findTenants,
  getTenantsList
} from '@/api/gaia/tenants'

// 全量引入格式化工具 请按需保留
import { formatDate } from '@/utils/format'
import { ref, reactive } from 'vue'

defineOptions({
    name: 'TenantList'
})

// 控制更多查询条件显示/隐藏状态
const showAllQuery = ref(false)

const searchRule = reactive({
  createdAt: [
    { validator: (rule, value, callback) => {
      if (searchInfo.value.startCreatedAt && !searchInfo.value.endCreatedAt) {
        callback(new Error('请填写结束日期'))
      } else if (!searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt) {
        callback(new Error('请填写开始日期'))
      } else if (searchInfo.value.startCreatedAt && searchInfo.value.endCreatedAt && (searchInfo.value.startCreatedAt.getTime() === searchInfo.value.endCreatedAt.getTime() || searchInfo.value.startCreatedAt.getTime() > searchInfo.value.endCreatedAt.getTime())) {
        callback(new Error('开始日期应当早于结束日期'))
      } else {
        callback()
      }
    }, trigger: 'change' }
  ],
})

const elSearchFormRef = ref()

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  elSearchFormRef.value?.validate(async(valid) => {
    if (!valid) return
    page.value = 1
    pageSize.value = 10
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
const getTableData = async() => {
  const table = await getTenantsList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
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
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
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
  const res = await findTenants({ id: row.id })
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

</script>

<style>

</style>