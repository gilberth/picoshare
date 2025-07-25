// thumbnail-loader.js
class ThumbnailLoader {
  constructor() {
    this.observer = new IntersectionObserver(
      this.handleIntersection.bind(this),
      { rootMargin: '50px' }
    );
    this.init();
  }

  init() {
    document.querySelectorAll('.file-thumbnail img[loading="lazy"]')
      .forEach(img => this.observer.observe(img));
  }

  handleIntersection(entries) {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        this.loadThumbnail(entry.target);
        this.observer.unobserve(entry.target);
      }
    });
  }

  loadThumbnail(img) {
    const container = img.closest('.file-thumbnail');
    container.classList.add('loading');
    
    img.onload = () => {
      container.classList.remove('loading');
      container.setAttribute('data-has-preview', 'true');
    };
    
    img.onerror = () => {
      container.classList.remove('loading');
      container.classList.add('fallback');
    };
  }
}

// Inicializar cuando el DOM esté listo
document.addEventListener('DOMContentLoaded', () => {
  new ThumbnailLoader();
});

export { ThumbnailLoader };