import { Show } from 'solid-js'

import { Spinner } from '@/shared/ui/Atoms'
import { Navigate, useLocation } from '@solidjs/router'

import { useAuth } from '../providers/AuthProvider'

/**
 * @param {Object} props
 * @param {import('solid-js').JSX.Element} props.children
 * @returns {import('solid-js').JSX.Element}
 */
export function EnsureAuthenticateRoute(props) {
    const location = useLocation()
    const { state } = useAuth()
    const isAllowedRoute = () =>
        Array.from(state.routes.keys()).some(
            route =>
                location.pathname === route ||
                location.pathname.startsWith(route + '/')
        )

    return (
        <>
            <Show
                when={state.checked}
                fallback={
                    <section class="flex items-center justify-center h-screen">
                        <Spinner size="lg" />
                    </section>
                }
            >
                <Show
                    when={
                        state.usr && isAllowedRoute()
                    }
                    fallback={() => <Navigate href="/login" />}
                >
                    {props.children}
                </Show>
            </Show>
        </>
    )
}
