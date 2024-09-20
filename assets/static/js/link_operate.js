function linkOperate(evt) {
  console.log(evt.target);
  let operate = evt.target.getAttribute('data-operate');
  const root = window.baseURL;
  if (operate == null) return;
  else {
    evt.preventDefault();
    const id = evt.target.getAttribute('data-id');
    if (operate == 'edit') {
      window.location.href = root + `/link/edit/${id}`;
    } else if (operate == 'delete') {
      if (confirm('确定要删除吗?')) {
        fetch(root + `/api/link/delete/${id}`, {
          method: 'GET',
          cache: 'no-cache',
          credentials: 'include',
        }).then(res => {
          window.location.reload();
        });
      }
    } else if (operate == 'read') {
      fetch(root + `/api/link/read/${id}`, {
        method: 'GET',
        cache: 'no-cache',
        credentials: 'include',
      }).then(res => {
        evt.target.lastElementChild.innerText = '标为未读';
        evt.target.setAttribute('data-operate', 'unread');
      });
    } else if (operate == 'unread') {
      fetch(root + `/api/link/unread/${id}`, {
        method: 'GET',
        cache: 'no-cache',
        credentials: 'include',
      }).then(res => {
        evt.target.lastElementChild.innerText = '标为已读';
        evt.target.setAttribute('data-operate', 'read');
      });
    } else if (operate == 'view') {
      fetch(root + `/api/link/read/${id}`, {
        method: 'GET',
        cache: 'no-cache',
        credentials: 'include',
      }).then(res => {});
      const url = evt.target.getAttribute('data-url');
      window.open(url);
    }
  }
}
