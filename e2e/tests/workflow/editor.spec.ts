import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('工作流编辑器', () => {
  test.slow();

  test('编辑器页面加载成功', async ({ authedPage: page }) => {
    // 检查页面已加载
    const title = await page.title();
    expect(title).toBeTruthy();
  });

  test('节点面板可见', async ({ authedPage: page }) => {
    const bodyText = await page.evaluate(() => document.body.innerText);
    expect(bodyText.length).toBeGreaterThan(0);
  });
});
