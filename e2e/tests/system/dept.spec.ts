import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('部门管理', () => {
  test.beforeEach(async ({ authedPage: page }) => {
    await expandAndClick(page, '系统管理', '部门管理');
  });

  test('部门管理页面加载', async ({ authedPage: page }) => {
    const container = page.locator('.vxe-grid, .vxe-table, .ant-tree, [class*="tree"]');
    await expect(container.first()).toBeVisible({ timeout: 10000 });
  });
});
