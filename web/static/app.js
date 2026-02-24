let login_token;

function logout() {
    const logoutBtn = document.querySelector('.Btn');

    logoutBtn.classList.add('loading');
    if (confirm('确定要退出登录吗？')) {
        // 清除本地存储的登录信息
        localStorage.setItem('justLoggedOut', 'true');
        localStorage.removeItem('authToken');
        localStorage.removeItem('userEmail');
        localStorage.removeItem('userId');
        localStorage.removeItem('userName');

        // 跳转到登录页
        localStorage.clear()
        window.location.href = '/';
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const sideMenu = document.querySelector("aside");
    const themeToggler = document.querySelector(".theme-toggler");
    const nextDay = document.getElementById('nextDay');
    const prevDay = document.getElementById('prevDay');

    // Profile button toggle for side menu
    // profileBtn.onclick = function () {
    //     sideMenu.classList.toggle('active');
    // }

    // Scroll event to remove side menu and add/remove header active class
    window.onscroll = () => {
        sideMenu.classList.remove('active');
        if (window.scrollY > 0) {
            document.querySelector('header').classList.add('active');
        } else {
            document.querySelector('header').classList.remove('active');
        }
    }

    // Theme toggle function
    const applySavedTheme = () => {
        const isDarkMode = localStorage.getItem('dark-theme') === 'true';
        if (isDarkMode) {
            document.body.classList.add('dark-theme');
            themeToggler.querySelector('span:nth-child(2)').classList.add('active');
            themeToggler.querySelector('span:nth-child(1)').classList.remove('active');
        } else {
            document.body.classList.remove('dark-theme');
            themeToggler.querySelector('span:nth-child(2)').classList.remove('active');
            themeToggler.querySelector('span:nth-child(1)').classList.add('active');
        }
    }

    // Set the initial theme based on localStorage
    applySavedTheme();

    // Toggle theme function
    themeToggler.onclick = function () {
        // Toggle dark theme class on body
        document.body.classList.toggle('dark-theme');

        // Toggle active class on the theme toggler spans
        themeToggler.querySelector('span:nth-child(1)').classList.toggle('active');
        themeToggler.querySelector('span:nth-child(2)').classList.toggle('active');

        // Save the theme preference in localStorage
        localStorage.setItem('dark-theme', document.body.classList.contains('dark-theme'));
    }

    // Function to set timetable data
    let setData = (day) => {
        document.querySelector('table tbody').innerHTML = '';  // Clear previous table data
        let daylist = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
        document.querySelector('.timetable div h2').innerHTML = daylist[day];

        // Define subjects for each day (you might need to update this with real data)
        let daySchedule = [];
        switch (day) {
            case 0: daySchedule = Sunday; break;
            case 1: daySchedule = Monday; break;
            case 2: daySchedule = Tuesday; break;
            case 3: daySchedule = Wednesday; break;
            case 4: daySchedule = Thursday; break;
            case 5: daySchedule = Friday; break;
            case 6: daySchedule = Saturday; break;
        }

        // Append timetable data to table
        daySchedule.forEach(sub => {
            const tr = document.createElement('tr');
            const trContent = `
                <td>${sub.time}</td>
                <td>${sub.roomNumber}</td>
                <td>${sub.subject}</td>
                <td>${sub.type}</td>
            `;
            tr.innerHTML = trContent;
            document.querySelector('table tbody').appendChild(tr);
        });
    }

    // Get current day and set timetable on page load
    let now = new Date();
    let today = now.getDay();  // Get current day (0 - 6)
    let day = today;  // To prevent today value from changing

    // Function to toggle timetable visibility
    function timeTableAll() {
        document.getElementById('timetable').classList.toggle('active');
        setData(today);
        document.querySelector('.timetable div h2').innerHTML = "Today's Timetable";
    }

    // Event listeners for next and previous day buttons
    nextDay.onclick = function () {
        day <= 5 ? day++ : day = 0;  // If-else one-liner
        setData(day);
    }

    prevDay.onclick = function () {
        day >= 1 ? day-- : day = 6;  // Move to previous day
        setData(day);
    }

    // Set data on page load
    setData(day);
    document.querySelector('.timetable div h2').innerHTML = "Today's Timetable";  // Set heading on load
});

