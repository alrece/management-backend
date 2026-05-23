import { test, expect } from '@playwright/test';

test.describe('登录', () => {
  test('admin 登录成功，跳转首页', async ({ page }) => {
    await page.goto('/#/auth/login');
    await page.waitForSelector('input[placeholder="请输入用户名"]', { timeout: 10000 });

    const usernameInput = page.locator('input[placeholder="请输入用户名"]').first();
    const passwordInput = page.locator('input[placeholder="请输入密码"]').first();

    await usernameInput.fill('admin');
    await passwordInput.fill('admin123');

    const loginBtn = page.locator('button:has-text("登录")').first();
    await loginBtn.click();

    // hash routing: 检查 hash 不含 login 即表示跳转成功
    await page.waitForFunction(
      () => !window.location.hash.includes('login'),
      { timeout: 15000 },
    );
    expect(new URL(page.url()).hash).not.toContain('login');
  });

  test('错误密码登录失败', async ({ page }) => {
    await page.goto('/#/auth/login');
    await page.waitForSelector('input[placeholder="请输入用户名"]', { timeout: 10000 });

    const usernameInput = page.locator('input[placeholder="请输入用户名"]').first();
    const passwordInput = page.locator('input[placeholder="请输入密码"]').first();

    await usernameInput.fill('admin');
    await passwordInput.fill('wrongpassword');

    const loginBtn = page.locator('button:has-text("登录")').first();
    await loginBtn.click();

    await page.waitForTimeout(2000);
    const stillOnLogin = page.url().includes('login');
    const hasError = await page.locator('text=/错误|失败|incorrect|invalid/i').count() > 0;
    expect(stillOnLogin || hasError).toBeTruthy();
  });
});
