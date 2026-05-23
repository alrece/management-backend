import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('菜单管理', () => {
  test.beforeEach(async ({ authedPage: page }) => {
    await expandAndClick(page, '系统管理', '菜单管理');
  });

  test('菜单管理页面加载', async ({ authedPage: page }) => {
    const container = page.locator('.vxe-grid, .vxe-table, .ant-tree, [class*="tree"]');
    await expect(container.first()).toBeVisible({ timeout: 10000 });
  });

  test('显示系统管理目录', async ({ authedPage: page }) => {
    const sysText = page.locator('text=系统管理');
    await expect(sysText.first()).toBeVisible({ timeout: 10000 });
  });
});
