# Omega Home 首页 UI 重新设计

## 目标
将当前的 "紫色糖果风" 改为 **暗色科技感** 风格，参考 Linear / Raycast / Vercel Dashboard。

## 设计规范

### 色彩方案
- **背景**: 纯深色 `#0a0a0a` 到 `#111` 渐变，不要花哨的彩色渐变
- **卡片背景**: `rgba(255,255,255,0.03)` ~ `0.05`，极淡的亮色
- **卡片边框**: `rgba(255,255,255,0.06)` ~ `0.08`，微妙的边界
- **hover 卡片边框**: `rgba(255,255,255,0.15)`，渐亮
- **主要文字**: `rgba(255,255,255,0.9)`
- **次要文字**: `rgba(255,255,255,0.5)`
- **强调色**: 不用渐变，用一个干净的蓝色 `#3b82f6` 做状态指示
- **在线指示灯**: `#22c55e` 绿色，微弱 glow
- **离线指示灯**: `#ef4444` 红色，无动画

### 卡片样式
- 圆角: `rounded-xl` (12px)
- 内边距: `p-5`
- hover 效果: 边框变亮 + 微弱的背景提升，**不要 translateY 上浮**（太幼稚）
- 过渡: `transition-all duration-200`
- 加一层微妙的 `shadow-[0_0_0_1px_rgba(255,255,255,0.06)]`

### 搜索栏
- 深色输入框，`bg-white/5 border border-white/10 rounded-xl`
- 左侧搜索图标用 SVG 而非 emoji
- placeholder: `text-white/30`
- focus: `border-white/20 ring-1 ring-white/10`

### 头部
- 左侧: 站点标题（白色粗体，不要 emoji logo）
- 右侧: 时钟（等宽字体，`text-white/40`）+ 暗色模式切换（图标按钮）
- 去掉 "问候语" 或改为极简一行（字小、color淡）

### 分组标题
- 小写字母风格，`text-xs uppercase tracking-wider text-white/40 font-medium`
- 不要 emoji 图标前缀（或可选）

### 服务卡片内部
- 左侧: 图标（如果是 emoji 放在 40x40 的圆角深色方块里；如果是图片直接显示）
- 右侧: 
  - 第一行: 服务名 + 状态圆点
  - 第二行: 描述（`text-sm text-white/40`）
- 不要显示 URL

### 书签栏
- 小药丸标签，`bg-white/5 border border-white/8 rounded-full px-3 py-1 text-xs`
- hover: `bg-white/10`

### Footer
- 极简: `text-white/20 text-xs`，右侧一个 "管理" 文字链接

### 动画
- 保留 fadeInUp 进场动画，但**幅度更小**（10px 而非 20px）
- 删除 breathe 呼吸灯，改为简单的稳定 glow: `shadow-[0_0_6px_rgba(34,197,94,0.4)]`

### 背景细节
- 可选: 在页面顶部加一个极淡的光晕效果（radial-gradient 从顶部中央扩散，极淡的蓝紫色，营造深度感）
- 类似: `radial-gradient(ellipse 80% 50% at 50% -20%, rgba(120,119,198,0.15), transparent)`

### 主题切换
- 保留多主题功能但重新定义主题：
  - 'midnight': 纯黑 #0a0a0a（默认）
  - 'deep-blue': 深蓝 from-[#0a1628] to-[#0d1117]
  - 'charcoal': 炭灰 from-[#18181b] to-[#09090b]
- 删除原来的彩色主题（太丑了）

## 技术约束
- 继续使用 Tailwind CDN + Alpine.js
- 保留所有现有的数据绑定和 API 调用逻辑不变
- 只改 HTML 模板和样式，不改 Go 后端
- 文件: `templates/index.html`

## 不要改的东西
- Alpine.js 数据绑定逻辑（portal() 函数的数据获取、过滤、时钟等）
- API 调用路径
- admin.html（这次不改）
- 任何 Go 代码

