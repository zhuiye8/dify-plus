<!--
    @auther: bypanghu<bypanghu@163.com>
    @date: 2024/5/8
!-->

<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column prop="ranking" label="排名" width="80" align="center" />
      <el-table-column prop="name" label="应用名称" width="230" show-overflow-tooltip />
<!--      <el-table-column prop="message_cost" label="对话记录费用($)">-->
<!--        <template #default="scope">-->
<!--          <el-tooltip :content="scope.row.message_cost">-->
<!--            <span>{{ truncateToOneDecimal(scope.row.message_cost) }}</span>-->
<!--          </el-tooltip>-->
<!--        </template>-->
<!--      </el-table-column>-->
<!--      <el-table-column prop="workflow_cost" label="工作流记录费用($)">-->
<!--        <template #default="scope">-->
<!--          <el-tooltip :content="scope.row.workflow_cost">-->
<!--            <span>{{ truncateToOneDecimal(scope.row.workflow_cost) }}</span>-->
<!--          </el-tooltip>-->
<!--        </template>-->
<!--      </el-table-column>-->
      <el-table-column prop="total_cost" label="总费用($)" width="120">
        <template #default="scope">
          <el-tooltip :content="scope.row.total_cost">
            <span>{{ truncateToOneDecimal(scope.row.total_cost) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="record_num" label="生成记录数" width="95" />
      <el-table-column label="平均每条($)" width="100">
        <template #default="scope">
          <el-tooltip :content="(scope.row.total_cost / scope.row.record_num).toFixed(5)">
            <span>{{ (scope.row.total_cost / scope.row.record_num).toFixed(5) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="use_num" label="界面使用数" width="120" />
      <el-table-column prop="mode" label="应用类型" width="150">
        <template #default="scope">
          <el-tooltip :content="scope.row.mode">
            <span :style="{ color: getAppModeColor(scope.row.mode) }">
              {{ getAppModeText(scope.row.mode) }}
            </span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="tenant_name" label="工作区名称" width="170" show-overflow-tooltip />
      <el-table-column prop="account_name" label="账号名称" />
    </el-table>

    <div class="gva-pagination-gaia">
      <el-pagination
        size="small"
        :page-size="pageSize"
        pager-count="5"
        layout="prev, pager, next"
        :total="total"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue";
import {getAppQuotaRankingData} from "@/api/gaia/dashboard";
import {getAppModeColor, getAppModeText, truncateToOneDecimal} from "@/view/gaia/dashboard/components/index";

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}
// 查询
const getTableData = async () => {
  const table = await getAppQuotaRankingData({page: page.value, pageSize: pageSize.value})
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()
</script>

<style scoped lang="scss">
.gva-pagination-gaia {
  @apply flex justify-end;
  .el-pagination__editor {
    .el-input__inner {
      @apply h-8;
    }
  }

  .is-active {
    @apply rounded text-white;
    background: var(--el-color-primary);
    color: #ffffff !important;
  }
}

</style>
