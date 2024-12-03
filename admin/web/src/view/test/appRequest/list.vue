<template>

  <div>
    <div class="gva-search-box">
      <el-form
        ref="searchForm"
        :inline="true"
        :model="searchInfo"
      >
        <el-form-item label="APP">
          <el-select
            v-model="searchInfo.apps"
            style="width: 500px"
            placeholder="请选择"
            filterable
            clearable
            multiple
          >
            <div
              v-for="(item, index) in appList"
              :key="index"
            >
              <el-option
                :value="item.value"
                :label="item.label"
              ></el-option>
            </div>
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select
            v-model="searchInfo.status"
            placeholder="请选择"
            clearable
          >
            <el-option label="全部" :value=0></el-option>
            <el-option label="成功" :value=1></el-option>
            <el-option label="失败" :value=2></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSubmit"
          >
            查询
          </el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <el-table
        :data="tableData"
        row-key="ID"
      >
        <el-table-column
          align="left"
          label="应用名"
          min-width="80"
          prop="name"
        />
        <el-table-column
          align="left"
          label="状态"
          min-width="50"
          prop="status"
        >
          <template #default="scope">
            <span v-if="scope.row.status" style="color: #42b983">成功</span>
            <span v-else style="color: #951d1d">失败</span>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="输入"
          min-width="280"
        >
          <template #default="scope">
            <div
              class="span-log"
              :title="scope.row.inputs"
              @click="copyText(scope.row.inputs)"
            >
              {{ scope.row.inputs }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="正确输出"
          min-width="380"
          prop="outputs"
        >
          <template #default="scope">
            <div
              class="span-log"
              :title="scope.row.outputs"
              @click="copyText(scope.row.outputs)"
            >
              {{ scope.row.outputs }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="输出"
          min-width="380"
          prop="comparison"
        >
          <template #default="scope">
            <div
              class="span-log"
              :title="scope.row.comparison"
              @click="copyText(scope.row.comparison)"
            >
              {{ scope.row.comparison }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="耗时"
          width="80"
          prop="elapsed_time"
        >
          <template #default="scope">
            <el-tooltip
              class="box-item"
              effect="dark"
              content="上次耗时"
              placement="top"
            >
              <span style="color: #8c939d">{{ scope.row.log_time }}</span>
            </el-tooltip>
            <div />
            <el-tooltip
              class="box-item"
              effect="dark"
              content="本次耗时"
              placement="top"
            >
              <span>{{ scope.row.elapsed_time }}</span>
            </el-tooltip>
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
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { gaiaAppRequestTesList } from "@/api/gaia/test";
import { useClipboard } from "@vueuse/core";

defineOptions({
  name: 'AppRequestTestList',
})


const route = useRoute()
const lock = ref(false)
const searchInfo = ref({
  status: 0,
  apps: [],
})

const onSubmit = () => {
  page.value = 1
  getTableData()
}

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
const appList = ref([])
const pageSize = ref(10)
const tableData = ref([])
const { copy, isSupported } = useClipboard()
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
const getTableData = async() => {
  const table = await gaiaAppRequestTesList({
    batch_id: route.query.id,
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo.value
  })
  if (table.code === 0) {
    lock.value = table.data.lock
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    appList.value = table.data.apps
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

const copyText = async (text) => {
  if (isSupported) {
    await copy(text)
    ElMessage({
      type: 'success',
      message: "复制成功",
    })
    return;
  }
  ElMessage({
    type: 'error',
    message: "请手动复制",
  })
}

</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
 }
  .span-log {
    overflow-wrap: break-word;
    overflow: hidden;
    word-break: break-all;
    white-space: normal;
    max-height: 300px;
    cursor:pointer;
  }
</style>
