<template>
  <div class="contacts">
    <h2 class="page-title">通讯录</h2>

    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <span>员工列表</span>
          <el-input
            v-model="keyword"
            placeholder="搜索姓名/部门..."
            style="width: 220px;"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </template>
      <el-table :data="filteredEmployees" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="department" label="部门" width="120">
          <template #default="scope">
            <el-tag size="small" effect="plain">{{ scope.row.department }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="position" label="职位" width="150" />
        <el-table-column prop="phone" label="手机" width="130" />
        <el-table-column prop="email" label="邮箱" show-overflow-tooltip />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { employees } from '../mock/data.js'

const keyword = ref('')

const filteredEmployees = computed(() => {
  const kw = keyword.value.toLowerCase().trim()
  if (!kw) return employees
  return employees.filter(e =>
    e.name.toLowerCase().includes(kw) || e.department.toLowerCase().includes(kw)
  )
})
</script>

<style scoped>
.page-title {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 24px 0;
  color: #fff;
}

.main-card {
  background: oklch(0.18 0.005 260);
  border: 1px solid oklch(0.25 0.005 260);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #e0e0e0;
  font-weight: 600;
}
</style>
