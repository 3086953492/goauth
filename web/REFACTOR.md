## Vite 项目规范整改计划（基于当前 web 子项目）

> 目标：让当前 `web` 前端项目逐步对齐《Vite 前端项目开发规范》。  
> 说明：按步骤执行，每个大步骤内部也可以拆成多次提交，避免一次性大改。

---

### 步骤 1：工程配置与静态检查补齐

#### 1.1 增加 ESLint 配置

- 现状：
  - `web` 目录下没有 `.eslintrc*` / `eslint.config.*` 文件。
  - `package.json` 中也没有 `lint` 脚本。
- 建议整改：
  1. 在 `web` 根目录新增 ESLint 配置文件（如 `eslint.config.mjs` 或 `.eslintrc.cjs`），启用：
     - `@typescript-eslint` + `eslint-plugin-vue` + 基础规则。
  2. 将 `src/**/*.ts`、`src/**/*.vue` 全部纳入 ESLint 检查范围。
  3. 在 `package.json` 中增加：
     - `lint` 脚本：运行 ESLint 检查全部源码。

---

### 步骤 2：样式与设计体系规范化

#### 2.1 建立 `src/styles` 目录并拆分全局样式

- 现状：
  - 当前全局样式集中在 `src/style.css` 中，包含：
    - 重置样式、滚动条样式、Element Plus 覆盖等。
  - 规范要求：使用 `src/styles` 目录拆分 `reset / variables / global`。
- 建议整改：
  1. 在 `src` 下新建 `styles` 目录，至少包含：
     - `reset.css`：基础重置（从 `style.css` 中抽取与重置相关的部分）。
     - `global.css`：全局布局、滚动条样式、Element Plus 通用覆盖等。
     - （可选）`variables.(css|scss)`：颜色、间距、字体变量。
  2. 将 `style.css` 中的内容按职责拆分迁移到上述文件中。
  3. 在 `main.ts` 中修改全局样式引入方式：
     - 从 `import './style.css'` 改为引入新的 `styles` 文件（例如 `import './styles/reset.css'`、`import './styles/global.css'` 等）。
  4. 保持组件内样式使用 `<style scoped>`，避免把组件私有样式放回全局。

#### 2.2 样式命名与设计约定固化

- 现状：
  - 各组件中类名整体还算语义化，例如：
    - `ErrorPage.vue` 中的 `.error-container` / `.error-card` 等。
    - `Users.vue` 中的 `.users-wrapper` / `.users-container` 等。
  - 但目前项目层面没有统一"命名规则文档"和"设计变量"落地。
- 建议整改：
  1. 在 `styles/global` 或单独文档中写清：
     - 类命名采用 BEM 或统一的 block/element 风格。
  2. 把常用的颜色、阴影、圆角统一抽象为 CSS 变量，放在 `variables` 文件中，以便未来统一换肤或调整风格。
- **命名约定**（已在 `DEVELOPMENT_GUIDELINES.md` 第十节明确）：
  - 统一采用「近似 BEM」命名：`block__element--modifier`。
  - 页面级根容器使用 `xxx-page` 格式（如 `error-page`、`users-page`）。
  - 内层元素使用 `__` 连接（如 `error-page__card`、`users-page__toolbar`）。
  - 状态/外观修饰符使用 `--` 连接（如 `btn--primary`、`card--disabled`）。
- **完成标准**：
  - ✅ 至少一个页面（如 `ErrorPage.vue`）的类名已改造为约定风格。
  - ✅ 关键颜色、圆角、阴影、间距已抽取到 `src/styles/variables.css` 并在全局/页面样式中使用。
  - ✅ `DEVELOPMENT_GUIDELINES.md` 中已补充命名规范与 CSS 变量使用示例。

---

### 步骤 3：Store 文件命名与结构对齐规范

#### 3.1 规范 `Pinia` Store 文件命名

- 现状：
  - Store 文件：`src/stores/auth.ts`。
  - Store 定义：`export const useAuthStore = defineStore('auth', () => { ... })`，逻辑设计良好。
  - 规范要求：Store 文件名采用 `useXxxStore.ts`。
- 建议整改：
  1. 将 `src/stores/auth.ts` 重命名为 `src/stores/useAuthStore.ts`。
  2. 全局替换导入路径：
     - 例如：
       - `main.ts` 中的 `import { useAuthStore } from './stores/auth'`。
       - `Navbar.vue`、`Home.vue`、`Users.vue`、`useAuth.ts` 等处的 `@/stores/auth`。
  3. 确保所有引用都改为 `@/stores/useAuthStore` 或相对路径的新文件名。

---

### 步骤 4：Composable 职责边界与 UI 解耦

#### 4.1 `useAuth.ts` 中的 UI 与逻辑解耦

- 现状（`src/composables/useAuth.ts`）：
  - Composable 内部直接：
    - 使用 `ElMessage` 做 UI 提示（成功/警告）。
    - 使用 `router.push` 做页面跳转。
  - 规范要求：组合式函数“只负责逻辑复用，不直接操作 DOM / 不耦合具体 UI 细节”。
