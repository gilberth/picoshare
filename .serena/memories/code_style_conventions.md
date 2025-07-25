# PicoShare - Code Style and Conventions

## Go Code Style
- Standard Go formatting with `gofmt`
- Package-level error handling with detailed logging
- HTTP handlers return `http.HandlerFunc`
- Use gorilla/mux for routing
- Embed templates and static files using `embed.FS`
- Context-based authentication checks

## Frontend Code Style
- **JavaScript**: ES6+ modules, no bundler
- **CSS**: CSS custom properties (variables) for theming
- **HTML**: Go templates with semantic markup
- **Icons**: FontAwesome classes with CSS variables for colors

## File Organization
- Templates in `handlers/templates/` with embed directive
- Static assets in `handlers/static/` 
- JavaScript modules in `handlers/static/js/`
- CSS with theme variables in `handlers/static/css/`

## Naming Conventions
- Go: camelCase for private, PascalCase for public
- CSS: kebab-case for classes, CSS custom properties with double dashes
- JavaScript: camelCase for variables/functions, PascalCase for classes
- Templates: kebab-case for filenames

## Error Handling
- Go: Return errors, log with context
- Frontend: User-friendly error messages with snackbar notifications
- HTTP: Proper status codes with meaningful messages

## Authentication
- Context-based auth state
- Support for multiple providers (shared secret, OAuth2)
- Middleware for route protection

## Accessibility
- Semantic HTML elements
- ARIA labels and roles
- Skip links for keyboard navigation
- High contrast color schemes (WCAG AA compliance)
- Screen reader friendly alt text