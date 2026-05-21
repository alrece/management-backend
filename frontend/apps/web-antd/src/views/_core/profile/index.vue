<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { Page } from '@vben/common-ui';
import { useUserStore } from '@vben/stores';

import { Card, Tabs } from 'ant-design-vue';

import { getAuthPermissionInfoApi } from '#/api';

import BaseInfo from './modules/base-info.vue';
import ProfileUser from './modules/profile-user.vue';
import ResetPwd from './modules/reset-pwd.vue';

const userStore = useUserStore();
const activeName = ref('basicInfo');

const profile = ref<any>();

/** 刷新个人信息 */
async function refreshProfile() {
  const authPermissionInfo = await getAuthPermissionInfoApi();
  userStore.setUserInfo(authPermissionInfo.user);
  profile.value = authPermissionInfo.user;
}

onMounted(refreshProfile);
</script>

<template>
  <Page auto-content-height>
    <div class="flex">
      <Card class="w-2/5" title="个人信息">
        <ProfileUser :profile="profile" @success="refreshProfile" />
      </Card>
      <Card class="ml-3 w-3/5">
        <Tabs v-model:active-key="activeName" class="-mt-4">
          <Tabs.TabPane key="basicInfo" tab="基本设置">
            <BaseInfo :profile="profile" @success="refreshProfile" />
          </Tabs.TabPane>
          <Tabs.TabPane key="resetPwd" tab="密码设置">
            <ResetPwd />
          </Tabs.TabPane>
        </Tabs>
      </Card>
    </div>
  </Page>
</template>
