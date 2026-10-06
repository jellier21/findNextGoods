const account = localStorage.getItem('userAccountNum');
function NavChange(){
//点击导航条转换
  let NavIssue = document.getElementById("item-issue");
  let NavReward = document.getElementById("item-reward");
  let NavContact = document.getElementById("item-contact");

  let Issue = document.getElementById("show-issue");
  let Reward = document.getElementById("show-reward");
  let Contact = document.getElementById("show-contact");

  NavIssue.onclick = () => {
    Issue.style.display = "block";
    Reward.style.display = "none";
    Contact.style.display = "none";
  }
  NavReward.onclick = () => {
    Issue.style.display = "none";
    Reward.style.display = "block";
    Contact.style.display = "none";
  }
  NavContact.onclick = () => {
    Issue.style.display = "none";
    Reward.style.display = "none";
    Contact.style.display = "block";
  }

}

function ShowOneContact(){ 
  // 显示联系方式
  let relWay = document.getElementById("software");
  let AccNum = document.getElementById("account-Num");
    let xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/api/site/mysite', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onreadystatechange = () => { 
        if (xhr.readyState === 4) {
            if (xhr.status === 200) {
                try {
                    var jsonData = JSON.parse(xhr.responseText);
                    console.log("联系平台响应数据：", jsonData); 
                    relWay.innerHTML = jsonData.SiteWay;
                    AccNum.innerHTML = jsonData.Location;
                } catch (e) {
                    console.error("JSON 解析失败:", e);
                    document.getElementById("software").innerHTML = xhr.responseText;
                    document.getElementById("account-Num").innerHTML = xhr.responseText;
                }
            } else {
                // alert("用户名返回失败：" + xhr.status);
                document.getElementById("software").innerHTML = "未设置";
                document.getElementById("account-Num").innerHTML = "未设置";
            }
        }
    };
    
    xhr.send(JSON.stringify({ Account: account }));
}

function changeContact(){ 
  // 修改联系方式
  document.getElementById("change-account").addEventListener("click", () => {
    document.getElementById("changeArea").style.display = "block";
    document.getElementById("confirmChange").addEventListener("click", () => {
      let newSoftware = document.getElementById("newSoftware").value;
      let newAccountNum = document.getElementById("newAccountNum").value;
      let data = {
        "account": account,
        "rel_way": newSoftware,
        "site": newAccountNum
      }
      let xhr = new XMLHttpRequest();
      xhr.open('POST', 'http://localhost:8080/api/site/upload', true);
      xhr.setRequestHeader('Content-Type', 'application/json');
      xhr.onreadystatechange = () => { 
          if (xhr.readyState === 4) {
              if (xhr.status === 200) {
                  try {
                      var jsonData = JSON.parse(xhr.responseText);
                      console.log("服务器响应数据：", jsonData); 
                      if (jsonData.bool === true) {
                          alert("修改成功！");
                          window.location.reload();
                      }
                  } catch (e) {
                      console.error("JSON 解析失败:", e);
                  }
              } else {
                  // alert("用户名返回失败：" + xhr.status);
                  console.error("请求失败，状态码:", xhr.status);
              }
          }
      };
      
      xhr.send(JSON.stringify(data));
    })
  });
  document.getElementById("cancelChange").addEventListener("click", () => {
    document.getElementById("changeArea").style.display = "none";
  });
}


function changeTypeNumToString(typeNum) { 
  if (typeNum === "1") {
    return "自制";
  } else {
    return "转卖";
  }
}


function deleteByid(id) {
  if (!id) {
    console.warn("无效的删除 ID");
    return;
  }
  const account = localStorage.getItem('userAccountNum');
  if (!account) {
    alert("请先登录！");
    return;
  }

  if (!confirm("确定要删除吗？")) return;

  let xh = new XMLHttpRequest();
  xh.open('POST', 'http://localhost:8080/api/selfmade/delete', true);
  xh.setRequestHeader('Content-Type', 'application/json');
  xh.onreadystatechange = () => {
    if (xh.readyState === 4) {
      if (xh.status === 200) {
        try {
          let res = JSON.parse(xh.responseText);
          if (res.success) {
            alert("删除成功！");
            window.location.reload();
          } else {
            alert("删除失败：" + (res.message || ""));
          }
        } catch (e) {
          console.error("响应解析失败:", e);
          console.log("响应内容:", xh.responseText);
          alert("服务器返回格式错误");
        }
      } else {
        console.error("请求失败，状态码:", xh.status);
        console.log("响应内容:", xh.responseText);
        alert("网络请求失败，请重试");
      }
    }
  };
  
  // 确保 id 是数字类型
  let deleteId = parseInt(id);
  if (isNaN(deleteId)) {
    alert("无效的 ID 格式");
    return;
  }
  
  xh.send(JSON.stringify({ id: deleteId }));
}


