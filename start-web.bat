@echo off
chcp 65001 >nul
echo 正在启动 Web 前端服务...
echo.

cd web/file-flow-web
npm run dev

echo.
echo Web 服务已启动在 http://localhost:3000
echo.
echo 请打开浏览器访问 http://localhost:3000
echo.
pause
