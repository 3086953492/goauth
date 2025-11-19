import js from '@eslint/js';
import tseslint from 'typescript-eslint';
import pluginVue from 'eslint-plugin-vue';
import vueParser from 'vue-eslint-parser';

export default [
  // 忽略构建产物和依赖目录
  {
    ignores: ['dist/**', 'node_modules/**', '*.config.js', '*.config.ts'],
  },

  // JavaScript 基础推荐规则
  js.configs.recommended,

  // TypeScript 文件配置
  ...tseslint.configs.recommendedTypeChecked.map(config => ({
    ...config,
    files: ['src/**/*.ts', 'src/**/*.tsx'],
  })),
  {
    files: ['src/**/*.ts', 'src/**/*.tsx'],
    languageOptions: {
      parserOptions: {
        project: './tsconfig.app.json',
        tsconfigRootDir: import.meta.dirname,
      },
    },
  },

  // Vue 文件配置
  ...pluginVue.configs['flat/recommended'],
  {
    files: ['src/**/*.vue'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: tseslint.parser,
        project: './tsconfig.app.json',
        tsconfigRootDir: import.meta.dirname,
        extraFileExtensions: ['.vue'],
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
    },
  },

  // 通用规则配置
  {
    files: ['src/**/*.ts', 'src/**/*.tsx', 'src/**/*.vue'],
    plugins: {
      '@typescript-eslint': tseslint.plugin,
    },
    rules: {
      // 禁止未使用的变量（TypeScript 已有相关检查，这里保持一致）
      '@typescript-eslint/no-unused-vars': [
        'warn',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
        },
      ],

      // 控制 console 使用，允许 warn 和 error
      'no-console': ['warn', { allow: ['warn', 'error'] }],

      // Vue 3 Composition API 相关规则调整
      'vue/multi-word-component-names': 'off', // 允许单词组件名（如 Home, Login）
      'vue/require-default-prop': 'off', // 不强制要求所有 prop 都有默认值

      // TypeScript 相关规则调整
      '@typescript-eslint/no-explicit-any': 'warn', // any 使用警告而非错误
      '@typescript-eslint/explicit-module-boundary-types': 'off', // 不强制要求导出函数的返回类型
    },
  },
];

