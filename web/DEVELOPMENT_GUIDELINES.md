## Vite 前端项目开发规范

### 一、项目基础

- **技术栈**
  - **构建工具**: Vite（Vue 官方推荐）
  - **框架**: Vue 3（Composition API + `<script setup>`）
  - **语言**: TypeScript（严格类型）
  - **路由**: Vue Router 4
  - **状态管理**: Pinia
  - **HTTP 客户端**: Axios（统一封装）

---

### 二、目录结构规范

- **推荐结构**

```
src/
  api/           # 接口封装（Axios 实例、各模块 API）
  assets/        # 静态资源（图片、字体、公共图标）
  components/    # 通用基础组件（与业务弱相关）
  composables/   # 组合式函数（useXxx）
  constants/     # 常量配置（枚举、路由常量、业务常量）
  router/        # 路由配置
  stores/        # Pinia Store（按业务模块划分）
  styles/        # 全局样式、变量、主题
  types/         # 全局 TS 类型定义
  utils/         # 工具函数（与 UI 无关）
  views/         # 页面级组件（按路由划分）
  App.vue
  main.ts
```

- **命名约定**
  - 目录与文件：统一使用小写 + 中划线或小驼峰，如 `user-profile` / `userProfile.ts`。
  - 组件文件：`PascalCase`，如 `UserInfoForm.vue`。
  - 组合式函数：`useXxx.ts`，如 `useUserList.ts`。
  - Store：`useXxxStore.ts`，如 `useAuthStore.ts`。
  - 类型文件：集中在 `types/` 下按模块拆分（如 `auth.ts`、`user.ts`）。

---

### 三、代码静态检查

- **语法与质量检查**
  - 使用 **ESLint + @typescript-eslint**，Vue 项目使用 `eslint-plugin-vue`。
  - 根目录定义 `.eslintrc` / `eslint.config`，包含：
    - TypeScript 相关规则。
    - Vue3 SFC 相关规则。
    - 基础代码质量规则（如禁止未使用变量、禁止多余 `console.log` 等）。
  - 所有 TS 文件、Vue SFC 必须通过 ESLint 检查后方可合并。

- **统一脚本约定**
  - `npm run lint`：执行 ESLint 检查。
  - `npm run format`：执行 Prettier 格式化（或通过 lint-staged 对暂存区格式化）。

---

### 四、TypeScript 使用规范

- **类型配置**
  - `tsconfig.json` 中推荐开启：
    - `strict: true`
    - `noImplicitAny: true`
    - `strictNullChecks: true`
  - 禁止随意使用 `any`；如确需使用必须加注释解释原因。

- **接口与数据模型**
  - 后端接口的请求/响应类型必须有明确的 TS 类型定义。
  - 公共类型放在 `src/types` 中，按业务模块拆分：
    - `types/auth.ts`
    - `types/user.ts`
    - `types/common.ts`
  - 避免在组件中大量使用 `as XxxType` 强制断言掩盖类型问题。

---

### 五、Vue 组件开发规范

- **SFC 基本要求**
  - 推荐统一使用 `<script setup lang="ts">`。
  - 组件文件内部顺序固定：
    1. `<template>`
    2. `<script setup lang="ts">`
    3. `<style>` / `<style scoped>`
  - 每个组件关注单一职责，避免超大组件；超过约 300 行须考虑拆分为子组件。

- **组件分类**
  - **基础组件**（`components/base`）：
    - 与业务无关，如 `BaseButton`、`BaseModal`、`BaseTable`。
  - **业务组件**（`components/module-name`）：
    - 绑定具体业务，如 `components/profile/UserInfoForm.vue`。
  - **页面组件**（`views` 下）：
    - 与路由一一对应，不直接在其他页面中复用。

- **Props 与事件**
  - 所有 Props 必须声明类型和必要的默认值。
  - 事件命名：
    - 代码中使用 `camelCase`（如 `onSubmit`）。
    - 模板中使用 `kebab-case`（如 `@submit`）。
  - 避免在父组件中直接操作子组件内部状态，使用 Props + Emits 传递。

- **组合式函数（Composables）**
  - 放在 `src/composables`，命名为 `useXxx.ts`。
  - 只负责逻辑复用，不直接操作 DOM，不耦合具体 UI 样式。
  - 开放明确的入参与返回值，类型完整。

---

### 六、路由与导航规范

- **路由结构**
  - 所有路由定义集中在 `src/router`，通常为 `router/index.ts`。
  - 大型项目推荐拆分模块路由：
    - `router/modules/auth.ts`
    - `router/modules/user.ts`
    - `router/modules/oauth.ts`
  - 在 `index.ts` 中统一汇总。

- **命名规范**
  - `name` 使用 `PascalCase`，如 `UserList`、`UserDetail`。
  - `path` 使用 `kebab-case`，如：
    - `/users`
    - `/user-detail/:id`

- **导航守卫与权限**
  - 在全局前置守卫中统一处理：
    - 登录态校验（token 是否存在/有效）。
    - 基础权限判断（如是否有访问该路由的权限）。
  - 复杂权限逻辑抽离到独立模块（如 `utils/permission.ts` 或 `composables/usePermission.ts`）。
  - 路由元信息（`meta`）用于标记权限、标题、缓存策略等。

---

### 七、状态管理（Pinia）

