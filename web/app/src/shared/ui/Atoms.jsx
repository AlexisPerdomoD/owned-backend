import { mergeProps } from 'solid-js'

/**
 * Línea divisora horizontal u vertical.
 *
 * @param {Object} props
 * @param {'horizontal' | 'vertical'} [props.orientation='horizontal']
 * @param {string} [props.class]
 * @returns {import('solid-js').JSX.Element}
 */
export function Divider(props) {
    return (
        <>
            {props.orientation === 'vertical' ? (
                <span
                    class={`inline-block w-px self-stretch bg-border-subtle ${props.class}`}
                    role="separator"
                    aria-orientation="vertical"
                />
            ) : (
                <hr
                    class={`border-none border-t border-border-subtle ${props.class}`}
                    role="separator"
                />
            )}
        </>
    )
}

// ─────────────────────────────────────────────────────────────────────────────

const avatarSizes = {
    sm: 'w-7 h-7 text-xs',
    md: 'w-9 h-9 text-sm',
    lg: 'w-12 h-12 text-base'
}

/**
 * Avatar circular con iniciales como fallback.
 *
 * @param {Object} props
 * @param {string} [props.src]               - URL de la imagen.
 * @param {string} [props.name]              - Nombre completo; se usa para generar iniciales si no hay imagen.
 * @param {'sm' | 'md' | 'lg'} [props.size='md']
 * @param {string} [props.class]
 * @returns {import('solid-js').JSX.Element}
 */
export function Avatar(props) {
    const mergedProps = mergeProps({ name: '', size: 'md', class: '' }, props)

    return (
        <span
            class={`
                inline-flex items-center justify-center
                rounded-full shrink-0 overflow-hidden
                bg-accent-pale text-accent
                font-sans font-normal select-none
                ${avatarSizes[mergedProps.size]} ${mergedProps.class}
            `}
        >
            {mergedProps.src ? (
                <img
                    src={mergedProps.src}
                    alt={mergedProps.name}
                    class="w-full h-full object-cover"
                />
            ) : (
                mergedProps.name
                    .trim()
                    .split(/\s+/)
                    .slice(0, 2)
                    .map(w => w[0]?.toUpperCase() ?? '')
                    .join('') || '?'
            )}
        </span>
    )
}

// ─────────────────────────────────────────────────────────────────────────────

const spinnerSizes = {
    sm: 'w-3.5 h-3.5 border',
    md: 'w-5 h-5 border-2',
    lg: 'w-7 h-7 border-2'
}

/**
 * Indicador de carga circular.
 *
 * @param {Object} props
 * @param {'sm' | 'md' | 'lg'} [props.size='md']
 * @param {string} [props.class]
 * @returns {import('solid-js').JSX.Element}
 */
export function Spinner(_props) {
    const props = mergeProps({ size: 'md', class: '' }, _props)
    return (
        <span
            role="status"
            aria-label="Cargando"
            class={`
                inline-block rounded-full animate-spin
                border-border border-t-ink
                ${spinnerSizes[props.size]} ${props.class}
            `}
        />
    )
}

// ─────────────────────────────────────────────────────────────────────────────

/**
 * Estado vacío genérico para listas y tablas sin resultados.
 *
 * @param {Object} props
 * @param {string} [props.title='Sin resultados']
 * @param {string} [props.description]
 * @param {import('solid-js').JSX.Element} [props.action]  - Slot para un CTA, e.g. un Button.
 * @param {import('solid-js').JSX.Element} [props.icon]    - Icono decorativo.
 * @param {string} [props.class]
 * @returns {import('solid-js').JSX.Element}
 */
export function EmptyState(_props) {
    const props = mergeProps({ title: 'Sin resultados', class: '' }, _props)
    return (
        <div
            class={`
                flex flex-col items-center justify-center gap-3
                py-16 px-8 text-center
                ${props.class}
            `}
        >
            {props.icon && (
                <span class="text-border mb-1 opacity-60">{props.icon}</span>
            )}
            <p class="font-serif text-base text-ink">{props.title}</p>
            {props.description && (
                <p class="text-sm text-muted max-w-xs leading-relaxed">
                    {props.description}
                </p>
            )}
            {props.action && <div class="mt-2">{props.action}</div>}
        </div>
    )
}
