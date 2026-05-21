/**
 * Vite Mock 插件
 * 当 VITE_NITRO_MOCK=true 时拦截 API 请求并返回 mock 数据
 */
import type { PluginOption } from 'vite';

export function viteMockPlugin(): PluginOption {
  return {
    name: 'vite:mock-api',
    enforce: 'pre',
    async configureServer(server) {
      const enabled = process.env.VITE_NITRO_MOCK === 'true';
      if (!enabled) return;

      const { handleMockRequest } = await import('./mock/handlers');

      server.middlewares.use(async (req, res, next) => {
        if (!req.url?.startsWith('/api/')) {
          next();
          return;
        }

        // 处理 CORS preflight
        if (req.method === 'OPTIONS') {
          res.setHeader('Access-Control-Allow-Origin', '*');
          res.setHeader('Access-Control-Allow-Methods', 'GET,POST,PUT,DELETE,PATCH');
          res.setHeader('Access-Control-Allow-Headers', 'Content-Type,Authorization');
          res.statusCode = 204;
          res.end();
          return;
        }

        const handled = await handleMockRequest(req, res);
        if (!handled) next();
      });

      console.log('[Mock] API mock server enabled');
    },
  };
}
