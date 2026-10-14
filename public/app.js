// 全局变量
let intervalId = null;
let refreshCount = 0;
let isMonitoring = false;
let intervalIdLimit = null;
let isLimit = false;

let refreshCountLimit = 0;

// DOM元素
const apiUrlInput = document.getElementById('api-url');
const limitUrlInput = document.getElementById('limit-url');

const startBtn = document.getElementById('start-btn');
const stopBtn = document.getElementById('stop-btn');
const statusText = document.getElementById('status-text');
const statusLight = document.getElementById('status-light');
const refreshCountElement = document.getElementById('refresh-count');
const refreshCountElementLimit = document.getElementById('refresh-limit');
const lastUpdateElement = document.getElementById('last-update');
const lastUpdateElementLimit = document.getElementById('last-limit');
const tableHeader = document.getElementById('table-header');
const tableBody = document.getElementById('table-body');

// 新增DOM元素
const bidNumInput = document.getElementById('bid-num');
const bidBtn = document.getElementById('bid-btn');
const totalAmountElement = document.getElementById('total-amount');
const bidTableHeader = document.getElementById('bid-table-header');
const bidTableBody = document.getElementById('bid-table-body');

// 当日活跃度
const activeNumInput = document.getElementById('active-num');
// const activeTab = document.getElementById('active-tab');
const activeBtn = document.getElementById('active-btn');
const activeTableHeader = document.getElementById('active-table-header');
const activeTableBody = document.getElementById('active-table-body');

// 连板股票
const limitTableHeader = document.getElementById('limit-table-header');
const limitTableBody = document.getElementById('limit-table-body');
const startBtnLimit = document.getElementById('start-btn-limit');
const stopBtnLimit = document.getElementById('stop-btn-limit');

// 注册事件监听
startBtn.addEventListener('click', startMonitoring);
stopBtn.addEventListener('click', stopMonitoring);
bidBtn.addEventListener('click', monitorBid);
activeBtn.addEventListener('click', getActiveData);
startBtnLimit.addEventListener('click', startLimit);
stopBtnLimit.addEventListener('click', stopLimit);

// 开始监控
function startMonitoring() {
    const apiUrl = apiUrlInput.value.trim();
    if (!apiUrl) {
        alert('请输入有效的API地址');
        return;
    }

    // 更新UI状态
    isMonitoring = true;
    startBtn.disabled = true;
    stopBtn.disabled = false;

    // 立即发起第一次请求
    fetchData(apiUrl)
        .then(() => {
            // 设置定时器，每1.5秒请求一次数据
            intervalId = setInterval(() => fetchData(apiUrl), 1500);
        })
        .catch(error => {
            console.error('初始连接失败:', error);
            updateStatus('连接失败', 'offline');
            stopMonitoring();
        });
}

// 停止监控
function stopMonitoring() {
    if (intervalId) {
        clearInterval(intervalId);
        intervalId = null;
    }
    
    isMonitoring = false;
    startBtn.disabled = false;
    stopBtn.disabled = true;
    updateStatus('监控已停止', 'offline');
}

// 开始监控
function startLimit() {
    const apiUrl = limitUrlInput.value.trim();
    if (!apiUrl) {
        alert('请输入有效的API地址');
        return;
    }

    // 更新UI状态
    isLimit = true;
    startBtnLimit.disabled = true;
    stopBtnLimit.disabled = false;
    // 立即发起第一次请求
    fetchDataLimit(apiUrl)
        .then(() => {
            // 设置定时器，每1.5秒请求一次数据
            intervalIdLimit = setInterval(() => fetchDataLimit(apiUrl), 1500);
        })
        .catch(error => {
            console.error('初始连接失败:', error);
            stopBtnLimit();
        });
}

// 停止监控
function stopLimit() {
    if (intervalIdLimit) {
        clearInterval(intervalIdLimit);
        intervalIdLimit = null;
    }
    
    isLimit = false;
    startBtnLimit.disabled = false;
    stopBtnLimit.disabled = true;
}

// 更新状态指示器
function updateStatus(message, status) {
    // statusText.textContent = message;
    // statusLight.className = `status-light ${status}`;
}

