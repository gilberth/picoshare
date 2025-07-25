/**
 * Theme Toggle - Sistema de alternancia de tema optimizado
 * Maneja el cambio entre modo claro y oscuro con mejor UX y accesibilidad
 */

document.addEventListener("DOMContentLoaded", () => {
  const toggleBtn = document.getElementById("theme-toggle-btn");
  if (!toggleBtn) return;

  // Estado del tema
  let isAnimating = false;
  
  /**
   * Detecta la preferencia inicial del usuario
   */
  function getInitialTheme() {
    const savedTheme = localStorage.getItem("theme");
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    
    if (savedTheme) {
      return savedTheme === "dark";
    }
    
    return prefersDark;
  }

  /**
   * Actualiza el icono del botón con animación suave
   */
  function updateToggleIcon(isDark, animate = true) {
    const icon = toggleBtn.querySelector("i");
    if (!icon) return;

    // Actualizar atributos de accesibilidad
    const ariaLabel = isDark ? "Cambiar a modo claro" : "Cambiar a modo oscuro";
    toggleBtn.title = ariaLabel;
    toggleBtn.setAttribute("aria-label", ariaLabel);
    toggleBtn.setAttribute("aria-pressed", isDark.toString());

    if (animate && !isAnimating) {
      // Animación suave del icono
      icon.style.transition = "transform 0.2s ease, opacity 0.15s ease";
      icon.style.transform = "rotate(90deg) scale(0.8)";
      icon.style.opacity = "0.5";
      
      setTimeout(() => {
        icon.className = isDark ? "fas fa-sun" : "fas fa-moon";
        icon.style.transform = "rotate(0deg) scale(1)";
        icon.style.opacity = "1";
      }, 150);
    } else {
      icon.className = isDark ? "fas fa-sun" : "fas fa-moon";
      icon.style.transform = "rotate(0deg) scale(1)";
      icon.style.opacity = "1";
    }
  }

  /**
   * Aplica el tema de forma eficiente
   */
  function applyTheme(isDark, withTransition = true) {
    // Aplicar/remover clase de modo oscuro
    document.body.classList.toggle("dark-mode", isDark);
    
    // Opcional: aplicar transición suave solo cuando se especifique
    if (withTransition) {
      document.body.style.transition = "background-color 0.3s ease, color 0.3s ease";
      setTimeout(() => {
        document.body.style.transition = "";
      }, 300);
    }
    
    // Actualizar meta theme-color para móviles
    updateThemeColor(isDark);
  }

  /**
   * Actualiza el color de tema para navegadores móviles
   */
  function updateThemeColor(isDark) {
    const themeColorMeta = document.querySelector('meta[name="theme-color"]');
    if (themeColorMeta) {
      themeColorMeta.content = isDark ? "#111827" : "#ffffff";
    }
  }

  /**
   * Maneja el clic del botón de alternancia
   */
  function handleToggle(e) {
    e.preventDefault();
    
    if (isAnimating) return;
    
    isAnimating = true;
    const isDark = !document.body.classList.contains("dark-mode");
    
    // Guardar preferencia
    localStorage.setItem("theme", isDark ? "dark" : "light");
    
    // Aplicar tema
    applyTheme(isDark, true);
    updateToggleIcon(isDark, true);
    
    // Animación del botón
    animateButton();
    
    // Disparar evento personalizado para otros componentes
    window.dispatchEvent(new CustomEvent('themechange', { 
      detail: { isDark } 
    }));
    
    setTimeout(() => {
      isAnimating = false;
    }, 300);
  }

  /**
   * Animación del botón con mejor feedback visual
   */
  function animateButton() {
    if (!toggleBtn.animate) return;
    
    const animation = toggleBtn.animate([
      { transform: "scale(1) rotate(0deg)" },
      { transform: "scale(0.95) rotate(-5deg)" },
      { transform: "scale(1.05) rotate(5deg)" },
      { transform: "scale(1) rotate(0deg)" }
    ], {
      duration: 250,
      easing: "cubic-bezier(0.34, 1.56, 0.64, 1)"
    });
    
    return animation;
  }

  /**
   * Maneja la navegación por teclado
   */
  function handleKeydown(e) {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      handleToggle(e);
    }
  }

  // Inicialización
  const initialTheme = getInitialTheme();
  applyTheme(initialTheme, false);
  updateToggleIcon(initialTheme, false);

  // Event listeners
  toggleBtn.addEventListener("click", handleToggle);
  toggleBtn.addEventListener("keydown", handleKeydown);

  // Escuchar cambios en la preferencia del sistema
  const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
  mediaQuery.addEventListener("change", (e) => {
    // Solo cambiar si no hay preferencia guardada
    if (!localStorage.getItem("theme")) {
      const isDark = e.matches;
      applyTheme(isDark, true);
      updateToggleIcon(isDark, true);
    }
  });

  // Mejorar estados hover/focus con CSS en lugar de JS
  toggleBtn.setAttribute("role", "switch");
  toggleBtn.setAttribute("aria-pressed", initialTheme.toString());
  
  // Añadir clase para identificar que JS está habilitado
  toggleBtn.classList.add("js-enabled");
});
