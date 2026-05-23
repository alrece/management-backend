import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('用户管理', () => {
  test.beforeEach(async ({ authedPage: page }) => {
    await expandAndClick(page, '系统管理', '用户管理');
  });

  test('用户列表页面加载', async ({ authedPage: page }) => {
    const grid = page.locator('.vxe-grid, .vxe-table, .vxe-body--row');
    await expect(grid.first()).toBeVisible({ timeout: 10000 });
  });

  test('用户列表显示 admin 用户', async ({ authedPage: page }) => {
    const adminText = page.locator('text=admin');
    await expect(adminText.first()).toBeVisible({ timeout: 10000 });
  });
});
