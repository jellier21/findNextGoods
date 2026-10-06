
const account = localStorage.getItem('userAccountNum');

function convertImageToJPG(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = (e) => {
            const img = new Image();
            img.onload = () => {
                const canvas = document.createElement('canvas');
                const ctx = canvas.getContext('2d');

                // 设置 canvas 尺寸（保持原图比例或指定最大尺寸）
                canvas.width = img.width;
                canvas.height = img.height;

                // 填充白色背景（避免透明 PNG 转 JPG 后变黑）
                ctx.fillStyle = '#ffffff';
                ctx.fillRect(0, 0, canvas.width, canvas.height);

                // 绘制原图
                ctx.drawImage(img, 0, 0);

                // 转为 JPG Data URL（质量 0.92）
                canvas.toBlob(blob => {
                    if (blob) {
                        resolve(blob);
                    } else {
                        reject(new Error('Canvas toBlob failed'));
                    }
                }, 'image/jpeg', 0.92);
            };
            img.onerror = reject;
            img.src = e.target.result;
        };
        reader.onerror = reject;
        reader.readAsDataURL(file);
    });
}

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

    let showHeader = () => {
        let myImageCode = account + ".jpg";
        let imgElement = document.getElementById('header-image');
        let imageUrl = `http://localhost:8080/api/image/header/show?filename=${encodeURIComponent(myImageCode)}`;
        // console.log('尝试加载图片:', imageUrl);

        // 添加错误监听
        imgElement.onerror = function () {
            console.error('本登录者图片加载失败！');
            // alert('图片加载失败！请检查：\n1. 后端服务是否启动\n2. 图片路径是否正确\n3. 浏览器控制台错误信息');
        };

        imgElement.onload = function () {
            console.log('图片加载成功！');
        };

        // 直接设置 img 的 src
        imgElement.src = imageUrl;
    };
    showHeader();
}



function showTheHomeGoods() {
    let xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/api/goods/allgoods', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = () => {
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    var jsonData = JSON.parse(xhr.responseText);
                    console.log("Response data:", jsonData);

                    if (Array.isArray(jsonData) && jsonData.length > 0) {
                        document.getElementById("showTheGoods").innerHTML = "";
                        console.log("ssssssssssssssssssssssssss:", jsonData.length);
                        function showGoodImage(num) {
                            let personImageCode = num;
                            let imageUrl = `http://localhost:8080/api/image/goods/show?filename=${encodeURIComponent(personImageCode)}`;
                            // console.log('尝试加载图片:', imageUrl);
                            return imageUrl;
                        }
                        for (let item of jsonData) {
                            try {
                                console.log("正在处理数据:", item.name);
                                document.getElementById("showTheGoods").innerHTML += `
                                    <div id="simple-obj" onclick="window.location.href='./exhibition.html?goodsId=${item.id}'">
                                        <div id="obj-image">
                                            <img src="${showGoodImage(item.image)}.jpg" alt="no image">
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

    

    /*个人信息展示*/
    showTheAccount();
    /*home展示*/
    showTheHomeGoods();
});    