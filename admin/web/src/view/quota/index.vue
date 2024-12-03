<template>

  <div>
    <div class="gva-search-box">
      <el-form
        ref="searchForm"
        :inline="true"
        :model="searchInfo"
      >
        <el-form-item label="搜索">
          <el-input
            v-model="searchInfo.keyword"
            placeholder="用户名"
          />
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
          label="头像"
          min-width="75"
        >
          <template #default="scope">
            <CustomPic
              style="margin-top:8px"
              :pic-src="scope.row.headerImg"
            />
          </template>
        </el-table-column>
        <el-table-column
          align="left"
          label="排名"
          min-width="50"
          prop="ranking"
        />
        <el-table-column
          align="left"
          label="成员"
          min-width="150"
          prop="name"
        />
        <el-table-column
          align="left"
          label="已使用配额"
          min-width="150"
          prop="used_quota"
        />
        <el-table-column
          align="left"
          label="配额"
          min-width="180"
          prop="total_quota"
        />
        <el-table-column
          align="left"
          label="余额"
          min-width="180"
        >
          <template #default="scope">
            {{ scope.row.total_quota - scope.row.used_quota }}
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          min-width="250"
          fixed="right"
        >
          <template #default="scope">
            <el-button
              type="primary"
              link
              icon="magic-stick"
              @click="resetQuotaFunc(scope.row)"
            >修改额度</el-button>
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

import { setUserQuota, getManagementList } from '@/api/user'

import { getAuthorityList } from '@/api/authority'
import CustomPic from '@/components/customPic/index.vue'

import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

defineOptions({
  name: 'QuotaList',
})

const searchInfo = ref({
  keyword: '',
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
const getTableData = async() => {
  const table = await getManagementList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
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

const resetQuotaFunc = (row) => {
  ElMessageBox.prompt(`是否修改 ${row.name} 的额度(单位美元)`, '请注意', {
    confirmButtonText: '修改',
    cancelButtonText: '取消',
    inputPattern: /^\d+(\.\d+)?$/,
    inputErrorMessage: '请配置用户月额度限制',
    inputValue: row.total_quota,
  }).then(async ({value}) => {
    const res = await setUserQuota({
      quota: Number(value.trim()),
      uid: row.uid,
    })
    if (res.code === 0) {
      initPage()
      ElMessage({
        type: 'success',
        message: res.msg,
      })
    } else {
      ElMessage({
        type: 'error',
        message: res.msg,
      })
    }
  })
}
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

</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
 }
</style>
