import { splitProps } from 'solid-js'

/**
 * Select nativo estilizado.
 *
 * @param {Object} props
 * @param {string} [props.label]
 * @param {string} [props.hint]
 * @param {string} [props.error]
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.SelectHTMLAttributes<HTMLSelectElement>} props
 * @returns {import('solid-js').JSX.Element}
 */
export function Select(props) {
    const [local, rest] = splitProps(props, [
        'label',
        'hint',
        'error',
        'class',
        'id',
        'children',
        'onChange',
        'value'
    ])

    const selectId =
        local.id ?? `select-${Math.random().toString(36).slice(2, 7)}`

    return (
        <div class={`flex flex-col gap-1 ${local.class}`}>
            {local.label && (
                <label
                    for={selectId}
                    class="text-xs text-muted tracking-wide uppercase"
                >
                    {local.label}
                </label>
            )}

            <div class="relative">
                <select
                    id={selectId}
                    onChange={e => local.onChange?.(e)}

                    class={`
                        w-full appearance-none
                        font-sans font-light text-sm
                        text-ink-dark
                        bg-bg border rounded-xs
                        px-3 py-2 pr-8
                        transition-colors duration-[--ease-base]
                        focus:outline-none focus:border-ink
                        disabled:opacity-40 disabled:cursor-not-allowed
                        ${local.error ? 'border-red-400' : 'border-border hover:border-muted'}
                    `}
                    {...rest}
                >
                    {local.children}
                </select>
                {/* chevron decorativo */}
                <span class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-[--color-muted]">
                    <svg width="10" height="6" viewBox="0 0 10 6" fill="none">
                        <path
                            d="M1 1l4 4 4-4"
                            stroke="currentColor"
                            stroke-width="1.2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        />
                    </svg>
                </span>
            </div>

            {local.error && <p class="text-xs text-red-500">{local.error}</p>}
            {!local.error && local.hint && (
                <p class="text-xs text-muted">{local.hint}</p>
            )}
        </div>
    )
}

// ─────────────────────────────────────────────────────────────────────────────

/**
 * Área de texto multilínea.
 *
 * @param {Object} props
 * @param {string} [props.label]
 * @param {string} [props.hint]
 * @param {string} [props.error]
 * @param {number} [props.rows=4]
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.TextareaHTMLAttributes<HTMLTextAreaElement>} props
 * @returns {import('solid-js').JSX.Element}
 */
export function Textarea(props) {
    const [local, rest] = splitProps(props, [
        'label',
        'hint',
        'error',
        'rows',
        'class',
        'id'
    ])
    const textareaId =
        local.id ?? `textarea-${Math.random().toString(36).slice(2, 7)}`

    return (
        <div class={`flex flex-col gap-1 ${local.class ?? ''}`}>
            {local.label && (
                <label
                    for={textareaId}
                    class="text-xs text-muted tracking-wide uppercase"
                >
                    {local.label}
                </label>
            )}

            <textarea
                id={textareaId}
                rows={local.rows !== undefined ? local.rows : 4}
                class={`
                    w-full resize-y
                    font-sans font-light text-sm
                    text-ink-dark placeholder:text-muted
                    bg-surface border rounded-xs
                    px-3 py-2
                    transition-colors duration-[--ease-base]
                    focus:outline-none focus:border-ink
                    disabled:opacity-40 disabled:cursor-not-allowed
                    ${local.error ? 'border-red-400' : 'border-border hover:border-muted'}
                `}
                {...rest}
            />

            {local.error && <p class="text-xs text-red-500">{local.error}</p>}
            {!local.error && local.hint && (
                <p class="text-xs text-muted">{local.hint}</p>
            )}
        </div>
    )
}

// ─────────────────────────────────────────────────────────────────────────────

/**
 * Checkbox accesible con label integrado.
 *
 * @param {Object} props
 * @param {string} [props.label]
 * @param {string} [props.hint]
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.InputHTMLAttributes<HTMLInputElement>} props
 * @returns {import('solid-js').JSX.Element}
 */
export function Checkbox(props) {
    const [local, rest] = splitProps(props, ['label', 'hint', 'class', 'id'])
    const checkId =
        local.id ?? `check-${Math.random().toString(36).slice(2, 7)}`

    return (
        <div class={`flex items-start gap-2.5 ${props.class}`}>
            <input
                type="checkbox"
                id={checkId}
                class={`
                    mt-0.5 w-3.5 h-3.5 shrink-0
                    rounded-xs border border-border
                    bg-surface
                    accent-ink-dark
                    cursor-pointer
                    disabled:opacity-40 disabled:cursor-not-allowed
                `}
                {...rest}
            />
            {(local.label || local.hint) && (
                <div class="flex flex-col gap-0.5">
                    {local.label && (
                        <label
                            for={checkId}
                            class="text-sm text-ink cursor-pointer leading-none"
                        >
                            {local.label}
                        </label>
                    )}
                    {local.hint && (
                        <p class="text-xs text-muted">{local.hint}</p>
                    )}
                </div>
            )}
        </div>
    )
}