// Public file upload (uses public endpoint added to backend)
document.addEventListener('DOMContentLoaded', () => {
    const uploadBtn = document.getElementById('uploadBtn');
    const fileInput = document.getElementById('fileInput');
    const fileDesc = document.getElementById('fileDesc');
    const uploadResult = document.getElementById('uploadResult');

    if (!uploadBtn) return;

    uploadBtn.addEventListener('click', () => {
        uploadResult.innerText = '';
        if (!fileInput || fileInput.files.length === 0) {
            uploadResult.innerText = '请选择一个文件后再上传';
            return;
        }
        const fd = new FormData();
        fd.append('file', fileInput.files[0]);
        fd.append('description', fileDesc ? fileDesc.value : '');

        fetch('/files/upload', {
            method: 'POST',
            body: fd,
        })
            .then(res => res.json())
            .then(data => {
                if (!isSuccess(data)) {
                    uploadResult.innerText = data.message || JSON.stringify(data);
                } else if (data && data.data) {
                    uploadResult.innerText = '上传成功：' + (data.data.filename || JSON.stringify(data.data));
                    alert("上传成功")
                } else {
                    uploadResult.innerText = '上传成功';
                }
            })
            .catch(err => {
                uploadResult.innerText = '上传失败：' + err;
            });
    });
});


function getUserProfile() {
    const token = localStorage.getItem('authToken'); // ⭐ 正确读取你存的 token

    if (!token) {
        console.error("未找到 authToken, 请重新登录");
        return;
    }

    fetch("/api/v1/user/profile", {
        method: "GET",
        headers: {
            "Authorization": "Bearer " + token,
            "Accept": "application/json"
        }
    })
    .then(res => {
        if (!res.ok) throw new Error("HTTP error " + res.status);
        return res.json();
    })
    .then(data => {
        console.log("用户资料：", data);
        
        // ⭐⭐ 修复这里：正确的邮箱路径是 data.data.email
        if (data && data.data && data.data.email) {
            const email = data.data.email;
            const formattedEmail = email.replace('@', '<wbr>@');
            document.getElementById("displayed_email").innerHTML = formattedEmail;
            console.log("邮箱已设置：" + data.data.email);
        } else {
            console.error("无法获取邮箱，响应数据：", data);
            document.getElementById("displayed_email").innerHTML = "邮箱获取失败";
        }
    })
    .catch(err => {
        console.error("获取用户资料失败:", err);
    });
}
document.addEventListener("DOMContentLoaded", getUserProfile);

const API_BASE = '/api/v1'

function isSuccess(resp) {
    return resp && (resp.code === 0 || resp.code === 200)
}

function register() {
    const email = document.getElementById("email") && document.getElementById("email").value
    const username = document.getElementById("userid") && document.getElementById("userid").value
    const password = document.getElementById("password") && document.getElementById("password").value
    const confirmed_password = document.getElementById("confirm") && document.getElementById("confirm").value

    if (!email || !username || !password) {
        alert('请填写所有必填项');
        return;
    }
    if (password !== confirmed_password) {
        alert("两次密码输入不一致，请重新输入！");
        return;
    }

    const data = { email, username, password, confirmed_password }

    fetch(API_BASE + "/auth/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data)
    })
        .then(res => res.json())
        .then(result => {
            document.getElementById("regResult").innerText = JSON.stringify(result);
            if (isSuccess(result)) {
                // registration success
                window.location.href = "/register_result"
            } else {
                alert(result.message || '注册失败')
            }
        })
        .catch(err => {
            console.error(err)
            alert('注册请求失败')
        })

}

