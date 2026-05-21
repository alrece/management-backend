import type { AuthPermissionInfo } from '@vben/types';

import { baseRequestClient, requestClient } from '#/api/request';
import { transformPermissionInfo } from '#/api/transform';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    password?: string;
    username?: string;
  }

  /** 登录接口返回值 */
  export interface LoginResult {
    accessToken: string;
    refreshToken: string;
  }
}

/** 登录 POST /api/auth/login */
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', data);
}

/** 刷新 accessToken POST /api/auth/refresh */
export async function refreshTokenApi(refreshToken: string) {
  return baseRequestClient.post('/auth/refresh', {
    refreshToken,
  });
}

/** 退出登录 POST /api/auth/logout */
export async function logoutApi(accessToken: string) {
  return baseRequestClient.post(
    '/auth/logout',
    {},
    {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    },
  );
}

/**
 * 获取权限信息 GET /api/system/auth/get-permission-info
 * 后端返回 PermissionInfoResp，这里转换为前端 AuthPermissionInfo 格式
 */
export async function getAuthPermissionInfoApi(): Promise<AuthPermissionInfo> {
  const raw = await requestClient.get<any>(
    '/system/auth/get-permission-info',
  );
  return transformPermissionInfo(raw);
}
