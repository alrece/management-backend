import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('角色管理', () => {
  test.beforeEach(async ({ authedPage: page }) => {
    await expandAndClick(page, '系统管理', '角色管理');
  });

  test('角色列表页面加载', async ({ authedPage: page }) => {
    const grid = page.locator('.vxe-grid, .vxe-table, .vxe-body--row');
    await expect(grid.first()).toBeVisible({ timeout: 10000 });
  });

  test('角色列表显示搜索表单', async ({ authedPage: page }) => {
    const bodyText = await page.evaluate(() => document.body.innerText);
    const hasRoleContent = bodyText.includes('角色名称') || bodyText.includes('角色标识');
    expect(hasRoleContent).toBeTruthy();
  });
});
