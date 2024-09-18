document.addEventListener('DOMContentLoaded', function () {
  const imgs = document.querySelectorAll('.img-lazy');
  const lazyLoad = (entries, observer) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const img = entry.target;
        img.src = img.getAttribute('data-src');
        img.onload = () => {
          img.style.opacity = 1;
          img.parentElement.style.backgroundImage = 'none';
        };
        img.onerror = () => {};
        observer.unobserve(img);
      }
    });
  };

  const observer = new IntersectionObserver(lazyLoad, {
    root: null,
    rootMargin: '0px',
    threshold: 0.1,
  });

  imgs.forEach(img => {
    observer.observe(img);
  });
});
