import { spawn } from 'node:child_process';
import { fileURLToPath } from 'node:url';
const site = fileURLToPath(new URL('../', import.meta.url));
const server = spawn(process.execPath, ['scripts/serve.mjs'], {cwd: site, windowsHide: true, stdio: ['ignore', 'inherit', 'inherit', 'ipc']});
try {
 await new Promise((resolve,reject) => {
  const timeout = setTimeout(() => reject(new Error('Static test server did not become ready')), 60_000);
  server.once('message', message => { if (message === 'ready') { clearTimeout(timeout); resolve(); } });
  server.once('error', error => { clearTimeout(timeout); reject(error); });
  server.once('exit', () => { clearTimeout(timeout); reject(new Error('Static test server could not start; check port 4174')); });
 });
 const code = await new Promise((resolve,reject) => {
  const tests = spawn(process.execPath, ['node_modules/@playwright/test/cli.js', 'test'], {cwd:site, windowsHide:true, stdio:['ignore','inherit','inherit']});
  tests.on('error',reject); tests.on('exit',resolve);
 });
 process.exitCode = code ?? 1;
} finally { server.kill(); }
