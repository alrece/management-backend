/**
 * 简化版上传 hook（原型阶段）
 */
export function useUpload() {
  return {
    accept: 'image/*',
    uploadUrl: '/api/system/file',
  };
}
