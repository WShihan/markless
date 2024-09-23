function setCookie(name, value) {
  var currentDate = new Date();
  var expiresDate = new Date(currentDate.getTime() + 7 * 24 * 60 * 60 * 1000); // 7天后的日期
  var expiresGMT = expiresDate.toGMTString();
  document.cookie = `${name}=${value}; expires=-1;path=/`;
  document.cookie = `${name}=${value}; expires=${expiresGMT};path=/;`;
}

function getCookie(name) {
  let mode = '';
  document.cookie.split(';').forEach(item => {
    if (item.split('=')[0].trim() == name) {
      mode = item.split('=')[1].trim();
    }
  });
  return mode;
}

function showMSG() {
  let msg = getCookie('message');
  let shown = getCookie('message_shown');
  console.log(msg, shown);
  if (msg == '' || shown == '' || shown == 'true') {
    return;
  } else {
    let popup = document.createElement('div');
    popup.className = 'popup';
    popup.id = 'popup';
    msg = decodeURIComponent(msg);
    popup.innerHTML = `<div class="popup-inner"><div class="msg">${msg}</div><div class="close" onclick="cleatTip()">x</div></div>`;
    document.body.appendChild(popup);
    setCookie('message_shown', 'true');
    setTimeout(() => {
      cleatTip();
    }, 5000);
  }
}
showMSG();

function cleatTip() {
  let popup = document.getElementById('popup');
  document.body.removeChild(popup);
}
