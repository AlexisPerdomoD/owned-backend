import { useAuth } from '@/features/auth/providers/AuthProvider'
import { Button } from '@/shared/ui'
import { A, useLocation } from '@solidjs/router'

/**
 * @type {object} NavItem
 * @property {string} path
 * @property {string} label
 * @property {string} icon
 */

/**
 * @param {Object} props
 * @param {boolean} [props.collapsed]
 * @param {Array<NavItem>} [props.navItems]
 * @returns {import('solid-js').JSX.Element}
 */
export function Navbar(props) {
    const location = useLocation()
    const { state, logout } = useAuth()

    const isActive = path =>
        location.pathname === path || location.pathname.startsWith(path + '/')

    const handleLogout = () => {
        logout().finally(() => {
            window.location.href = '/login'
        })
    }

    return (
        <nav
            class={`
                flex items-center justify-between
                h-14 px-4
                bg-bg-2 border-b border-border
                ${props.collapsed ? 'w-14 flex-col' : ''}
            `}
        >
            <div
                class={`flex items-center gap-1 ${props.collapsed ? 'flex-col' : ''}`}
            >
                {props.navItems?.map(item => (
                    <A
                        href={item.path}
                        class={`
                        flex items-center gap-2 px-3 py-2 rounded-xs
                        text-sm font-light
                        transition-colors duration-[--ease-base]
                        ${
                            isActive(item.path)
                                ? 'bg-surface text-ink-dark'
                                : 'text-ink hover:bg-surface'
                        }
                    `}
                    >
                        <span style="font-size:16px">📁</span>
                        {!props.collapsed && <span>{item.label}</span>}
                    </A>
                ))}
            </div>

            <div class="flex items-center gap-3">
                <span class="text-sm text-muted">
                    {props.collapsed ? '' : state.usr?.username}
                </span>
                <Button variant="ghost" size="sm" onClick={handleLogout}>
                    {props.collapsed ? '↩' : 'Logout'}
                </Button>
            </div>
        </nav>
    )
}
