import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), "");
    return {
        plugins: [react(), tailwindcss()],
        server: {
            proxy: {
                "/api": {
                    target: env.SERVER_URL,
                    secure: false,
                    rewrite: (path) => path.replace(/^\/api/, "")
                },
                "^/[a-zA-Z0-9]{6}/?$": {
                    target: env.SERVER_URL,
                    secure: false,
                    rewrite: (path) => path
                }
            }
        }
    };
});