- 建议整改：
  1. 拆分纯逻辑与 UI/导航：
     - 将 `handleLogin`、`handleRegister` 拆成：
       - 一个只负责调用 API、更新 Store、返回结果状态的逻辑函数；
       - UI 层（组件）根据返回的状态决定是否弹出 `ElMessage`、是否跳转路由。
  2. 例如：
     - `useAuth` 返回的函数只返回：
       - `success` 布尔值。
       - 可选的错误消息 / 数据。
     - `Login.vue` / `Register.vue` 页面中根据结果做：
       - `ElMessage.success / warning`。
       - `router.push` 或其他跳转逻辑。
  3. 保留 composable 中使用 `FormInstance` 做表单校验是可以接受的，但可以考虑也下放到组件，由组件控制校验节奏。

---

### 步骤 5：Axios 请求封装与错误处理规范化

#### 5.1 成功 / 失败判断逻辑优化（`src/api/request.ts`）✅ 已完成

- 现状（关键逻辑）：
  - 响应拦截器中：
    - `if (res.success === true || response.status === 200) { return res }`
    - 否则才会执行业务错误分支，`ElMessage.error` 并 `Promise.reject`。
  - 问题：
    - 当 `HTTP 200` 但 `res.success !== true` 时，仍被当作成功返回，不符合"统一业务错误处理"的期望。
- 建议整改：
  1. 重新定义"成功"的判断条件，例如：
     - **严格以业务字段 `success === true` 为成功**；
     - 或者统一约定：所有正常业务返回都保证 `success === true`，否则一律视为错误。
  2. 将 `if` 条件改为只在业务成功时返回，其他情况走错误提示分支。
  3. 保留统一的错误消息提示，并保证业务层不需要重复判断 `success`。
- **已完成的改动**：
  1. 修改了 `src/api/request.ts` 响应拦截器，严格以 `res.success === true` 为成功标准。
  2. 移除了所有业务代码中对 `response.success` 的冗余判断，改为依赖 try/catch 处理失败。
  3. 涉及文件：
     - `src/api/request.ts`：核心拦截器逻辑修改
     - `src/composables/useUserList.ts`：移除 success 判断
     - `src/composables/useOAuthClientList.ts`：移除 success 判断
     - `src/composables/useOAuthClientForm.ts`：移除 success 判断
     - `src/views/OAuthClients.vue`：移除多处 success 判断

#### 5.2 错误路由与全局错误页联动 ✅

- 现状：
  - 已有 `ErrorPage.vue`，主要用于展示 OAuth 授权相关的错误（读取 `error`、`error_description` query）。
  - 路由中有 `/error` 路由配置。
  - 但当前 Axios 拦截器中对 4xx/5xx 错误只做 `ElMessage.error` 提示，没有与 `/error` 页联动。
- 已完成整改：
  1. **扩展 ErrorPage.vue**：
     - 增加 HTTP 状态码错误映射（403、404、500、502、503）
     - 增加通用错误类型（forbidden、not_found、internal_error、network_error、timeout）
     - 为不同错误类型提供友好的默认描述
     - 保持对 OAuth 错误的完全兼容
  2. **改造 Axios 拦截器**：
     - 403 权限错误：跳转到错误页，携带错误码和描述
     - 500/502/503 服务器错误：跳转到错误页，携带错误码和描述
     - 404 接口错误：仅弹出 ElMessage 提示（不跳转，因为通常是接口不存在）
     - 401 认证错误：保持原有逻辑（刷新 token 或跳转登录）
     - 其他错误：弹出 ElMessage 提示
     - 增加防抖机制，避免多次跳转错误页
  3. **使用场景明确**：
     - 普通接口错误（400、404 等）：ElMessage 提示
     - 严重错误（403、5xx）：跳转错误页面
     - OAuth 错误：保持原有处理方式

---

### 步骤 6：路由结构与模块化（可选增强）

#### 6.1 路由模块拆分（基于 `src/router/index.ts`）

- 现状：
  - 所有路由均集中在 `router/index.ts`：
    - 包括登录、注册、主页、个人信息、用户列表、OAuth 客户端、OAuth 授权、错误页。
  - 对于当前体量，这样没有功能性问题，但与「按模块拆分路由配置」的规范存在差距。
- 建议整改（可在项目复杂度上升后再做）：
  1. 在 `src/router` 下新建 `modules` 目录：
     - `auth.ts`：登录、注册相关路由。
     - `user.ts`：用户列表、个人信息路由。
     - `oauth.ts`：OAuth 客户端、授权路由。
  2. 在每个模块文件中导出 `RouteRecordRaw[]`，在 `index.ts` 中统一合并。
  3. 保持路由 `name` 使用 `PascalCase`（当前 `Login` / `Register` / `Home` 等已符合）。
