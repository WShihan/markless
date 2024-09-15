function linkHandler(evt) {
  evt.preventDefault();
  let id = evt.target.getAttribute('data-id');
  let operation = evt.target.getAttribute('data-operation');
  if (id == null) {
    id = evt.target.parentElement.getAttribute('data-id');
    operation = evt.target.parentElement.getAttribute('data-operation');
  }
  fetch('/link/' + operation, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      id: id,
    }),
  });
}