- **Store 结构与命名**
  - 所有 Store 放在 `src/stores`。
  - 文件命名：`useXxxStore.ts`，例如：
    - `useAuthStore.ts`
    - `useUserStore.ts`
  - Store 名称与文件名保持一致，便于查找。

- **使用规范**
  - Store 负责跨组件共享状态和复杂业务逻辑。
  - 避免将纯页面局部状态（如某个对话框是否打开）全部放入 Store。
  - 避免在 Store 外部直接修改 Store 内部的响应式对象结构。

---

### 八、网络请求与接口封装

- **Axios 实例封装**
  - 在 `src/api/request.ts` 中统一创建 Axios 实例，配置：
    - `baseURL`
    - `timeout`
    - 请求拦截器（如自动携带 token）
    - 响应拦截器（统一错误处理、数据结构解包）

- **模块接口文件**
  - 按业务模块拆分：
    - `api/auth.ts`
    - `api/user.ts`
    - `api/oauth_client.ts`
  - 每个接口函数：
    - 明确入参与返回类型。
    - 只做 HTTP 调用与类型转换，不包含页面逻辑。

- **错误与异常处理**
  - 在响应拦截器中统一处理通用错误（如 401、403、500）。
  - 在业务层根据返回码进行细粒度处理（如表单错误提示）。

---

### 九、错误展示与用户提示

- **全局错误页**
  - 提供统一的错误页面（如 `ErrorPage.vue`），用于展示：
    - 404（页面不存在）
    - 403（无权限）
    - 5xx（服务器错误）
  - 通过路由或全局异常处理统一跳转。

- **局部错误提示**
  - 使用统一的通知组件（如 `Message`、`Notification`）展示接口失败等信息。
  - 避免在最终代码中保留原生 `alert` 和大量 `console.error`。

---

### 十、样式与设计规范

- **样式组织**
  - 在 `src/styles` 目录中维护：
    - `reset.(css|scss)`：浏览器重置样式。
    - `variables.(css|scss)`：颜色、字体、间距等变量。
    - `global.(css|scss)`：全局通用样式。
  - 组件私有样式放在组件内 `<style scoped>` 中。

- **命名与规范**
  - **统一采用「近似 BEM」的命名风格**：
    - 格式：`block__element--modifier`
    - **Block（块）**：独立的功能单元，页面级根容器建议使用 `xxx-page` 或 `page-xxx` 作为 block 名。
      - 示例：`error-page`、`users-page`、`login-page`
    - **Element（元素）**：block 内部的组成部分，使用 `__` 连接。
      - 示例：`error-page__card`、`error-page__content`、`users-page__toolbar`
    - **Modifier（修饰符）**：表示状态或外观变化，使用 `--` 连接。
      - 示例：`btn--primary`、`error-page__action--secondary`、`card--disabled`
  - **命名原则**：
    - 避免无语义的类名（如 `.red`、`.big`、`.container`），使用语义化命名。
    - 类名应清晰表达其职责与层级关系。
    - 组件内私有样式使用 `<style scoped>`，类名可适当简化但仍需保持语义。
  - **示例对比**：
    - ❌ 不推荐：`.container` / `.wrapper` / `.box`
    - ✅ 推荐：`.users-page` / `.users-page__header` / `.users-page__table`

- **配色与主题**
  - 使用现代、简洁的设计风格，颜色通过变量统一管理。
  - 严禁使用高饱和、刺眼的颜色作为主色，尤其避免让界面显得浮夸。
  - 主题切换（如暗色模式）通过 CSS 变量或类名切换实现。

- **CSS 变量使用规范**
  - **强制要求**：所有颜色、圆角、阴影、间距等设计属性必须优先使用 `src/styles/variables.css` 中定义的 CSS 变量。
  - **禁止行为**：在组件样式中直接硬编码 hex 颜色值（如 `#303133`）、固定像素值的圆角/阴影等。
  - **变量分类**：
    - 颜色变量：`--color-primary`、`--color-text-primary`、`--color-background` 等。
    - 圆角变量：`--border-radius-small`、`--border-radius-button`、`--border-radius-card` 等。
    - 阴影变量：`--shadow-card`、`--shadow-layer` 等。
    - 间距变量：`--spacing-sm`、`--spacing-md`、`--spacing-lg` 等。
    - 字号变量：`--font-size-base`、`--font-size-title` 等。
  - **使用示例**：
    ```css
    /* ❌ 不推荐：硬编码颜色与圆角 */
    .my-card {
      background: #ffffff;
      color: #303133;
      border-radius: 12px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }

    /* ✅ 推荐：使用 CSS 变量 */
    .my-card {
      background: var(--color-card-background);
      color: var(--color-text-primary);
      border-radius: var(--border-radius-card);
      box-shadow: var(--shadow-card);
    }
    ```
  - **扩展变量**：当需要新的设计属性时，应先在 `variables.css` 中定义变量，再在组件中使用。

---

### 十一、环境与配置管理

- **环境变量**
  - 使用 `.env`、`.env.development`、`.env.production` 等区分环境。
  - 所有需要在前端使用的环境变量必须以 `VITE_` 前缀开头。
  - 禁止在代码中直接写死敏感信息（如后端地址、密钥）。

- **Vite 配置**
  - 使用别名简化路径：
    - `@` 指向 `src` 目录。
  - 开发服务器配置仅包含开发相关项（如 `port`、`open`、`proxy`），生产差异通过环境变量控制。

---
