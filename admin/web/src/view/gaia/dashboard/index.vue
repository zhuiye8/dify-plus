<template>
  <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-7 py-2 gap-4 md:gap-2 gva-container2">
    <gva-card title="成员使用分析" custom-class="col-span-1 md:col-span-2 row-span-2">
      <gaia-account-money-table />
    </gva-card>

    <gva-card title="应用使用分析（以下数据定时 10分钟 更新一次，缓存了前3页数据）" custom-class="col-span-1 md:col-span-5 row-span-2">
      <gaia-app-quota-table />
    </gva-card>

    <gva-card title="密钥使用分析" custom-class="col-span-1 md:col-span-7 row-span-2">
      <gaia-app-token-quota-table />
    </gva-card>

    <gva-card title="每日密钥额度花费（单位：$）" custom-class="col-span-1 md:col-span-5 lg:col-span-7 row-span-2">
      <gva-chart :type="4" />
    </gva-card>
  </div>
</template>

<script setup>
import { GaiaAppQuotaTable, GaiaAccountMoneyTable, GvaChart, GvaCard, GaiaAppTokenQuotaTable } from "./components"
defineOptions({
  name: 'GaiaDashboard'
})
</script>

<style lang="scss" scoped>
</style>
