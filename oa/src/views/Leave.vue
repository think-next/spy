<template>
  <div class="leave">
    <h2 class="page-title">请假审批</h2>

    <!-- 申请请假 -->
    <el-card class="form-card">
      <template #header>
        <span class="card-header-text">申请请假</span>
      </template>
      <el-form :model="form" label-width="100px" inline>
        <el-form-item label="姓名">
          <el-input v-model="form.applicant" placeholder="请输入姓名" style="width: 140px;" />
        </el-form-item>
        <el-form-item label="部门">
          <el-select v-model="form.department" placeholder="选择部门" style="width: 140px;">
            <el-option label="研发部" value="研发部" />
            <el-option label="市场部" value="市场部" />
            <el-option label="财务部" value="财务部" />
            <el-option label="人力资源部" value="人力资源部" />
            <el-option label="行政部" value="行政部" />
          </el-select>
        </el-form-item>
        <el-form-item label="请假类型">
          <el-select v-model="form.type" placeholder="选择类型" style="width: 120px;">
            <el-option label="年假" value="年假" />
            <el-option label="病假" value="病假" />
            <el-option label="事假" value="事假" />
            <el-option label="调休" value="调休" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker
            v-model="form.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始"
            end-placeholder="结束"
            value-format="YYYY-MM-DD"
            style="width: 260px;"
          />
        </el-form-item>
        <el-form-item label="事由">
          <el-input v-model="form.reason" placeholder="请假原因" style="width: 200px;" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submitLeave">提交申请</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 审批列表 -->
    <el-card class="list-card" style="margin-top: 20px;">
      <template #header>
        <span class="card-header-text">审批列表</span>
      </template>
      <el-table :data="leaves" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="applicant" label="申请人" width="90" />
        <el-table-column prop="department" label="部门" width="100" />
        <el-table-column prop="type" label="类型" width="80" />
        <el-table-column prop="startDate" label="开始" width="120" />
        <el-table-column prop="endDate" label="结束" width="120" />
        <el-table-column prop="days" label="天数" width="70" />
        <el-table-column prop="reason" label="事由" show-overflow-tooltip />
        <el-table-column label="状态" width="90">
          <template #default="scope">
            <el-tag v-if="scope.row.status === 'pending'" type="warning" size="small" class="status-tag">待审批</el-tag>
            <el-tag v-else-if="scope.row.status === 'approved'" type="success" size="small" class="status-tag">已通过</el-tag>
            <el-tag v-else type="danger" size="small" class="status-tag">已驳回</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="scope">
            <template v-if="scope.row.status === 'pending'">
              <el-button size="small" type="success" class="action-btn" @click="approve(scope.row)">通过</el-button>
              <el-button size="small" type="danger" class="action-btn" @click="reject(scope.row)">驳回</el-button>
            </template>
            <span v-else style="color: #606060; font-size: 13px;">已处理</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { leaveRequests as initData } from '../mock/data.js'
import { ElMessage } from 'element-plus'

const leaves = ref([...initData])

const form = reactive({
  applicant: '',
  department: '',
  type: '',
  dateRange: [],
  reason: '',
})

function submitLeave() {
  if (!form.applicant || !form.department || !form.type || !form.dateRange || form.dateRange.length < 2 || !form.reason) {
    ElMessage.warning('请填写完整的请假信息')
    return
  }
  const [startDate, endDate] = form.dateRange
  const days = Math.ceil((new Date(endDate) - new Date(startDate)) / 86400000) + 1
  leaves.value.unshift({
    id: leaves.value.length + 1,
    applicant: form.applicant,
    department: form.department,
    type: form.type,
    startDate,
    endDate,
    days,
    reason: form.reason,
    status: 'pending',
  })
  form.applicant = ''
  form.department = ''
  form.type = ''
  form.dateRange = []
  form.reason = ''
  ElMessage.success('请假申请已提交')
}

function approve(row) {
  row.status = 'approved'
  ElMessage.success(`已通过 ${row.applicant} 的请假申请`)
}

function reject(row) {
  row.status = 'rejected'
  ElMessage.warning(`已驳回 ${row.applicant} 的请假申请`)
}
</script>

<style scoped>
.page-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 24px 0;
  color: #fff;
  position: relative;
  display: inline-block;
  padding-bottom: 8px;
}
.page-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 60px;
  height: 3px;
  border-radius: 2px;
  background: linear-gradient(90deg, oklch(0.75 0.15 250), oklch(0.65 0.18 270));
}

.form-card,
.list-card {
  background: oklch(0.18 0.005 260 / 0.6);
  border: 1px solid oklch(0.25 0.005 260 / 0.5);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  animation: fadeInUp 0.5s ease-out backwards;
}
.form-card {
  animation-delay: 0ms;
}
.list-card {
  animation-delay: 100ms;
}

.form-card :deep(.el-card__header),
.list-card :deep(.el-card__header) {
  border-bottom: 1px solid oklch(0.25 0.005 260);
}

.card-header-text {
  color: #e0e0e0;
  font-weight: 600;
}

/* 状态标签精致样式 */
.status-tag {
  border-radius: 10px;
  padding: 4px 14px;
  box-shadow: 0 2px 6px oklch(0 0 0 / 0.15);
}

/* 操作按钮 hover 颜色过渡 */
.action-btn {
  transition: all 0.2s ease !important;
}
</style>
