// 搜索
function searchLink() {
  let searchBtn = document.getElementById('link-search');
  if (searchBtn instanceof HTMLButtonElement) {
    searchBtn.addEventListener('click', () => {
      let keyword = '';
      keyword = document.getElementById('keyword')?.value;
      if (keyword == '') return;
      window.location.href = `${window.baseURL}/?keyword=${keyword}`;
    });
  }
}
searchLink();
