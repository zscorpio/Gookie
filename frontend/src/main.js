import './style.css';
import './app.css';

import {CreateBrowser, GetCookie} from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
    <div class="container">
      <div class="button-group">
        <button class="btn" id="newBrowserBtn">新建浏览器</button>
        <button class="btn" id="getCookieBtn">获取Cookie</button>
      </div>
      <div class="cookie-display" id="cookieDisplay">
        <div class="input-group">
          <input class="input" id="cookieInput" type="text" readonly placeholder="这里将显示Cookie" />
          <button class="copy-btn" id="copyBtn" title="复制Cookie">复制</button>
        </div>
      </div>
    </div>
`;

// 设置"新建浏览器"按钮的事件处理
document.getElementById("newBrowserBtn").addEventListener("click", function() {
  CreateBrowser()
    .then((result) => {
      alert(result); // 显示创建结果
    })
    .catch((err) => {
      console.error(err);
      alert("创建浏览器失败");
    });
});

// 设置"获取Cookie"按钮的事件处理
document.getElementById("getCookieBtn").addEventListener("click", function() {
  // 切换cookie展示区域的显示状态
  const cookieDisplay = document.getElementById("cookieDisplay");
  
  if (!cookieDisplay.classList.contains('visible')) {
    // 显示cookie区域
    cookieDisplay.classList.add('visible');
    
    // 获取cookie
    GetCookie()
      .then((result) => {
        document.getElementById("cookieInput").value = result;
      })
      .catch((err) => {
        console.error(err);
        document.getElementById("cookieInput").value = "获取Cookie失败";
      });
  } else {
    // 隐藏cookie区域
    cookieDisplay.classList.remove('visible');
  }
});

// 设置"复制"按钮的事件处理
document.getElementById("copyBtn").addEventListener("click", function() {
  const cookieInput = document.getElementById("cookieInput");
  
  // 选中输入框内容
  cookieInput.select();
  cookieInput.setSelectionRange(0, 99999); // 适用于移动设备
  
  // 复制内容到剪贴板
  try {
    navigator.clipboard.writeText(cookieInput.value)
      .then(() => {
        // 显示复制成功的视觉反馈
        const copyBtn = document.getElementById("copyBtn");
        const originalText = copyBtn.textContent;
        
        copyBtn.textContent = "已复制!";
        copyBtn.classList.add("copied");
        
        setTimeout(() => {
          copyBtn.textContent = originalText;
          copyBtn.classList.remove("copied");
        }, 1500);
      })
      .catch(err => {
        console.error('复制失败:', err);
        alert("复制失败，请手动复制");
      });
  } catch (err) {
    // 如果浏览器不支持clipboard API，使用传统方式
    document.execCommand('copy');
    
    // 显示复制成功的视觉反馈
    const copyBtn = document.getElementById("copyBtn");
    const originalText = copyBtn.textContent;
    
    copyBtn.textContent = "已复制!";
    copyBtn.classList.add("copied");
    
    setTimeout(() => {
      copyBtn.textContent = originalText;
      copyBtn.classList.remove("copied");
    }, 1500);
  }
});
