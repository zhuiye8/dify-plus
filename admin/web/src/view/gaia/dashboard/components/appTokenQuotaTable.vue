<!--
    @author: maiyouming
    @date: 2024/11/08
!-->

<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column prop="ranking" label="排名" width="80" align="center"/>
      <el-table-column prop="app_token" label="Token" width="300"/>
      <el-table-column prop="name" label="应用" width="250"/>
      <el-table-column prop="accumulated_quota" label="已使用($)">
        <template #default="scope">
          <el-tooltip :content="scope.row.accumulated_quota">
            <span>{{ truncateToOneDecimal(scope.row.accumulated_quota) }}</span>
          </el-tooltip>
        </template>
      </el-table-column>

      <el-table-column label="日使用/日限额($)">
        <template #default="scope">
          <el-tooltip :content="`${scope.row.day_used_quota} / ${scope.row.day_limit_quota}`">
            <span :class="getColorClass(scope.row.day_used_quota, scope.row.day_limit_quota)">
              {{ truncateToOneDecimal(scope.row.day_used_quota) }} /
              {{ scope.row.day_limit_quota === -1 ? '无限制' : truncateToOneDecimal(scope.row.day_limit_quota) }}
            </span>
          </el-tooltip>
        </template>
      </el-table-column>

      <el-table-column label="月使用/月限额($)">
        <template #default="scope">
          <el-tooltip :content="`${scope.row.month_used_quota} / ${scope.row.month_limit_quota}`">
            <span :class="getColorClass(scope.row.month_used_quota, scope.row.month_limit_quota)">
              {{ truncateToOneDecimal(scope.row.month_used_quota) }} /
              {{ scope.row.month_limit_quota === -1 ? '无限制' : truncateToOneDecimal(scope.row.month_limit_quota) }}
            </span>
          </el-tooltip>
        </template>
      </el-table-column>
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
import {getAppTokenQuotaRankingData} from "@/api/gaia/dashboard";
import {getColorClass, truncateToOneDecimal} from "./index";

const tableData = ref([])


const page = ref(1)
const total = ref(0)
const pageSize = ref(10)

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
  const {data, code} = await getAppTokenQuotaRankingData({page: page.value, pageSize: pageSize.value})
  if (code === 0) {
    tableData.value = data.list
    total.value = data.total
    page.value = data.page
    pageSize.value = data.pageSize
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

// 预警颜色
.green {
  color: green;
}

.yellow {
  color: #FFD700;
}

.light-red {
  color: #f84f30;
}

.deep-red {
  color: #8B0000;
}
</style>
