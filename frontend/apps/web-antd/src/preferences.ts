import {
  defineOverridesPreferences,
  definePreferencesExtension,
} from '@vben/preferences';

interface WebAntdPreferencesExtension {
  defaultTableSize: number;
  enableFormFullscreen: boolean;
}

export const overridesPreferences = defineOverridesPreferences({
  app: {
    /** 后端路由模式 */
    accessMode: 'backend',
    name: import.meta.env.VITE_APP_TITLE,
    enableRefreshToken: true,
    /** 默认色弱模式 */
    colorWeakMode: false,
    /** 灰色模式 */
    grayMode: false,
    /** 主题色 */
    primaryColor: 'hsl(211, 98%, 48%)',
  },
  theme: {
    /** 默认亮色模式，支持 light / dark / auto */
    mode: 'light',
    /** 半亮色侧边栏 */
    semiDarkSidebar: false,
    /** 半亮色头部 */
    semiDarkHeader: false,
  },
  sidebar: {
    /** 侧边栏默认展开 */
    collapsed: false,
  },
  footer: {
    enable: false,
    fixed: false,
  },
  copyright: {
    companyName: import.meta.env.VITE_APP_TITLE,
    companySiteLink: '',
  },
});

export const preferencesExtension =
  definePreferencesExtension<WebAntdPreferencesExtension>({
    tabLabel: 'preferences.antd.tabLabel',
    title: 'preferences.antd.title',
    fields: [
      {
        component: 'switch',
        defaultValue: true,
        key: 'enableFormFullscreen',
        label: 'preferences.antd.fields.enableFormFullscreen.label',
        tip: 'preferences.antd.fields.enableFormFullscreen.tip',
      },
      {
        component: 'number',
        componentProps: { max: 200, min: 10, step: 10 },
        defaultValue: 20,
        key: 'defaultTableSize',
        label: 'preferences.antd.fields.defaultTableSize.label',
      },
    ],
  });
