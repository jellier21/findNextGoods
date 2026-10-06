
const account = localStorage.getItem('userAccountNum');

function findMessageByAccount(account_s) {
    return new Promise((resolve, reject) => {
        let hr = new XMLHttpRequest();
        hr.open('POST', 'http://localhost:8080/api/user/account-nick', true);
        hr.setRequestHeader('Content-Type', 'application/json');
        hr.onreadystatechange = () => {
            if (hr.readyState === 4) {
                if (hr.status === 200) {
                    try {
                        let jsonData = JSON.parse(hr.responseText);
                        resolve(jsonData.nick || null);
                    } catch (e) {
                        console.error("JSON 解析失败:", e);
                        reject(e);
                    }
                } else {
                    console.error("请求失败，状态码:", hr.status);
                    reject(new Error(`HTTP ${hr.status}`));
                }
            }
        };
        hr.send(JSON.stringify({ account: account_s }));
    });
}

function showTheAccount() {
    let showNick = document.getElementById("myName");

    let xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/api/user/account-nick', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = () => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    let jsonData = JSON.parse(xhr.responseText);
                    console.log("昵称响应数据：", jsonData);
                    showNick.innerHTML = jsonData.nickname;
                } catch (e) {
                    console.error("JSON 解析失败:", e);
                    showNick.innerHTML = xhr.responseText;
                }
            } else {
                // alert("用户名返回失败：" + xhr.status);
                console.error("请求失败，状态码:", xhr.status);
            }
        }
    };

    xhr.send(JSON.stringify({ Account: account }));
}



function showTheHomeGoods() {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/api/goods/allgoods', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = ()=> {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    var jsonData = JSON.parse(xhr.responseText); 
                    console.log("Response data:", jsonData);

                    if (Array.isArray(jsonData) && jsonData.length > 0) {
                        document.getElementById("showTheGoods").innerHTML = "";
                        for (let item of jsonData) {
                            try {
                                document.getElementById("showTheGoods").innerHTML += `
                                    <div id="simple-obj" onclick="window.location.href='./exhibition.html?goodsId=${item.id}'">
                                        <div id="obj-image">
                                            <img src="./images/universal/no-img.png" alt="no image">
                                        </div>
                                        <div id="obj-name">${item.name}</div>
                                    </div>
                                `
                            } catch (e) {
                                console.error('获取失败:', e);
                            }
                        }
                    } else {
                        document.getElementById("showTheGoods").innerHTML = "无数据";
                    }
                } catch (e) {
                    console.log('JSON 解析失败:', e);
                }
            } else {
                console.log('请求失败，状态码:', xhr.status);
            }
        }
    };

    xhr.send();
}



document.addEventListener("DOMContentLoaded", () => {
    /*是否登录*/
    document.getElementById("myPortrait").addEventListener("click", () => {
        if (localStorage.getItem('userAccountNum') == null) {
            window.location.href = "./index.html";
        } else {
            window.location.href = "./myInterface.html";
        }
    })




    /*切换导航页面*/
    document.getElementById("home").addEventListener("click", () => {
        window.location.href = "./homepage.html";
    });
    document.getElementById("post").addEventListener("click", () => {
        window.location.href = "./find_Needs.html";
    });
    document.getElementById("self-made").addEventListener("click", () => {
        window.location.href = "./Self-made.html";
    });

    /*鼠标 移入+移出 个人信息显示*/
    document.getElementById("myPortrait").addEventListener("mouseover", () => {
        document.getElementById("HoverItem").style.display = "block";
    });
    document.getElementById("myPortrait").addEventListener("mouseout", () => {
        document.getElementById("HoverItem").style.display = "none";
    });
    document.getElementById("HoverItem").addEventListener("mouseover", () => {
        document.getElementById("HoverItem").style.display = "block";
    });
    document.getElementById("HoverItem").addEventListener("mouseout", () => {
        document.getElementById("HoverItem").style.display = "none";
    });

    /*个人信息展示*/
    showTheAccount();
    /*home展示*/
    showTheHomeGoods();
});    