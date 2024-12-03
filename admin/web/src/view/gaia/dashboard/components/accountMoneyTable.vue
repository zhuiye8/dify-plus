<!--
    @auther: maiyouming
    @date: 2024/11/08
!-->

<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column prop="ranking" label="排名" width="80" align="center" />
      <el-table-column prop="name" label="成员" show-overflow-tooltip/>
      <el-table-column label="已使用($)">
        <template #default="scope">
          <el-tooltip :content="`${scope.row.used_quota} / ${scope.row.total_quota}`">
            <span :class="getColorClass(scope.row.used_quota, scope.row.total_quota)">
              {{ truncateToOneDecimal(scope.row.used_quota) }} /
              {{ scope.row.total_quota === -1 ? '无限制' : truncateToOneDecimal(scope.row.total_quota) }}
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
import {getAccountQuotaRankingData} from "@/api/gaia/dashboard";
import {getColorClass, truncateToOneDecimal} from "@/view/gaia/dashboard/components/index";

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
  const table = await getAccountQuotaRankingData({page: page.value, pageSize: pageSize.value})
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
