# Ownned Web Client (web/app)

SolidJS frontend of the Ownned monorepo. Global/backend context: see `../AGENTS.md`

## Commands

```bash
pnpm install
pnpm dev                  # Dev server on :5173
pnpm build                # Build to web/app/dist (copied to web/dist by make build-web)
```

## Stack

- SolidJS 1.9, @solidjs/router 0.16, TailwindCSS 4 (no `tailwind.config.js`, `@import 'tailwindcss'` in CSS), Vite 8, Zod 4
- ESLint only, prettier integrated via eslint config (no separate config)

## Path Aliases

- `@/*` → `src/`, `@entities/*`, `@features/*`, `@pages/*`, `@shared/*`

## Structure & Naming

- `entities/{feature}/api` — DTOs and API functions
- `features/{feature}/{providers,ui,usecase}` — feature-specific logic
- `pages/{View}/index.jsx` — routes
- `shared/{ui,api,config}` — design system and shared helpers
- Components PascalCase, files kebab-case (`login-form.jsx`), hooks `use*` camelCase
- Design system in `shared/ui`, imported via `@/shared/ui`

## SolidJS Reactivity (read before editing components)

- Props that receive dynamic expressions become getters on the props object, evaluated lazily per access.
- NEVER destructure props — or rest-destructure then spread —— it snapshots values once at mount and they never update.
  Real bug: `LoginForm`'s showPassword toggle did nothing because `Input.jsx` and `Button.jsx` rest-destructured props, freezing `type` and `children` as static values.
- To forward props to a native element use `splitProps(props, [...localKeys])` and spread `{...rest}`.
- When consuming your own props, access them directly (`props.x`) inside JSX, not via destructuring.

## Reference

- `build_indications.md` — full stack, structure, routes, views and design-system documentation.
