import { requestClient } from '#/api/request';

export namespace SystemDeptApi {
  /** 部门信息（匹配后端 DeptTreeResp） */
  export interface Dept {
    id?: number;
    parentId: number;
    deptName: string;
    ancestors?: string;
    sort: number;
    leader: string;
    status: number;
    children?: Dept[];
  }
}

/** 部门树 */
export function getDeptTree() {
  return requestClient.get<SystemDeptApi.Dept[]>('/system/dept/tree');
}

/** 新增部门 */
export function createDept(data: SystemDeptApi.Dept) {
  return requestClient.post('/system/dept', data);
}

/** 修改部门 */
export function updateDept(data: SystemDeptApi.Dept) {
  return requestClient.put('/system/dept', data);
}

/** 删除部门 */
export function deleteDept(id: number) {
  return requestClient.delete(`/system/dept/${id}`);
}

/** 兼容旧引用：部门列表 = 部门树 */
export const getSimpleDeptList = getDeptTree;
export const getDeptList = getDeptTree;
export function getDept(_id: number) {
  return Promise.resolve(null);
}
export function deleteDeptList(_ids: number[]) {
  return Promise.resolve();
}