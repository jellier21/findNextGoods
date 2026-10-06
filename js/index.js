function LoginByApi() { 
    let username = document.getElementById('username_l').value;
    let password = document.getElementById('password_l').value;
    let data = {
        accDB: username,
        passkey: password
    }
    let xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/api/user/login', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = () => { 
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    var jsonData = JSON.parse(xhr.responseText);
                    console.log("登录响应数据：", jsonData);

                    if (jsonData.Account) {
                        // alert("登录成功！");
                        alert(jsonData.Nickname)
                        localStorage.setItem('userAccountNum', jsonData.Account);
                        localStorage.setItem('userPassword', jsonData.Passkey);
                        location.replace('./homepage.html');
                    } else {
                        alert("登录失败：" + xhr.responseText);
                    }
                } catch (e) {
                    console.error("JSON 解析失败:", e);
                    alert("服务器返回数据格式错误：" + xhr.responseText);
                }
            } else {
                console.log('请求失败，状态码:', xhr.status);
                alert("登录失败，状态码：" + xhr.status);
            }
        }
    };
    xhr.send(JSON.stringify(data));
}

function AddUserByApi() { 
    let nickname = document.getElementById('nickname_r').value;
    let account = document.getElementById('account_r').value;
    let password = document.getElementById('password_r').value;
    let data = {
        Nickname: nickname,
        Account: account,
        Passkey: password
    }
    let xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/api/user/add', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = () => { 
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                alert("注册成功！");
                location.reload();
            } else {
                alert("注册失败,该账号已存在，请重新选择账号！");
                log.console.log('请求失败，状态码:', xhr.status);
            }
        }
    };
    xhr.send(JSON.stringify(data));
}


document.addEventListener('DOMContentLoaded', ()=>{
    document.getElementById('loginBtn').addEventListener('click', LoginByApi);//登录按钮事件
    document.getElementById('registerBtn').addEventListener('click', AddUserByApi);//注册按钮事件

    //切换注册<->登录
    document.getElementById('switchToRegister').addEventListener('click', ()=>{
        document.getElementById('login').style.display = 'none';
        document.getElementById('register').style.display = 'block';
        document.getElementById('switchToLogin').style.display = 'block';
        document.getElementById('switchToRegister').style.display = 'none';
    });
    
    document.getElementById('switchToLogin').addEventListener('click', ()=>{
        document.getElementById('login').style.display = 'block';
        document.getElementById('register').style.display = 'none';
        document.getElementById('switchToLogin').style.display = 'none';
        document.getElementById('switchToRegister').style.display = 'block';
    });

    //注册时两次密码是否一致
    document.getElementById('confirm_password_r').addEventListener('input', ()=>{
        let password = document.getElementById('password_r').value;
        let confirmPassword = document.getElementById('confirm_password_r').value;
        if (password !== confirmPassword) {
            document.getElementById('confirm_password_r').style.borderColor = 'red';
            document.getElementById('submitBtn').style.backgroundColor = 'gray';
            document.getElementById('refreshBtn').style.backgroundColor = 'gray';
            document.getElementById('submitBtn').style.pointerEvents = 'none';
            document.getElementById('refreshBtn').style.pointerEvents = 'none';
        } else {
            document.getElementById('confirm_password_r').style.borderColor = 'green';
            document.getElementById('submitBtn').style.backgroundColor = '#007bff';
            document.getElementById('refreshBtn').style.backgroundColor = '#007bff';
            document.getElementById('submitBtn').style.pointerEvents = 'auto';
            document.getElementById('refreshBtn').style.pointerEvents = 'auto';
        }
    });

    
});