// 获取数据
async function fetchData(url) {
    try {
        updateStatus('正在获取数据...', 'pending');
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            mode: 'cors' // 明确指定跨域模式
        });
        
        if (!response.ok) {
            throw new Error(`HTTP错误! 状态: ${response.status}`);
        }
        
        const responseJson = await response.json();
        // console.log('收到的响应数据:', responseJson); // 添加调试日志
        
        // 检查后端返回的业务状态码
        if (responseJson.code !== 0) {
            throw new Error(`业务错误! 信息: ${responseJson.message}`);
        }
        
        // 从响应中提取真正的数据
        const data = responseJson.data;
        
        if (!Array.isArray(data)) {
            throw new Error('返回的数据格式不正确，期望数组类型');
        }
        
        // 更新UI
        refreshCount++;
        refreshCountElement.textContent = refreshCount;
        lastUpdateElement.textContent = new Date().toLocaleTimeString();
        updateStatus('已连接', 'online');
        
        // 渲染数据表格
        renderTable(data);
    } catch (error) {
        console.error('获取数据失败:', error);
        updateStatus(`获取数据失败: ${error.message}`, 'offline');
    }
}

// 渲染数据表格
function renderTable(data) {
    // 清空现有表格
    tableHeader.innerHTML = '';
    tableBody.innerHTML = '';
    
    // 处理数组数据
    if (Array.isArray(data)) {
        if (data.length === 0) {
            showEmptyTable('没有数据');
            return;
        }
        
        // 获取第一个对象的所有键作为表头
        const firstItem = data[0];
        const headers = Object.keys(firstItem);
        const headers2 = ["代码", "名称", "涨幅", "价格", "涨停原因", "活跃度", "流通", "成交额", "换手"];
        // 创建表头
        headers2.forEach(header => {
            const th = document.createElement('th');
            th.textContent = header;
            tableHeader.appendChild(th);
        });
        
        // 创建数据行
        data.forEach(item => {
            const row = document.createElement('tr');
            
            headers.forEach(header => {
                const cell = document.createElement('td');
                cell.textContent = formatCellValue(item[header]);
                row.appendChild(cell);
            });
            // // 在页面加载完成后，发起一个 GET 请求
            // fetch('https://api.example.com/data')
            //   .then(response => response.json())
            //   .then(data => {
            //     console.log('获取到的数据:', data);
            //   });
            
            tableBody.appendChild(row);
        });
    } 
    // 处理对象数据
    else if (typeof data === 'object' && data !== null) {
        // 创建表头: 键和值
        const keyHeader = document.createElement('th');
        keyHeader.textContent = '键';
        tableHeader.appendChild(keyHeader);
        
        const valueHeader = document.createElement('th');
        valueHeader.textContent = '值';
        tableHeader.appendChild(valueHeader);
        
        // 创建数据行
        Object.entries(data).forEach(([key, value]) => {
            const row = document.createElement('tr');
            
            const keyCell = document.createElement('td');
            keyCell.textContent = key;
            row.appendChild(keyCell);
            
            const valueCell = document.createElement('td');
            valueCell.textContent = formatCellValue(value);
            row.appendChild(valueCell);
            
            tableBody.appendChild(row);
        });
    } else {
        showEmptyTable('无效的数据格式');
    }
}

// 格式化单元格值
function formatCellValue(value) {
    if (value === null || value === undefined) {
        return '-';
    } else if (typeof value === 'object') {
        return JSON.stringify(value);
    } else {
        return String(value);
    }
}

// 显示空表格消息
function showEmptyTable(message) {
    const emptyHeader = document.createElement('th');
    emptyHeader.textContent = '信息';
    tableHeader.appendChild(emptyHeader);
    
    const row = document.createElement('tr');
    const cell = document.createElement('td');
    cell.textContent = message;
    cell.style.textAlign = 'center';
    row.appendChild(cell);
    tableBody.appendChild(row);
}

