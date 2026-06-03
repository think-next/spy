<template>
  <div class="announcement">
    <h2 class="page-title">公告管理</h2>

    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <span>公告列表</span>
          <el-button type="primary" @click="showDialog = true">
            <el-icon><Plus /></el-icon> 新增公告
          </el-button>
        </div>
      </template>
      <el-table :data="list" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="title" label="标题" show-overflow-tooltip />
        <el-table-column prop="author" label="发布者" width="120" />
        <el-table-column prop="date" label="发布日期" width="130" />
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button size="small" type="primary" link @click="viewDetail(scope.row)">查看</el-button>
            <el-button size="small" type="danger" link @click="removeItem(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增公告弹窗 -->
    <el-dialog v-model="showDialog" title="新增公告" width="500px" destroy-on-close>
      <el-form :model="form" label-width="80px">
        <el-form-item label="标题">
          <el-input v-model="form.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-form-item label="发布者">
          <el-input v-model="form.author" placeholder="请输入发布者" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="form.content" type="textarea" :rows="4" placeholder="请输入公告内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="addItem">确定</el-button>
      </template>
    </el-dialog>

    <!-- 查看详情弹窗 -->
    <el-dialog v-model="showDetail" :title="detailItem.title" width="500px" class="detail-dialog">
      <div class="detail-meta">
        <span>发布者：{{ detailItem.author }}</span>
        <span>日期：{{ detailItem.date }}</span>
      </div>
      <el-divider />
      <p class="detail-content">{{ detailItem.content }}</p>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { announcements as initData } from '../mock/data.js'
import { ElMessage } from 'element-plus'

const list = ref([...initData])
const showDialog = ref(false)
const showDetail = ref(false)
const detailItem = reactive({ id: 0, title: '', author: '', date: '', content: '' })

const form = reactive({ title: '', author: '', content: '' })

function addItem() {
  if (!form.title || !form.content) {
    ElMessage.warning('标题和内容不能为空')
    return
  }
  const today = new Date().toISOString().slice(0, 10)
  list.value.unshift({
    id: list.value.length + 1,
    title: form.title,
    author: form.author || '管理员',
    date: today,
    content: form.content,
  })
  form.title = ''
  form.author = ''
  form.content = ''
  showDialog.value = false
  ElMessage.success('公告已发布')
}

function removeItem(id) {
  list.value = list.value.filter(item => item.id !== id)
  ElMessage.success('已删除')
}

function viewDetail(row) {
  Object.assign(detailItem, row)
  showDetail.value = true
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

.main-card {
  background: oklch(0.18 0.005 260 / 0.6);
  border: 1px solid oklch(0.25 0.005 260 / 0.5);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
  animation: fadeInUp 0.5s ease-out backwards;
  animation-delay: 0ms;
}
.main-card :deep(.el-card__header) {
  border-bottom: 1px solid oklch(0.25 0.005 260);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #e0e0e0;
  font-weight: 600;
}

.detail-meta {
  display: flex;
  gap: 20px;
  color: #909090;
  font-size: 13px;
}

.detail-content {
  color: #c0c0c0;
  line-height: 1.8;
}
</style>
