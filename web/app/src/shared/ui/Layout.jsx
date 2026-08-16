import { Divider } from './Atoms'

/**
 * Layout de aplicación con sidebar colapsable.
 * Compuesto por `AppShell`, `AppShell.Sidebar` y `AppShell.Content`.
 *
 * @example
 * <AppShell>
 *   <AppShell.Sidebar>
 *     <NavItem href="/nodes" icon={<FolderIcon />}>Archivos</NavItem>
 *     <NavItem href="/groups">Grupos</NavItem>
 *   </AppShell.Sidebar>
 *   <AppShell.Content>
 *     <Outlet />
 *   </AppShell.Content>
 * </AppShell>
 */

/**
 * @param {Object} props
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.Element} props.children
 * @returns {import('solid-js').JSX.Element}
 */
export function AppShell(props) {
    return (
        <div class={`flex h-screen overflow-hidden bg-bg ${props.class ?? ''}`}>
            {props.children}
        </div>
    )
}

/**
 * Barra lateral de la app.
 *
 * @param {Object} props
 * @param {string} [props.brand]                  - Nombre/logo de la app.
 * @param {import('solid-js').JSX.Element} [props.footer] - Slot inferior (e.g. avatar de usuario).
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.Element} props.children
 * @returns {import('solid-js').JSX.Element}
 */
AppShell.Sidebar = function Sidebar(props) {
    return (
        <aside
            class={`
                w-52 shrink-0 flex flex-col
                bg-bg-2 border-r border-border
                h-full overflow-y-auto
                ${props.class ?? ''}
            `}
        >
            {/* Brand */}
            {props.brand && (
                <div class="px-4 py-4 border-b border-border-subtle">
                    <span class="font-serif text-base text-ink-dark">
                        {props.brand}
                    </span>
                </div>
            )}

            {/* Nav items */}
            <nav class="flex-1 py-3 px-2 flex flex-col gap-0.5">
                {props.children}
            </nav>

            {/* Footer slot */}
            {props.footer && (
                <>
                    <Divider />
                    <div class="px-3 py-3">{props.footer}</div>
                </>
            )}
        </aside>
    )
}

/**
 * Área de contenido principal scrolleable.
 *
 * @param {Object} props
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.Element} props.children
 * @returns {import('solid-js').JSX.Element}
 */
AppShell.Content = function Content(props) {
    return (
        <main class={`flex-1 overflow-y-auto ${props.class ?? ''}`}>
            {props.children}
        </main>
    )
}

// ─────────────────────────────────────────────────────────────────────────────

/**
 * Ítem de navegación del sidebar.
 *
 * @param {Object} props
 * @param {string} [props.href]
 * @param {boolean} [props.active=false]
 * @param {import('solid-js').JSX.Element} [props.icon]
 * @param {() => void} [props.onClick]
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.Element} props.children
 * @returns {import('solid-js').JSX.Element}
 */
export function NavItem(props) {
    const active_style =
        'bg-surface text-ink-dark border border-border font-normal'
    const inactive_style = 'text-ink hover:bg-surface hover:text-ink-dark'
    const base = `
        flex items-center gap-2.5 w-full
        px-3 py-2 rounded-xs
        text-sm font-sans font-light
        transition-colors duration-[--ease-base]
        cursor-pointer`

    return (
        <>
            {props.href ? (
                <a
                    href={props.href}
                    class={
                        base +
                        ' ' +
                        (props.active ? active_style : inactive_style) +
                        (props.class ? ` ${props.class}` : '')
                    }
                    onClick={props.onClick}
                >
                    {props.icon && (
                        <span
                            class="shrink-0 opacity-60"
                            style="font-size:15px"
                        >
                            {props.icon}
                        </span>
                    )}
                    {props.children}
                </a>
            ) : (
                <button
                    class={
                        base +
                        ' ' +
                        (props.active ? active_style : inactive_style) +
                        (props.class ? ` ${props.class}` : '')
                    }
                    onClick={props.onClick}
                >
                    {props.icon && (
                        <span
                            class="shrink-0 opacity-60"
                            style="font-size:15px"
                        >
                            {props.icon}
                        </span>
                    )}
                    {props.children}
                </button>
            )}
        </>
    )
}

// ─────────────────────────────────────────────────────────────────────────────

/**
 * Cabecera de página estandarizada.
 *
 * @param {Object} props
 * @param {string} props.title
 * @param {string} [props.subtitle]
 * @param {import('solid-js').JSX.Element} [props.actions]  - Slot derecho para botones.
 * @param {import('solid-js').JSX.Element} [props.breadcrumb]
 * @param {() => void} [props.onBack] - Callback para botón atrás.
 * @param {string} [props.backTo] - Ruta para botón atrás (usa navigate si no hay callback).
 * @param {string} [props.class]
 * @returns {import('solid-js').JSX.Element}
 */
export function PageHeader(props) {
    return (
        <div
            class={`flex items-start justify-between gap-4 mb-6 ${props.class ?? ''}`}
        >
            <div class="flex flex-col gap-1">
                {(props.breadcrumb || props.onBack || props.backTo) && (
                    <div class="flex items-center gap-2 mb-1">
                        {props.onBack && (
                            <button
                                type="button"
                                onClick={props.onBack}
                                class="text-sm text-accent hover:underline"
                            >
                                ← Back
                            </button>
                        )}
                        {props.backTo && !props.onBack && (
                            <a
                                href={props.backTo}
                                class="text-sm text-accent hover:underline"
                            >
                                ← Back
                            </a>
                        )}
                    </div>
                )}
                {props.breadcrumb}
                <h1 class="font-[--font-serif] text-2xl text-ink-dark leading-tight">
                    {props.title}
                </h1>
                {props.subtitle && (
                    <p class="text-sm text-muted">{props.subtitle}</p>
                )}
            </div>
            {props.actions && (
                <div class="flex items-center gap-2 shrink-0 mt-1">
                    {props.actions}
                </div>
            )}
        </div>
    )
}
