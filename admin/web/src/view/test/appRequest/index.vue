<template>
  <div>
    <div class="api-request-search-box">
      <el-form-item style="position: absolute;right: 20px;">
        <el-button
          type="warning"
          icon="refresh"
          :disabled="lock"
          @click="toSyncTest"
        >
          开始批量测试API应用
        </el-button>
      </el-form-item>
    </div>
    <div class="gva-table-box">
      <el-table
        :data="tableData"
        row-key="ID"
      >
        <el-table-column
          align="left"
          label="状态"
          min-width="50"
          prop="status"
        >
          <template #default="scope">
            <span v-if="scope.row.status===1" style="color: #d68a00">执行中</span>
            <span v-else style="color: #1f9121">已完成</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="累计测试数"
          min-width="120"
          prop="sum"
        />
        <el-table-column
          align="left"
          label="累计应用数"
          min-width="120"
          prop="app"
        />
        <el-table-column
          align="left"
          label="创建时间"
          min-width="210"
        >
          <template #default="scope">
            <span>{{ toTime(scope.row.create_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="结束时间"
          min-width="210"
        >
          <template #default="scope">
            <span>{{ toTime(scope.row.end_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="成功数"
          min-width="90"
        >
          <template #default="scope">
            <span style="color: #67c23a">{{ scope.row.success_count }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="失败数"
          min-width="90"
        >
          <template #default="scope">
            <span style="color: #951d1d">{{ scope.row.failure_count }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="总耗时"
          min-width="90"
        >
          <template #default="scope">
            <span style="color: #d68a00">{{ toFormatTime(scope.row.end_time, scope.row.create_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="查看"
          min-width="90"
        >
          <template #default="scope">
            <el-link @click="toList(scope.row.id)" type="primary">
              查看详情
            </el-link>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>


import { getAuthorityList } from '@/api/authority'

import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {gaiaAppRequestBatch, gaiaAppRequestTest} from "@/api/gaia/test";
import {useRouter} from "vue-router";
import {formatTimeToStr, formatTime} from "@/utils/date";
const lock = ref(false)
const router = useRouter()
defineOptions({
  name: 'AppRequestTestBatch',
})


// 初始化相关
const setAuthorityOptions = (AuthorityData, optionsData) => {
  AuthorityData &&
        AuthorityData.forEach(item => {
          if (item.children && item.children.length) {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              children: []
            }
            setAuthorityOptions(item.children, option.children)
            optionsData.push(option)
          } else {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName
            }
            optionsData.push(option)
          }
        })
}

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}
// 时间戳转日期
const toTime = (val) => {
  if (val === 0) {
    return ""
  }
  return formatTimeToStr(val * 1000)
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}
// 自动设置时分秒
const toFormatTime = (end, start) => {
  if (end === 0)
    return ""
  return formatTime(end - start)
}
// 跳转到列表页
const toList = (id) => {
  const query = {};
  query.id = id;
  router.push({name: 'AppRequestTestList', query, params: { id } });
// :href="'/#/layout/AppRequestTestList?id=' + scope.row.id"
}
// 查询
const getTableData = async() => {
  const table = await gaiaAppRequestBatch({ page: page.value, pageSize: pageSize.value })
  if (table.code === 0) {
    lock.value = table.data.lock
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

watch(() => tableData.value, () => {
  setAuthorityIds()
})

const initPage = async() => {
  getTableData()
  const res = await getAuthorityList()
  setOptions(res.data)
}

initPage()

const setAuthorityIds = () => {
  tableData.value && tableData.value.forEach((user) => {
    user.authorityIds = user.authorities && user.authorities.map(i => {
      return i.authorityId
    })
  })
}

const authOptions = ref([])
const setOptions = (authData) => {
  authOptions.value = []
  setAuthorityOptions(authData, authOptions.value)
}

const toSyncTest = async () => {
  ElMessageBox.confirm('该功能会调用【API生产】空间所有被调用过的API应用2次，是否确认开始进行测试？')
  .then(async () => {
    const res = await gaiaAppRequestTest()
    if (res.code === 0) {
      ElMessage.success('执行成功')
      await getTableData()
    }
  })
  .catch(() => {
    this.$message({
      type: 'info',
      message: '已取消测试'
    });
  })
}

</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
 }
  .api-request-search-box {
    height: 40px;
    padding: 0;
    margin: 0;
  }
  .span-log {
    overflow-wrap: break-word;
    overflow: hidden;
    word-break: break-all;
    white-space: normal;
    max-height: 300px;
  }
</style>