function ShowMyIssue(){ 
  // 显示我的发布
      let xhr_s = new XMLHttpRequest();
      xhr_s.open('POST', 'http://localhost:8080/api/selfmade/showmyself', true);
      xhr_s.setRequestHeader('Content-Type', 'application/json');
      xhr_s.onreadystatechange = () => { 
          if (xhr_s.readyState === 4) {
              if (xhr_s.status === 200) {
                  try {
                      let jsonData = JSON.parse(xhr_s.responseText);
                      document.getElementById("issue-container").innerHTML = ""; // 清空容器
                      if (Array.isArray(jsonData) && jsonData.length > 0) {
                        jsonData.forEach(item => {
                          if (item.type === 1) {
                            document.getElementById("issue-container").innerHTML += `
                              <div id="my-one-issue">
                                  <div id="issue-image">
                                      <img src="./images/universal/worring.png" alt="问题图片">
                                  </div>
                                  <div id="issue-info">
                                      <div id="issue-title">${item.name}</div>
                                      <div id="issue-description">${item.introduction}</div>
                                  </div>
                                  <div id="issue-operation">
                                      <div id="issue-type">自制</div>
                                      <div id="issue-price">${"$" + item.price}</div>
                                      <div id="issue-delete" onclick="deleteByid('${item.id}')">删除</div>
                                  </div>
                              </div>                    
                            ` 
                          }
                        });
                    } else {
                      console.log("没有自制品数据");
                    }
                      
                  } catch (e) {
                      console.error("JSON 解析失败:", e);
                  }
              } else {
                  console.error("请求失败，状态码:", xhr_s.status);
              }
          }
      };
      
      xhr_s.send(JSON.stringify({ Account: account }));

      let xhr_g = new XMLHttpRequest();
      xhr_g.open('POST', 'http://localhost:8080/api/goods/showmygoods', true);
      xhr_g.setRequestHeader('Content-Type', 'application/json');
      xhr_g.onreadystatechange = () => { 
          if (xhr_g.readyState === 4) {
              if (xhr_g.status === 200) {
                  try {
                      let jsonData_goods = JSON.parse(xhr_g.responseText);
                      // document.getElementById("issue-container").innerHTML = "";
                      if (Array.isArray(jsonData_goods) && jsonData_goods.length > 0) {
                        jsonData_goods.forEach(item => {
                          if (item.type === 2) {
                            document.getElementById("issue-container").innerHTML += `
                              <div class="myGoods">
                                <div class="showArea_g">
                                  <div class="good-image-area">
                                    <img src="./images/universal/no-img.png" alt="no img" class="active">
                                    <img src="./images/universal/worringImg.png" alt="">
                                    <img src="./images/universal/test.png" alt="无图片">
                                    <img src="./images/universal/myheader.jpg" alt="无图片">
                                  </div>
                                  <div class="changeimg">
                                    <div class="left">←</div>
                                    <div class="right">→</div>
                                  </div>
                                </div>
                                <div class="good-info-area">
                                  <div class="good-name">${item.name}</div>
                                  <div class="good-description">${item.description}</div>
                                </div>
                                <div class="good-info_operation">
                                  <div class="good-type">二手</div>
                                  <div class="good-price">$${item.price}</div>
                                  <div class="good-delete" onclick="deleteByid('${item.id}')">删除</div>
                                </div>
                              </div>
                            `;
                          }
                        });
                    } else {
                        console.log("无二手商品数据");
                    }
                      
                  } catch (e) {
                      console.error("JSON 解析失败:", e);
                  }
              } else {
                  // alert("用户名返回失败：" + xhr_g.status);
                  console.error("请求失败，状态码:", xhr_g.status);
              }
          }
      };
      
      xhr_g.send(JSON.stringify({ Account: account }));
}


