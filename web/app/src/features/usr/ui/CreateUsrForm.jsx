import { For, createSignal } from 'solid-js'

import { Button } from '@/shared/ui'
import { ROLES } from '@/entities/usrs'

/**
 * @param {Object} props
 * @param {(username: string, password: string, firstname: string, lastname: string, role: string, roles: string[]) => void} props.onSubmit
 * @param {() => void} [props.onClose]
 * @param {boolean} [props.open]
 * @param {import('solid-js').JSX.Element} [props.children]
 */
export function CreateUsrForm(props) {
    let usernameInput
    let passwordInput
    let firstnameInput
    let lastnameInput
    let roleInput

    const [errors, setErrors] = createSignal({})

    const handleSubmit = e => {
        e.preventDefault()
        setErrors({})
        if (!props.onSubmit) {
            setErrors({ general: 'No handler provided. Aborting.' })
            return
        }

        const username = usernameInput?.value?.trim()
        const password = passwordInput?.value
        const firstname = firstnameInput?.value?.trim()
        const lastname = lastnameInput?.value?.trim()
        const role = roleInput?.value

        if (!username || !password || !firstname || !lastname || !role) {
            setErrors({
                general: 'All fields are required.'
            })
            return
        }

        props.onSubmit(username, password, firstname, lastname, role, [])
    }

    return (
        <>
            {props.open && (
                <div
                    class="fixed inset-0 z-50 flex items-center justify-center bg-black/30"
                    onClick={() => props.onClose?.()}
                >
                    <div
                        class="bg-surface border border-border rounded-md w-full max-w-md p-6"
                        onClick={e => e.stopPropagation()}
                    >
                        <h3 class="font-serif text-lg text-ink-dark mb-4">
                            New User
                        </h3>

                        <form
                            onSubmit={handleSubmit}
                            class="flex flex-col gap-4"
                        >
                            <div class="flex flex-col gap-1">
                                <label class="text-xs text-muted uppercase tracking-wide">
                                    Email *
                                </label>
                                <input
                                    ref={r => (usernameInput = r)}
                                    type="email"
                                    placeholder="email@example.com"
                                    class="
                                w-full font-sans font-light text-sm
                                text-[--color-ink-dark]
                                bg-[--color-bg] border border-[--color-border] rounded-xs
                                px-3 py-2
                                focus:outline-none focus:border-[--color-ink]
                            "
                                    required
                                />
                            </div>

                            <div class="flex gap-2">
                                <div class="flex flex-col gap-1 flex-1">
                                    <label class="text-xs text-[--color-muted] uppercase tracking-wide">
                                        First Name *
                                    </label>
                                    <input
                                        ref={r => (firstnameInput = r)}
                                        type="text"
                                        placeholder="John"
                                        class="
                                    w-full font-sans font-light text-sm
                                    text-[--color-ink-dark]
                                    bg-[--color-bg] border border-[--color-border] rounded-xs
                                    px-3 py-2
                                    focus:outline-none focus:border-[--color-ink]
                                "
                                        required
                                    />
                                </div>
                                <div class="flex flex-col gap-1 flex-1">
                                    <label class="text-xs text-[--color-muted] uppercase tracking-wide">
                                        Last Name *
                                    </label>
                                    <input
                                        ref={r => (lastnameInput = r)}
                                        type="text"
                                        placeholder="Doe"
                                        class="
                                    w-full font-sans font-light text-sm
                                    text-[--color-ink-dark]
                                    bg-[--color-bg] border border-[--color-border] rounded-xs
                                    px-3 py-2
                                    focus:outline-none focus:border-[--color-ink]
                                "
                                        required
                                    />
                                </div>
                            </div>

                            <div class="flex flex-col gap-1">
                                <label class="text-xs text-[--color-muted] uppercase tracking-wide">
                                    Role *
                                </label>
                                <select
                                    ref={r => (roleInput = r)}
                                    defaultValue="normal_usr_role"
                                    class="
                                w-full font-sans font-light text-sm
                                text-[--color-ink-dark]
                                bg-[--color-bg] border border-[--color-border] rounded-xs
                                px-3 py-2
                                focus:outline-none focus:border-[--color-ink]
                            "
                                    required
                                >
                                    <For each={ROLES}>
                                        {opt => (
                                            <option value={opt.value}>
                                                {opt.label}
                                            </option>
                                        )}
                                    </For>
                                </select>
                            </div>

                            <div class="flex flex-col gap-1">
                                <label class="text-xs text-[--color-muted] uppercase tracking-wide">
                                    Password *
                                </label>
                                <input
                                    ref={r => (passwordInput = r)}
                                    type="password"
                                    minLength={8}
                                    placeholder="Min 8 characters"
                                    class="
                                w-full font-sans font-light text-sm
                                text-[--color-ink-dark]
                                bg-[--color-bg] border border-[--color-border] rounded-xs
                                px-3 py-2
                                focus:outline-none focus:border-[--color-ink]
                            "
                                    required
                                />
                            </div>

                            {errors()?.general && (
                                <p class="text-xs text-[--color-danger]">
                                    {errors().general}
                                </p>
                            )}

                            <div class="flex justify-end gap-2">
                                <Button
                                    type="button"
                                    variant="ghost"
                                    onClick={props.onClose}
                                >
                                    Cancel
                                </Button>
                                <Button type="submit" loading={props.loading}>
                                    Create
                                </Button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </>
    )
}
