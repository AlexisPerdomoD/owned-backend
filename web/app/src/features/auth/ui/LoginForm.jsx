import { createEffect, createSignal } from 'solid-js'

import { Button, Input, toast } from '@/shared/ui'

import { usePwdLogin } from '../usecase/usePwdLogin'

export function LoginForm() {
    const { fields, issues, loading, setField, submit } = usePwdLogin()
    const onUsernameChange = e => setField('username', e.target.value)
    const onPasswordChange = e => setField('password', e.target.value)
    const [showPassword, setShowPassword] = createSignal(false)

    const TOAST_CODE = 'LOGIN_FORM_GENERAL_ISSUES'
    createEffect(() => {
        const generalIssues = issues().general
        if (!generalIssues) {
            return
        }

        toast.clearCode(TOAST_CODE)
        toast.error(generalIssues, TOAST_CODE)
    })

    return (
        <>
            <form onSubmit={submit} class="flex flex-col gap-4 w-full">
                <Input
                    label="Username"
                    placeholder="dropdibemol@gmail.com"
                    value={fields.username}
                    onInput={onUsernameChange}
                    error={issues.username}
                    hint="Introduce your email"
                    type="email"
                    autocomplete="email"
                />
                <Input
                    label="Password"
                    placeholder="**********"
                    value={fields.password}
                    onInput={onPasswordChange}
                    error={issues.password}
                    hint="Introduce your password"
                    type={showPassword() ? 'text' : 'password'}
                    class="flex-1"
                    autocomplete="current-password"
                    minlength={8}
                    suffix={
                        <Button
                            type="button"
                            class="w-12 shrink-0 rounded-sm"
                            variant="ghost"
                            size="sm"
                            onClick={() => setShowPassword(prev => !prev)}
                            children={showPassword() ? 'Hide' : 'Show'}
                        />

                    }
                />


                <Button
                    type="submit"
                    variant="primary"
                    size="md"
                    loading={loading()}
                    class="w-full mt-1"
                    children="Login"
                />
            </form>
        </>
    )
}
