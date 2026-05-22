<script lang="ts" setup>
import type { CategoryApi } from '#/api/workflow/category';

import { onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { message, Modal, Tree } from 'ant-design-vue';

import { deleteCategory, getCategoryTree } from '#/api/workflow/category';

import Form from './modules/form.vue';

const [FormModal, formModalApi] = useVbenModal({
  connectedComponent: Form,
  destroyOnClose: true,
});

const treeData = ref<CategoryApi.TreeNode[]>([]);
const loading = ref(false);

async function loadTree() {
  loading.value = true;
  try {
    treeData.value = await getCategoryTree();
  } finally {
    loading.value = false;
  }
}

function handleCreate(parentId?: number) {
  formModalApi.setData({ parentId: parentId ?? 0, treeData: treeData.value }).open();
}

function handleEdit(node: CategoryApi.TreeNode) {
  formModalApi.setData({ ...node, treeData: treeData.value }).open();
}

function handleDelete(node: CategoryApi.TreeNode) {
  if (node.children?.length) {
    message.warning('请先删除子分类');
    return;
  }
  Modal.confirm({
    title: `确认删除分类「${node.name}」吗？`,
    async onOk() {
      await deleteCategory(node.id);
      message.success('删除成功');
      await loadTree();
    },
  });
}

const onSelect = (_keys: any[], _info: { node: CategoryApi.TreeNode }) => {
  // 预留：点击节点可跳转或筛选
};

onMounted(loadTree);
</script>

<template>
  <Page auto-content-height>
    <FormModal @success="loadTree" />
    <div class="p-4">
      <div class="mb-4 flex items-center justify-between">
        <h3 class="text-lg font-medium">流程分类</h3>
        <a-button type="primary" @click="handleCreate()">新建分类</a-button>
      </div>
      <a-spin :spinning="loading">
        <Tree
          :tree-data="treeData"
          :field-names="{ title: 'name', key: 'id', children: 'children' }"
          default-expand-all
          @select="onSelect"
        >
          <template #title="node">
            <div class="flex items-center justify-between gap-2 py-1">
              <span>{{ node.name }}</span>
              <span class="flex gap-1">
                <a-button type="link" size="small" @click.stop="handleCreate(node.id)">添加</a-button>
                <a-button type="link" size="small" @click.stop="handleEdit(node)">编辑</a-button>
                <a-button type="link" size="small" danger @click.stop="handleDelete(node)">删除</a-button>
              </span>
            </div>
          </template>
        </Tree>
        <a-empty v-if="!loading && !treeData.length" description="暂无分类数据" />
      </a-spin>
    </div>
  </Page>
</template>
