import { requestClient } from '#/api/request';

export namespace CategoryApi {
  /** 分类树节点 */
  export interface TreeNode {
    id: number;
    name: string;
    parentId: number;
    sort: number;
    status: number;
    children?: TreeNode[];
  }

  /** 分类表单 */
  export interface CategoryForm {
    id?: number;
    name: string;
    parentId: number;
    sort: number;
    status?: number;
  }
}

/** 获取分类树 */
export function getCategoryTree() {
  return requestClient.get<CategoryApi.TreeNode[]>('/workflow/category/tree');
}

/** 创建分类 */
export function createCategory(data: CategoryApi.CategoryForm) {
  return requestClient.post('/workflow/category', data);
}

/** 更新分类 */
export function updateCategory(data: CategoryApi.CategoryForm) {
  return requestClient.put('/workflow/category', data);
}

/** 删除分类 */
export function deleteCategory(id: number) {
  return requestClient.delete(`/workflow/category/${id}`);
}
