import { test as base, expect, type Page } from '@playwright/test';

type AuthFixtures = {
  authedPage: Page;
};

async function login(page: Page) {
  await page.goto('/#/auth/login');
  await page.waitForSelector('input[placeholder="请输入用户名"]', { timeout: 10000 });
  const usernameInput = page.locator('input[placeholder="请输入用户名"]').first();
  const passwordInput = page.locator('input[placeholder="请输入密码"]').first();
  await usernameInput.fill('admin');
  await passwordInput.fill('admin123');
  const loginBtn = page.locator('button:has-text("登录")').first();
  await loginBtn.click();
  await page.waitForFunction(
    () => !window.location.hash.includes('login') && window.location.hash.length > 0,
    { timeout: 15000 },
  );
}

async function clickMenuItem(page: Page, menuLabel: string) {
  const item = page.locator(`[class*="menu"] a:has-text("${menuLabel}"), [class*="menu"] [role=menuitem]:has-text("${menuLabel}")`).first();
  await item.click();
  await page.waitForTimeout(3000);
}

async function expandAndClick(page: Page, parentLabel: string, childLabel: string) {
  const parent = page.locator(`[class*="menu"] :text-is("${parentLabel}"), [class*="sidebar"] :text-is("${parentLabel}")`).first();
  if (await parent.isVisible()) {
    await parent.click();
    await page.waitForTimeout(1000);
  }
  await clickMenuItem(page, childLabel);
}

export const test = base.extend<AuthFixtures>({
  authedPage: async ({ page }, use) => {
    await login(page);
    await use(page);
  },
});

export { expect, login, clickMenuItem, expandAndClick };
