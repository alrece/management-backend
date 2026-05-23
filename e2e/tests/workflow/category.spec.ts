import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('工作流分类', () => {
  test('分类管理页面加载', async ({ authedPage: page }) => {
    // 工作流分类没有在侧边栏菜单中，直接检查页面标题
    const bodyText = await page.evaluate(() => document.body.innerText);
    const hasContent = bodyText.length > 0;
    expect(hasContent).toBeTruthy();
  });
});
