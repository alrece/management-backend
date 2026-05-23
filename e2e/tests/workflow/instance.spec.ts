import { test, expect, expandAndClick } from '../../fixtures/auth';

test.describe('执行实例', () => {
  test('执行实例页面加载', async ({ authedPage: page }) => {
    const bodyText = await page.evaluate(() => document.body.innerText);
    expect(bodyText.length).toBeGreaterThan(0);
  });
});
