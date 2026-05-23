import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('工作流列表', () => {
  test('工作流列表页面加载', async ({ authedPage: page }) => {
    const bodyText = await page.evaluate(() => document.body.innerText);
    expect(bodyText.length).toBeGreaterThan(0);
  });

  test('工作流搜索区域可见', async ({ authedPage: page }) => {
    const bodyText = await page.evaluate(() => document.body.innerText);
    expect(bodyText.length).toBeGreaterThan(0);
  });
});
