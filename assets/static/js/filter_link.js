let page = 0;
page = document.getElementById('current-page')?.getAttribute('data-page');
page = Number(page);
function filterLink(pagennum) {
  let keyword = '';
  let baseurl = window.baseURL == '' ? '/' : window.baseURL;
  keyword = document.getElementById('keyword')?.value;
  let url;
  if (keyword == '') {
    url = baseurl + '?page=' + pagennum;
  } else {
    if (keyword.includes('#')) {
      url = baseurl + '?tag=' + keyword.replace('#', '') + '&page=' + pagennum;
    } else {
      url = baseurl + '?keyword=' + keyword + '&page=' + pagennum;
    }
  }
  console.log(url);
  window.location = url;
}

document
  .getElementById('link-filter-btn')
  ?.addEventListener('click', () => filterLink(page));
document
  .getElementById('next-page')
  ?.addEventListener('click', () => filterLink(page + 1));
document
  .getElementById('pre-page')
  ?.addEventListener('click', () => filterLink(page - 1));
