import { defineConfig } from 'vite';

export default defineConfig({
	build: {
		outDir: 'build', // 指定输出目录
		assetsDir: 'static', // 静态资源目录
		sourcemap: true, // 生成 source map 文件
		minify: false, // 不压缩代码
		cssMinify: false,
		brotliSize: false, // 不使用 brotli 压缩
		terserOptions: {
      compress: false // 不使用 terser 压缩
    },
		rollupOptions: {
			output: {
				inlineDynamicImports: false,
				manualChunks: (id) => {
					// Generate a separate chunk for each file
					const moduleId = id.split('/').pop().split('.')[0];
					return moduleId;
				}
			}
		}
	},
	base: './' // 使用相对路径
});