// 切换tab
function openTab(evt, tabName) {
    const tabContents = document.getElementsByClassName('tab-content');
    for (let i = 0; i < tabContents.length; i++) {
        tabContents[i].style.display = 'none';
    }
    const tabButtons = document.getElementsByClassName('tab-button');
    for (let i = 0; i < tabButtons.length; i++) {
        tabButtons[i].className = tabButtons[i].className.replace(' active', '');
    }
    document.getElementById(tabName).style.display = 'block';
    evt.currentTarget.className += ' active';
}

// 监控竞价
async function monitorBid() {
    const bidNum = bidNumInput.value.trim();
    if (!bidNum) {
        alert('请输入有效的竞价编号');
        return;
    }

    const url = `http://localhost:8000/monitor/hq/bid/${bidNum}`;
    bidBtn.disabled = true; // 禁用按钮
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            mode: 'cors'
        });

        if (!response.ok) {
            throw new Error(`HTTP错误! 状态: ${response.status}`);
        }

        const responseJson = await response.json();

        if (responseJson.code !== 0) {
            throw new Error(`业务错误! 信息: ${responseJson.message}`);
        }

        const bidData = responseJson.data;
        // 假设 bidData.totalamount 是一个整数
        const totalAmount = bidData.totalamount;
        // 昨日成交量
        const yesterAmount = bidData.yesterAmount;
        // 将数字转换为以"亿"为单位的格式
        const billionValue = (totalAmount / 1e8).toFixed(2); // 转换为亿并保留两位小数（四舍五入）
        // 将数字转换为以"亿"为单位的格式
        const yesterValue = (yesterAmount / 1e8).toFixed(2); // 转换为亿并保留两位小数（四舍五入）

        const differ = ((totalAmount - yesterAmount) / 1e8).toFixed(2)
        if (differ >= 0) {
            // 设置文本内容
            totalAmountElement.textContent = `${bidData.date}日 ${billionValue}亿 较昨日放量 ${differ}亿`;
        } else {
            // 设置文本内容
            totalAmountElement.textContent = `${bidData.date}日 ${billionValue}亿 较昨日缩量 ${differ}亿`;
        }


        renderBidTable(bidData.bidhq);
    } catch (error) {
        console.error('获取竞价数据失败:', error);
        alert(`获取竞价数据失败: ${error.message}`);
    } finally {
        bidBtn.disabled = false; // 收到数据后启用按钮
    }
}

// 渲染竞价表格
function renderBidTable(data) {
    // 清空现有表格
    bidTableHeader.innerHTML = '';
    bidTableBody.innerHTML = '';

    if (!Array.isArray(data) || data.length === 0) {
        showEmptyBidTable('没有竞价数据');
        return;
    }

    // 获取第一个对象的所有键作为表头
    const firstItem = data[0];
    const headers = Object.keys(firstItem);

    // 创建表头
    headers.forEach(header => {
        const th = document.createElement('th');
        th.textContent = header;
        bidTableHeader.appendChild(th);
    });

    // 创建数据行
    data.forEach(item => {
        const row = document.createElement('tr');

        headers.forEach(header => {
            const cell = document.createElement('td');
            cell.textContent = formatCellValue(item[header]);
            row.appendChild(cell);
        });

        bidTableBody.appendChild(row);
    });
}

// 显示空竞价表格消息
function showEmptyBidTable(message) {
    const emptyHeader = document.createElement('th');
    emptyHeader.textContent = '信息';
    bidTableHeader.appendChild(emptyHeader);

    const row = document.createElement('tr');
    const cell = document.createElement('td');
    cell.textContent = message;
    cell.style.textAlign = 'center';
    row.appendChild(cell);
    bidTableBody.appendChild(row);
} 

async function getActiveData() {
    const stkNum = activeNumInput.value.trim();
    if (!stkNum) {
        alert('请输入有效的竞价编号');
        return;
    }
    const url = `http://localhost:8001/active/limit/${stkNum}`;
    activeBtn.disabled = true; // 禁用按钮
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            mode: 'cors'
        });

        if (!response.ok) {
            throw new Error(`HTTP错误! 状态: ${response.status}`);
        }

        const responseJson = await response.json();

        if (responseJson.code !== 0) {
            throw new Error(`业务错误! 信息: ${responseJson.message}`);
        }

        const activeData = responseJson.data;

        renderActiveTable(activeData);
    } catch (error) {
        console.error('获取竞价数据失败:', error);
        alert(`获取竞价数据失败: ${error.message}`);
    } finally {
         activeBtn.disabled = false; // 收到数据后启用按钮
    }
}