function login() {
    console.log('login() function called');

    // 从表单获取数据
    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const resultElement = document.getElementById("loginResult");

    console.log('Email:', email);
    console.log('Password:', password);

    if (!email && !password) {
        alert('请输入邮箱和密码');
        return false;
    } else if(!email) {
        alert('请输入邮箱');
        return false;
    } else if (!password) {
        alert('请输入密码');
        return false;
    }

    // 构建请求数据
    const requestData = {
        email: email,
        password: password
    };

    console.log('Sending request to:', API_BASE + "/auth/login");
    console.log('Request data:', requestData);

    // 发送POST请求
    fetch(API_BASE + "/auth/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json"
        },
        body: JSON.stringify(requestData),
        
    })
    
        .then(response => {
            console.log('Response status:', response.status);
            console.log('Response headers:', response.headers);

            // 首先检查HTTP状态码
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            // 尝试解析JSON
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                return response.json();
            } else {
                // 如果不是JSON，获取文本
                return response.text().then(text => {
                    console.log('Non-JSON response:', text);
                    throw new Error('Server returned non-JSON response');
                });
            }
        })
        .then(result => {
            console.log('Login response:', result);

            if (isSuccess(result)) {
                // 登录成功，存储token（如果有的话）
                if (result.data && result.data.token) {
                    localStorage.setItem('authToken', result.data.token);
                    localStorage.setItem('userEmail', email);
                    
                }
                // 跳转到首页
                console.log('Redirecting to /index');
                window.location.href = "/index";
            } else {
                // 显示错误信息
                const errorMsg = result.message || result.msg || '登录失败';
                resultElement.textContent = errorMsg;
                alert(errorMsg);
            }
        })
        .catch(err => {
            console.error('Login error:', err);

            // 更详细的错误信息
            let errorMsg = '登录请求失败';

            if (err.message.includes('Failed to fetch')) {
                errorMsg = '无法连接到服务器，请检查服务器是否运行';
            } else if (err.message.includes('non-JSON')) {
                errorMsg = '服务器返回了非JSON响应，请检查API';
            } else {
                errorMsg = err.message;
            }

            resultElement.textContent = errorMsg;
            alert(errorMsg);
        });

    return false; // 阻止表单提交
}

function clearLoginForm() {
    const emailInput = document.getElementById('email');
    const passwordInput = document.getElementById('password');
    
    if (emailInput) {
        emailInput.value = '';
        // 设置空值后再设置一次以确保清空
        emailInput.setAttribute('value', '');
    }
    
    if (passwordInput) {
        passwordInput.value = '';
        // 设置空值后再设置一次以确保清空
        passwordInput.setAttribute('value', '');
    }
    
    // 移除登出标志
    localStorage.removeItem('justLoggedOut');
}

// 页面加载完成后绑定事件
document.addEventListener('DOMContentLoaded', function () {
     // 检查是否是从登出跳转过来的
    const justLoggedOut = localStorage.getItem('justLoggedOut');
    if (justLoggedOut === 'true') {
        clearLoginForm();
    }
    
    // 额外的：如果用户已经登录，直接跳转到首页
    const logout_token = localStorage.getItem('authToken');
    if (logout_token && window.location.pathname === '/login') {
        window.location.href = '/index';
    }

    console.log('DOM loaded, initializing login form');

    const loginForm = document.getElementById('loginForm');
    if (loginForm) {
        console.log('Found login form');
        loginForm.addEventListener('submit', function (e) {
            console.log('Form submit event triggered');
            e.preventDefault(); // 阻止默认表单提交
            login();
        });
    } else {
        console.error('Login form not found! Check the HTML.');
    }

    // 为按钮添加点击事件作为备用
    const loginButton = document.querySelector('.btn-login');
    if (loginButton) {
        console.log('Found login button');
        loginButton.addEventListener('click', function (e) {
            console.log('Button click event triggered');
            e.preventDefault();
            login();
        });
    }

    // 检查用户是否已登录
    const token = localStorage.getItem('authToken');
    if (token) {
        console.log('User is already logged in with token');
        // 可以选择自动跳转或显示已登录状态
        // window.location.href = "/index";
    }
});
