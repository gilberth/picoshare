// Botón de alternancia de tema oscuro/claro mejorado
// Agrega/remueve la clase 'dark-mode' en <body> y guarda preferencia en localStorage

document.addEventListener("DOMContentLoaded", () => {
  const toggleBtn = document.getElementById("theme-toggle-btn");
  if (!toggleBtn) return;

  // Función para actualizar el icono del botón con animación
  function updateToggleIcon(isDark, animate = true) {
    const icon = toggleBtn.querySelector("i");
    if (icon) {
      if (animate) {
        // Animación de rotación al cambiar icono
        icon.style.transform = "rotate(180deg) scale(0.8)";
        icon.style.opacity = "0.6";
        
        setTimeout(() => {
          icon.className = isDark ? "fas fa-sun" : "fas fa-moon";
          icon.style.transform = "rotate(0deg) scale(1)";
          icon.style.opacity = "1";
        }, 150);
      } else {
        icon.className = isDark ? "fas fa-sun" : "fas fa-moon";
      }
    }
    toggleBtn.title = isDark ? "Cambiar a modo claro" : "Cambiar a modo oscuro";
    toggleBtn.setAttribute("aria-label", isDark ? "Cambiar a modo claro" : "Cambiar a modo oscuro");
  }

  // Función para aplicar tema con transición suave
  function applyTheme(isDark, withTransition = true) {
    if (withTransition) {
      // Agregar clase de transición temporal
      document.body.style.setProperty('--transition-duration', '0.4s');
      
      // Aplicar el tema
      document.body.classList.toggle("dark-mode", isDark);
      
      // Forzar estilos del dropdown si es modo oscuro
      if (isDark) {
        setTimeout(() => {
          forceDropdownStyles();
        }, 100);
      }
      
      // Remover la transición después de completarse
      setTimeout(() => {
        document.body.style.removeProperty('--transition-duration');
      }, 400);
    } else {
      document.body.classList.toggle("dark-mode", isDark);
      if (isDark) {
        setTimeout(() => {
          forceDropdownStyles();
        }, 100);
      }
    }
  }

  // Función para forzar estilos del dropdown usando clases CSS (CSP compatible)
  function forceDropdownStyles() {
    const dropdowns = document.querySelectorAll('.navbar-dropdown');
    dropdowns.forEach(dropdown => {
      if (document.body.classList.contains('dark-mode')) {
        // Solo usar clases CSS para compatibilidad con CSP
        dropdown.classList.add('forced-dark-dropdown');
        
        // Forzar re-aplicación de estilos CSS
        dropdown.offsetHeight; // trigger reflow
        
        // Asegurar que el dropdown específico tenga el ID correcto
        if (dropdown.parentElement && dropdown.parentElement.querySelector('.navbar-link')) {
          const linkText = dropdown.parentElement.querySelector('.navbar-link').textContent.trim();
          if (linkText === 'System') {
            dropdown.id = 'system-dropdown';
          }
        }
      } else {
        // Modo claro - remover clase forzada
        dropdown.classList.remove('forced-dark-dropdown');
      }
    });
  }

  // Función simplificada que solo manipula clases
  function continuouslyForceStyles() {
    forceDropdownStyles();
    
    // Forzar re-renderizado
    const systemDropdown = document.getElementById('system-dropdown');
    if (systemDropdown && document.body.classList.contains('dark-mode')) {
      systemDropdown.classList.add('forced-dark-dropdown');
      // Trigger reflow para forzar aplicación de estilos
      systemDropdown.offsetHeight;
    }
  }

  // Inicializar según preferencia guardada o sistema
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const savedTheme = localStorage.getItem("theme");
  const shouldUseDark = savedTheme === "dark" || (!savedTheme && prefersDark);
  
  // Aplicar tema inicial sin animación
  applyTheme(shouldUseDark, false);
  updateToggleIcon(shouldUseDark, false);

  // Escuchar cambios en la preferencia del sistema
  window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", (e) => {
    if (!localStorage.getItem("theme")) {
      const isDark = e.matches;
      applyTheme(isDark);
      updateToggleIcon(isDark);
    }
  });

  // Manejar el click del botón con mejores animaciones
  toggleBtn.addEventListener("click", (e) => {
    e.preventDefault();
    
    // Prevenir clicks múltiples rápidos
    if (toggleBtn.classList.contains("animating")) return;
    toggleBtn.classList.add("animating");
    
    const isDark = !document.body.classList.contains("dark-mode");
    
    // Guardar preferencia
    localStorage.setItem("theme", isDark ? "dark" : "light");
    
    // Aplicar tema con transición
    applyTheme(isDark);
    updateToggleIcon(isDark);
    
    // Animación del botón con rebote
    const keyframes = [
      { transform: "scale(1) rotate(0deg)" },
      { transform: "scale(0.9) rotate(-10deg)" },
      { transform: "scale(1.1) rotate(10deg)" },
      { transform: "scale(1) rotate(0deg)" }
    ];
    
    const animation = toggleBtn.animate(keyframes, {
      duration: 300,
      easing: "cubic-bezier(0.68, -0.55, 0.265, 1.55)"
    });
    
    animation.onfinish = () => {
      toggleBtn.classList.remove("animating");
    };
  });

  // Mejores estados de hover y focus
  toggleBtn.addEventListener("mouseenter", () => {
    if (!toggleBtn.classList.contains("animating")) {
      toggleBtn.style.transform = "translateY(-1px) rotate(15deg)";
    }
  });

  toggleBtn.addEventListener("mouseleave", () => {
    if (!toggleBtn.classList.contains("animating")) {
      toggleBtn.style.transform = "translateY(0) rotate(0deg)";
    }
  });

  // Accesibilidad: soporte para tecla Enter y Espacio
  toggleBtn.addEventListener("keydown", (e) => {
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault();
      toggleBtn.click();
    }
  });

  // Observer para detectar cuando se crean dropdowns y aplicar estilos
  const observer = new MutationObserver((mutations) => {
    mutations.forEach((mutation) => {
      if (mutation.type === 'childList') {
        mutation.addedNodes.forEach((node) => {
          if (node.nodeType === 1) { // Element node
            if (node.classList && node.classList.contains('navbar-dropdown')) {
              if (document.body.classList.contains('dark-mode')) {
                setTimeout(() => forceDropdownStyles(), 0);
              }
            } else if (node.querySelector && node.querySelector('.navbar-dropdown')) {
              if (document.body.classList.contains('dark-mode')) {
                setTimeout(() => forceDropdownStyles(), 0);
              }
            }
          }
        });
      }
    });
  });

  // Observar cambios en el DOM
  observer.observe(document.body, {
    childList: true,
    subtree: true
  });

  // Aplicar estilos iniciales si ya hay dropdowns
  if (document.body.classList.contains('dark-mode')) {
    setTimeout(() => forceDropdownStyles(), 100);
  }

  // SOLUCION AGRESIVA: Forzar estilos cada 100ms para asegurar que se apliquen
  setInterval(() => {
    if (document.body.classList.contains('dark-mode')) {
      continuouslyForceStyles();
    }
  }, 100);

  // Ejecutar inmediatamente al cargar
  continuouslyForceStyles();
});
