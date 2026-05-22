import { requestClient } from '#/api/request';

export namespace CategoryApi {
  /** 分类树节点 */
  export interface TreeNode {
    id: number;
    name: string;
    parentId: number;
    sort: number;
    children?: TreeNode[];
  }

  /** 分类表单 */
  export interface CategoryForm {
    name: string;
    parentId: number;
    sort: number;
  }
}

/** 获取分类树 */
export function getCategoryTree() {
  return requestClient.get<CategoryApi.TreeNode[]>('/workflow/categories/tree');
}

/** 创建分类 */
export function createCategory(data: CategoryApi.CategoryForm) {
  return requestClient.post('/workflow/categories', data);
}

/** 更新分类 */
export function updateCategory(id: number, data: CategoryApi.CategoryForm) {
  return requestClient.put(`/workflow/categories/${id}`, data);
}

/** 删除分类 */
export function deleteCategory(id: number) {
  return requestClient.delete(`/workflow/categories/${id}`);
}
