import { requestClient } from '#/api/request';

export namespace SystemUserProfileApi {
  export interface UserProfile {
    id: number;
    username: string;
    nickname: string;
    email?: string;
    mobile?: string;
    sex?: number;
    avatar?: string;
    createTime: string;
  }

  export interface UpdatePasswordReq {
    oldPassword: string;
    newPassword: string;
  }

  export interface UpdateProfileReq {
    nickname?: string;
    email?: string;
    mobile?: string;
    sex?: number;
  }
}

/** 获取个人信息 GET /api/system/user/profile */
export function getUserProfile() {
  return requestClient.get<SystemUserProfileApi.UserProfile>(
    '/system/user/profile',
  );
}

/** 修改个人信息 PUT /api/system/user/profile */
export function updateUserProfile(
  data: SystemUserProfileApi.UpdateProfileReq,
) {
  return requestClient.put('/system/user/profile', data);
}

/** 修改密码 PUT /api/system/user/password */
export function updateUserPassword(
  data: SystemUserProfileApi.UpdatePasswordReq,
) {
  return requestClient.put('/system/user/password', data);
}
