import { test, expect } from '@playwright/test';
import { readFileSync } from 'node:fs';

const selectedComparison = JSON.parse(readFileSync(new URL('../.source/comparison.json', import.meta.url), 'utf8')) as { tag: string };

test('home, keyboard entry, persistent appearance and public resources', async ({ page }) => {
 const errors: string[] = []; page.on('pageerror', error => errors.push(error.message));
 await page.goto('/');
 await expect(page.getByRole('heading', {level:1})).toContainText('Put your workflows on schedule');
 await page.keyboard.press('Tab'); await expect(page.getByRole('link', {name:'Skip to content'})).toBeFocused();
 await page.keyboard.press('Enter'); await expect(page.locator('#main-content')).toBeFocused();
 const toggle=page.getByRole('button', {name:/Switch to .* mode/});
 const old=await toggle.getAttribute('aria-label'); await toggle.click();
 await expect(toggle).not.toHaveAttribute('aria-label', old!);
 const selected=await toggle.getAttribute('aria-label'); await page.reload();
 await expect(toggle).toHaveAttribute('aria-label',selected!);
 await expect(page.getByRole('link',{name:'ShruggieTech',exact:true})).toHaveAttribute('href','https://shruggie.tech');
 await expect(page.getByRole('link',{name:'Latest release',exact:true})).toHaveAttribute('href','https://github.com/shruggietech/go-schedule/releases/latest');
 expect(errors).toEqual([]);
});

test('static documentation search and navigation work without a server API',async({page})=>{
 await page.goto('/docs/');
 await page.getByRole('button',{name:/Search/}).first().click();
 const search=page.getByRole('dialog').getByRole('combobox'); await search.fill('timezone');
 await expect(page.getByRole('dialog').getByRole('option').first()).toBeVisible();
 await search.press('Enter');
 await expect(page.getByRole('dialog')).not.toBeVisible();
 await expect(page).toHaveURL(/\/docs\/.+/);
 await expect(page.getByRole('heading',{level:1})).toBeVisible();
});

test('legacy alias preserves query and fragment; missing routes have recovery',async({page})=>{
 await page.goto('/INSTALL-windows.html?origin=legacy#installation');
 await expect(page).toHaveURL(/\/docs\/install\/windows\/\?origin=legacy#installation$/);
 await page.goto('/does-not-exist/');
 await expect(page.getByRole('heading',{level:1})).toContainText('This page is no longer here');
 await expect(page.getByRole('link',{name:/documentation/i}).first()).toBeVisible();
});

for(const width of [390,768,1440])test('responsive home and comparison at '+width+'px',async({page})=>{
 await page.setViewportSize({width,height:900});
 for(const route of ['/','/compare/','/docs/install/windows/']){
  await page.goto(route); await expect(page.getByRole('heading',{level:1})).toBeVisible();
  expect(await page.evaluate(()=>document.documentElement.scrollWidth<=window.innerWidth)).toBe(true);
 }
 await page.goto('/compare/'); await expect(page.getByRole('table')).toBeVisible();
 await expect(page.getByText(selectedComparison.tag,{exact:true}).first()).toBeVisible();
 await expect(page.getByRole('columnheader')).toHaveCount(9);
});

test('diagrams render and appearance works when preference storage is unavailable',async({page})=>{
 await page.addInitScript(()=>{Object.defineProperty(window,'localStorage',{get(){throw new Error('Storage unavailable in test');}});});
 await page.goto('/docs/website/');
 await expect(page.locator('.diagram svg')).toBeVisible();
 await page.goto('/');
 const toggle=page.getByRole('button',{name:/Switch to .* mode/});
 const before=await toggle.getAttribute('aria-label'); await toggle.click();
 await expect(toggle).not.toHaveAttribute('aria-label',before!);
});

test('capture the production review surfaces',async({page})=>{
 await page.setViewportSize({width:1440,height:1000});
 for(const mode of ['dark','light'] as const){
  await page.emulateMedia({colorScheme:mode}); await page.goto('/');
  await expect(page.locator('html')).toHaveClass(mode);
  await page.screenshot({path:'../.test-temp/s106-home-'+mode+'.png',fullPage:true});
  if(mode==='dark') await page.screenshot({path:'../.test-temp/s106-home-hero.png'});
 }
 await page.goto('/docs/website/'); await expect(page.locator('.diagram svg')).toBeVisible();
 await page.screenshot({path:'../.test-temp/s106-documentation.png',fullPage:true});
 await page.goto('/compare/'); await page.locator('.comparison-filters').scrollIntoViewIfNeeded();
 await page.screenshot({path:'../.test-temp/s106-comparison.png'});
 await page.setViewportSize({width:390,height:900}); await page.goto('/');
 await page.screenshot({path:'../.test-temp/s106-mobile.png',fullPage:true});
});

test('comparison filters retain go-schedule and narrow the question',async({page})=>{
 await page.goto('/compare/');
 await page.getByRole('combobox',{name:'Compare with',exact:true}).selectOption('opcon');
 await expect(page.getByRole('columnheader')).toHaveCount(3);
 await expect(page.getByRole('columnheader',{name:'go-schedule',exact:true})).toBeVisible();
 await page.getByRole('combobox',{name:'Capability category',exact:true}).selectOption('Scheduling');
 await expect(page.getByRole('rowheader')).toHaveCount(1);
 await page.getByRole('button',{name:'Reset filters'}).click();
 await expect(page.getByRole('columnheader')).toHaveCount(9);
 await expect(page.getByRole('rowheader')).toHaveCount(10);
});

test('old home sections and original document headings retain destinations',async({page})=>{
 await page.goto('/?origin=legacy#reporting-problems');
 await expect(page).toHaveURL(/\/docs\/\?origin=legacy#reporting-problems$/);
 await expect(page.locator('#reporting-problems')).toBeVisible();
 await page.goto('/INSTALL-windows.html#installing-go-schedule-on-windows');
 await expect(page.locator('h1')).toHaveCount(1);
 await expect(page.locator('h1')).toHaveAttribute('id','installing-go-schedule-on-windows');
});
