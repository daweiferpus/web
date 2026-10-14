# 数据监控面板

一个简单的前端应用，用于定时发送HTTP请求并以表格形式展示JSON数据。

## 功能

- 发送并接收HTTP请求
- 解析JSON格式数据并以表格形式展示
- 每3秒自动刷新一次数据
- 支持启动/停止监控
- 显示连接状态和刷新统计

## 部署说明

### 前提条件

- 一个支持静态网站托管的云服务器（如Nginx, Apache等）

### 部署步骤

1. 将`public`目录中的所有文件上传到云服务器的web根目录

   ```
   scp -r public/* user@your-server:/var/www/html/
   ```

2. 配置您的web服务器指向这些文件

   **Nginx配置示例:**

   ```nginx
   server {
       listen 80;
       server_name your-domain.com;
       root /var/www/html;
       index index.html;

       location / {
           try_files $uri $uri/ =404;
       }
   }
   ```

   **Apache配置示例:**

   ```apache
   <VirtualHost *:80>
       ServerName your-domain.com
       DocumentRoot /var/www/html
       <Directory /var/www/html>
           Options Indexes FollowSymLinks
           AllowOverride All
           Require all granted
       </Directory>
   </VirtualHost>
   ```

3. 重启web服务器

   ```
   # Nginx
   sudo systemctl restart nginx

   # Apache
   sudo systemctl restart apache2
   ```

## 使用说明

1. 打开浏览器访问您的网站
2. 在"API地址"字段中输入您想要请求的API URL
3. 点击"开始监控"按钮开始定时请求数据
4. 数据将以表格形式显示，并每3秒自动刷新一次
5. 点击"停止监控"按钮停止请求

## 注意事项

- 请确保API返回的是有效的JSON格式数据
- 如果要从不同域名请求数据，API需要支持CORS或者您需要配置代理
- 应用支持两种数据格式：
  - JSON数组（包含对象）：会显示为多行表格，每个对象一行
  - JSON对象：会显示为两列表格，键和值

## 自定义

- 修改`styles.css`文件可以自定义界面样式
- 修改`app.js`文件中的`fetchData`函数可以自定义请求逻辑 