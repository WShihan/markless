const modeCookieKey = 'markless-view-mode';
const NORMALMODE = 'normal';
const DARKMODE = 'dark';

function setCookie(name, value) {
  var currentDate = new Date();
  var expiresDate = new Date(currentDate.getTime() + 7 * 24 * 60 * 60 * 1000); // 7天后的日期
  var expiresGMT = expiresDate.toGMTString();
  document.cookie = `${name}=${value}; expires=-1;path=/`;
  document.cookie = `${name}=${value}; expires=${expiresGMT};path=/;`;
}

function getCookie(name) {
  let mode = NORMALMODE;
  document.cookie.split(';').forEach(item => {
    if (item.split('=')[0].trim() == name) {
      mode = item.split('=')[1].trim();
    }
  });
  return mode;
}

function switchMode(mode) {
  if (mode == NORMALMODE) {
    document.body.classList.add(NORMALMODE);
    document.body.classList.remove(DARKMODE);
  } else {
    document.body.classList.add(DARKMODE);
    document.body.classList.remove(NORMALMODE);
  }
}
// 设置默认显示模式
switchMode(getCookie(modeCookieKey));

document.addEventListener('DOMContentLoaded', () => {
  document.getElementById('theme-switch')?.addEventListener('click', function (evt) {
    let mode = getCookie(modeCookieKey);
    mode = mode == NORMALMODE ? DARKMODE : NORMALMODE;
    switchMode(mode);
    setCookie(modeCookieKey, mode);
    console.log(mode);
  });
});
