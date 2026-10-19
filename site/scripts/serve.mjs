import http from 'node:http';
import fs from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../out');
const types = {'.html':'text/html; charset=utf-8','.js':'text/javascript','.css':'text/css','.json':'application/json','.svg':'image/svg+xml','.png':'image/png','.woff2':'font/woff2','.webmanifest':'application/manifest+json','.ico':'image/x-icon'};
const server = http.createServer(async (request, response) => {
 try {
  const url = new URL(request.url, 'http://localhost');
  let file = path.resolve(root, '.' + decodeURIComponent(url.pathname));
  if (file !== root && !file.startsWith(root + path.sep)) { response.writeHead(403).end(); return; }
  try { if ((await fs.stat(file)).isDirectory()) file = path.join(file, 'index.html'); } catch {}
  let body;
  try { body = await fs.readFile(file); } catch { file = path.join(root, '404.html'); body = await fs.readFile(file); response.statusCode = 404; }
  response.setHeader('Content-Type', types[path.extname(file)] || 'application/octet-stream'); response.end(body);
 } catch { response.writeHead(400).end('Invalid path'); }
});
server.listen(Number(process.env.PORT || 4174), '127.0.0.1', () => { console.log('Static website ready on loopback'); process.send?.('ready'); });
for (const signal of ['SIGINT', 'SIGTERM']) process.on(signal, () => server.close());