// 渲染竞价表格
function renderActiveTable(data) {
    // 清空现有表格
    activeTableHeader.innerHTML = '';
    activeTableBody.innerHTML = '';

    if (!Array.isArray(data) || data.length === 0) {
        showEmptyBidTable('没有竞价数据');
        return;
    }

    // 获取第一个对象的所有键作为表头
    const firstItem = data[0];
    const headers = Object.keys(firstItem);

    // 创建表头
    headers.forEach(header => {
        const th = document.createElement('th');
        th.textContent = header;
        activeTableHeader.appendChild(th);
    });

    // 创建数据行
    data.forEach(item => {
        const row = document.createElement('tr');

        headers.forEach(header => {
            const cell = document.createElement('td');
            cell.textContent = formatCellValue(item[header]);
            row.appendChild(cell);
        });

        activeTableBody.appendChild(row);
    });
}

// 获取数据
async function fetchDataLimit(url) {
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            mode: 'cors' // 明确指定跨域模式
        });
        
        if (!response.ok) {
            throw new Error(`HTTP错误! 状态: ${response.status}`);
        }
        
        const responseJson = await response.json();
        // console.log('收到的响应数据:', responseJson); // 添加调试日志
        
        // 检查后端返回的业务状态码
        if (responseJson.code !== 0) {
            throw new Error(`业务错误! 信息: ${responseJson.message}`);
        }
        
        // 从响应中提取真正的数据
        const data = responseJson.data;
        
        if (!Array.isArray(data)) {
            throw new Error('返回的数据格式不正确，期望数组类型');
        }
        
        // 更新UI
        refreshCountLimit++;
        refreshCountElementLimit.textContent = refreshCountLimit;
        lastUpdateElementLimit.textContent = new Date().toLocaleTimeString();
        
        // 渲染数据表格
        renderTableLimit(data);
    } catch (error) {
        console.error('获取数据失败:', error);
    }
}

// 渲染数据表格
function renderTableLimit(data) {
    // 清空现有表格
    limitTableHeader.innerHTML = '';
    limitTableBody.innerHTML = '';
    
    // 处理数组数据
    if (Array.isArray(data)) {
        if (data.length === 0) {
            showEmptyTableLimit('没有数据');
            return;
        }
        
        // 获取第一个对象的所有键作为表头
        const firstItem = data[0];
        const headers = Object.keys(firstItem);
        const headers2 = ["涨停天数", "股票代码", "股票名称", "涨幅", "价格", "涨停原因", "活跃度", "市值", "成交额", "换手率"];
        // 创建表头
        headers2.forEach(header => {
            const th = document.createElement('th');
            th.textContent = header;
            limitTableHeader.appendChild(th);
        });
        
        // 创建数据行
        data.forEach(item => {
            const row = document.createElement('tr');
            
            headers.forEach(header => {
                const cell = document.createElement('td');
                cell.textContent = formatCellValue(item[header]);
                row.appendChild(cell);
            });
            // // 在页面加载完成后，发起一个 GET 请求
            // fetch('https://api.example.com/data')
            //   .then(response => response.json())
            //   .then(data => {
            //     console.log('获取到的数据:', data);
            //   });
            
            limitTableBody.appendChild(row);
        });
    } else {
        showEmptyTableLimit('无效的数据格式');
    }
}

// 显示空表格消息
function showEmptyTableLimit(message) {
    const emptyHeader = document.createElement('th');
    emptyHeader.textContent = '信息';
    LimitTableHeader.appendChild(emptyHeader);
    
    const row = document.createElement('tr');
    const cell = document.createElement('td');
    cell.textContent = message;
    cell.style.textAlign = 'center';
    row.appendChild(cell);
    limitTableBody.appendChild(row);
}