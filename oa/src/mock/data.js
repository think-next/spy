// 公告数据
export const announcements = [
  { id: 1, title: '关于2026年端午节放假安排的通知', author: '行政部', date: '2026-06-01', content: '端午节6月6日至6月8日放假3天，6月9日正常上班。' },
  { id: 2, title: '办公室搬迁通知', author: '行政部', date: '2026-05-28', content: '因公司业务扩展，研发部将于6月15日搬迁至B栋12层。' },
  { id: 3, title: '6月份团建活动报名', author: '人力资源部', date: '2026-05-25', content: '6月20日（周六）组织户外拓展训练，请于6月10日前报名。' },
  { id: 4, title: '网络安全培训通知', author: '技术部', date: '2026-05-20', content: '全员网络安全意识培训将于6月5日下午2点在会议室A举行。' },
]

// 请假单数据
export const leaveRequests = [
  { id: 1, applicant: '张三', department: '研发部', type: '年假', startDate: '2026-06-10', endDate: '2026-06-12', days: 3, reason: '回家探亲', status: 'pending' },
  { id: 2, applicant: '李四', department: '市场部', type: '病假', startDate: '2026-06-03', endDate: '2026-06-04', days: 2, reason: '身体不适需休息', status: 'pending' },
  { id: 3, applicant: '王五', department: '财务部', type: '事假', startDate: '2026-06-05', endDate: '2026-06-05', days: 1, reason: '处理个人事务', status: 'approved' },
  { id: 4, applicant: '赵六', department: '研发部', type: '年假', startDate: '2026-05-28', endDate: '2026-05-30', days: 3, reason: '旅游出行', status: 'rejected' },
  { id: 5, applicant: '孙七', department: '人力资源部', type: '调休', startDate: '2026-06-06', endDate: '2026-06-06', days: 1, reason: '补上周六加班', status: 'pending' },
]

// 员工数据
export const employees = [
  { id: 1, name: '张三', department: '研发部', position: '高级前端工程师', phone: '138****1234', email: 'zhangsan@company.com' },
  { id: 2, name: '李四', department: '市场部', position: '市场经理', phone: '139****5678', email: 'lisi@company.com' },
  { id: 3, name: '王五', department: '财务部', position: '财务主管', phone: '137****9012', email: 'wangwu@company.com' },
  { id: 4, name: '赵六', department: '研发部', position: '后端工程师', phone: '136****3456', email: 'zhaoliu@company.com' },
  { id: 5, name: '孙七', department: '人力资源部', position: 'HR专员', phone: '135****7890', email: 'sunqi@company.com' },
  { id: 6, name: '周八', department: '研发部', position: '技术总监', phone: '134****2345', email: 'zhouba@company.com' },
  { id: 7, name: '吴九', department: '市场部', position: '市场专员', phone: '133****6789', email: 'wujiu@company.com' },
  { id: 8, name: '郑十', department: '行政部', position: '行政助理', phone: '132****0123', email: 'zhengshi@company.com' },
  { id: 9, name: '陈一一', department: '财务部', position: '会计', phone: '131****4567', email: 'chenyy@company.com' },
  { id: 10, name: '林一二', department: '研发部', position: '测试工程师', phone: '130****8901', email: 'linye@company.com' },
]