function deleteNeedByid(id) {
  if (!id) {
    console.warn("无效的删除 ID");
    return;
  }
  if (!account) {
    alert("请先登录！");
    return;
  }

  if (!confirm("确定要删除吗？")) return;

  let xh = new XMLHttpRequest();
  xh.open('POST', 'http://localhost:8080/api/need/delete', true);
  xh.setRequestHeader('Content-Type', 'application/json');
  xh.onreadystatechange = () => {
    if (xh.readyState === 4) {
      if (xh.status === 200) {
        try {
          let res = JSON.parse(xh.responseText);
          if (res.success) {
            alert("删除成功！");
            window.location.reload();
          } else {
            alert("删除失败：" + (res.message || ""));
          }
        } catch (e) {
          console.error("响应解析失败:", e);
          console.log("响应内容:", xh.responseText);
          alert("服务器返回格式错误");
        }
      } else {
        console.error("请求失败，状态码:", xh.status);
        console.log("响应内容:", xh.responseText);
        alert("网络请求失败，请重试");
      }
    }
  };
  
  // 确保 id 是数字类型
  let deleteId = parseInt(id);
  if (isNaN(deleteId)) {
    alert("无效的 ID 格式");
    return;
  }
  
  xh.send(JSON.stringify({ id: deleteId }));
}
function ShowMyNeedRequire(){ 
  // 显示我的发布
      let xhr = new XMLHttpRequest();
      xhr.open('POST', 'http://localhost:8080/api/need/showmy', true);
      xhr.setRequestHeader('Content-Type', 'application/json');
      xhr.onreadystatechange = () => { 
          if (xhr.readyState === 4) {
              if (xhr.status === 200) {
                  try {
                      let jsonData = JSON.parse(xhr.responseText);
                      console.log("jsonData", jsonData);
                      document.getElementById("reward-container").innerHTML = ""; // 清空容器
                      if (Array.isArray(jsonData) && jsonData.length > 0) {
                        jsonData.forEach(item => {
                          console.log("type", item.type);
                          document.getElementById("reward-container").innerHTML += `
                            <div id="my-one-reward">
                                <div id="reward-image">
                                    <img src="./images/universal/worring.png" alt="悬赏图片">
                                </div>
                                <div id="reward-info">
                                    <div id="reward-description">${item.description}</div>
                                    <div id="reward-delete" onclick="deleteNeedByid('${item.id}')">删除</div>
                                </div>
                            </div>                   
                          `
                           // onclick="deleteIssue('${item.id}')"
                        });
                    } else {
                        document.getElementById("reward-container").innerHTML = "无数据";
                    }
                      
                  } catch (e) {
                      console.error("JSON 解析失败:", e);
                  }
              } else {
                  console.error("请求失败，状态码:", xhr.status);
              }
          }
      };
      
      xhr.send(JSON.stringify({ Account: account }));
}



document.addEventListener("DOMContentLoaded", () => {

  console.log("用户名：", account);
  NavChange();
  ShowOneContact();
  changeContact();
  ShowMyIssue();
  ShowMyNeedRequire();
  

  // 使用事件委托监听整个 issue-container 的点击
  document.getElementById("issue-container").addEventListener("click", function(e) {
    if (e.target.classList.contains("left") || e.target.classList.contains("right")) {
      // 找到最近的 .good-image-area
      const imageArea = e.target.closest(".myGoods").querySelector(".good-image-area");
      const images = imageArea.querySelectorAll("img");
      
      // 找出当前 active 的图片索引
      let currentIndex = -1;
      images.forEach((img, index) => {
        if (img.classList.contains("active")) {
          currentIndex = index;
        }
      });

      // 如果没有 active，默认第一个
      if (currentIndex === -1) {
        currentIndex = 0;
        images[0].classList.add("active");
      }

      // 计算新索引
      let newIndex;
      if (e.target.classList.contains("left")) {
        newIndex = (currentIndex - 1 + images.length) % images.length;
      } else {
        newIndex = (currentIndex + 1) % images.length;
      }

      // 移除所有 active，设置新的
      images.forEach(img => img.classList.remove("active"));
      images[newIndex].classList.add("active");
    }
  });
});