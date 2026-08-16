import { Show, createEffect, mergeProps, onCleanup, splitProps } from 'solid-js'
import { Portal } from 'solid-js/web'

/**
 * Modal accesible con Portal y trap de foco básico.
 * Cierra con Escape y click en el backdrop.
 *
 * @example
 * const [open, setOpen] = createSignal(false)
 *
 * <Modal open={open()} onClose={() => setOpen(false)} title="Confirmar acción">
 *   <p>¿Estás seguro?</p>
 *   <Modal.Footer>
 *     <Button variant="ghost" onClick={() => setOpen(false)}>Cancelar</Button>
 *     <Button variant="primary">Confirmar</Button>
 *   </Modal.Footer>
 * </Modal>
 */

const sizes = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-2xl'
}

/**
 * @param {Object} _props
 * @param {boolean} _props.open                           - Controla visibilidad.
 * @param {() => void} _props.onClose                    - Callback de cierre.
 * @param {string} [_props.title]
 * @param {string} [_props.description]
 * @param {'sm' | 'md' | 'lg' | 'xl'} [_props.size='md']
 * @param {boolean} [_props.closeOnBackdrop=true]
 * @param {import('solid-js').JSX.Element} _props.children
 * @returns {import('solid-js').JSX.Element}
 */
export function Modal(_props) {
    const props = mergeProps(
        {
            open: false,
            onClose: () => {},
            title: '',
            description: '',
            size: 'md',
            closeOnBackdrop: true
        },

        _props
    )
    createEffect(() => {
        if (!props.open) return
        const onKey = e => {
            if (e.key === 'Escape') props.onClose()
        }
        document.addEventListener('keydown', onKey)
        document.body.style.overflow = 'hidden'
        onCleanup(() => {
            document.removeEventListener('keydown', onKey)
            document.body.style.overflow = ''
        })
    })

    return (
        <Portal>
            <Show when={open}>
                {/* Backdrop */}
                <div
                    class="fixed inset-0 z-50 flex items-center justify-center p-4"
                    style="background: rgba(42, 37, 32, 0.5)"
                    onClick={e => {
                        const isCloseOnBackdropClick =
                            props.closeOnBackdrop &&
                            e.target === e.currentTarget

                        if (!isCloseOnBackdropClick) {
                            return
                        }

                        props.onClose()
                    }}
                    role="dialog"
                    aria-modal="true"
                    aria-labelledby={props.title ? 'modal-title' : undefined}
                >
                    {/* Panel */}
                    <div
                        class={`
                            relative w-full ${sizes[props.size]}
                            bg-surface border border-border
                            rounded-sm shadow-sm
                            flex flex-col
                        `}
                    >
                        {/* Header */}
                        {(props.title || props.description) && (
                            <div class="px-5 py-4 border-b border-[--color-border-subtle]">
                                {props.title && (
                                    <h2
                                        id="modal-title"
                                        class="font-serif text-lg text-ink-dark leading-snug"
                                    >
                                        {props.title}
                                    </h2>
                                )}
                                {props.description && (
                                    <p class="mt-1 text-sm text-muted">
                                        {props.description}
                                    </p>
                                )}
                            </div>
                        )}

                        {/* Close button */}
                        <button
                            onClick={() => props.onClose()}
                            aria-label="Cerrar"
                            class="
                                absolute top-3.5 right-4
                                text-muted hover:text-ink
                                transition-colors duration-[--ease-base]
                                cursor-pointer
                            "
                        >
                            <svg
                                width="14"
                                height="14"
                                viewBox="0 0 14 14"
                                fill="none"
                            >
                                <path
                                    d="M1 1l12 12M13 1L1 13"
                                    stroke="currentColor"
                                    stroke-width="1.2"
                                    stroke-linecap="round"
                                />
                            </svg>
                        </button>

                        {/* Body */}
                        <div class="px-5 py-4 flex-1">{props.children}</div>
                    </div>
                </div>
            </Show>
        </Portal>
    )
}

/**
 * Pie del modal, alineado a la derecha.
 *
 * @param {Object} props
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.HTMLAttributes<HTMLDivElement>} props
 * @returns {import('solid-js').JSX.Element}
 */
Modal.Footer = function ModalFooter(props) {
    const [local, rest] = splitProps(props, ['class', 'children'])
    return (
        <div
            class={`
                flex items-center justify-end gap-2
                pt-3 mt-2 border-t border-[--color-border-subtle]
                ${local.class ?? ''}
            `}
            {...rest}
        >
            {local.children}
        </div>
    )
}
