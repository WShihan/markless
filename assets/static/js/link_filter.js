let page = 0;
page = document.getElementById('current-page')?.getAttribute('data-page');
page = Number(page);
function filterLink(pagennum) {
  let keyword = '';
  let root = window.location.href.split('?')[0];
  keyword = document.getElementById('keyword')?.value;
  let url;
  if (keyword == '') {
    url = root + '?page=' + pagennum;
  } else {
    if (keyword.includes('#')) {
      url = root + '?tag=' + keyword.replace('#', '') + '&page=' + pagennum;
    } else {
      url = root + '?keyword=' + keyword + '&page=' + pagennum;
    }
  }
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

  document.getElementById('keyword')?.addEventListener('keyup', function (event) {
    if (event.key == 'Enter') {
      filterLink(page);
    }
  })