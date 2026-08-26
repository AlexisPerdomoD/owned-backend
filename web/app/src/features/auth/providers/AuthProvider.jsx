import { createContext, onMount, useContext } from 'solid-js'
import { createStore } from 'solid-js/store'

import { apiGetMe, apiLogin, apiLogout } from '@/entities/usrs/api'
import { getRoutesMap } from '@/shared/config/views'

/**
 * @typedef {Object} AuthCtx
 * @property {Object} state
 * @property {import('@/entities/usrs').Usr | null} state.usr
 * @property {boolean} state.checked
 * @property {Map<string, import('@/shared/config/views').Route>} state.routes
 *
 *
 * @property {(credentials: import('@entities/usrs/api/login').LoginPwdDTO) => Promise<void>} loginPwd
 * @property {() => Promise<void>} logout
 */

/**
 * @type {import('solid-js').Context<AuthCtx>}
 */
const AuthCtx = createContext()

/**
 * @returns {AuthCtx}
 */
export function useAuth() {
    const ctx = useContext(AuthCtx)
    if (!ctx) {
        throw new Error('useAuth must be used within a AuthProvider')
    }

    return ctx
}

/**
 * @param {Object} props
 * @param {import('solid-js').JSX.Element} props.children
 * @returns {import('solid-js').JSX.Element}
 */
export function AuthProvider(props) {
    /**
     *  @type AuthCtx['state']
     */
    const initialState = {
        usr: null,
        checked: false,
        routes: new Map()
    }

    const [state, setState] = createStore(initialState)

    let authCheckVersion = 0

    onMount(
        () => {
            const version = authCheckVersion

            apiGetMe().then(usr => {
                if (version !== authCheckVersion) {
                    return
                }

                setState({
                    usr,
                    checked: true,
                    routes: usr ? getRoutesMap(usr.role) : new Map()
                })
            })
        }
        // TODO: check if here errors need to be handled
    )

    /**
     * Login the user.
     * @param { import('@/entities/usrs/api/loginUsr').LoginPwdDTO } sanitizedCredentials
     *
     */
    const loginPwd = sanitizedCredentials =>
        apiLogin(sanitizedCredentials).then(usr => {
            authCheckVersion += 1
            setState({ usr, checked: true, routes: getRoutesMap(usr.role) })
            return usr
        })
    // TODO: check if here errors need to be handled

    const logout = async () => {
        try {
            await apiLogout()
        } finally {
            setState({ usr: null, checked: true, routes: new Map() })
        }
    }

    return (
        <AuthCtx.Provider value={{ state, loginPwd, logout }}>
            {props.children}
        </AuthCtx.Provider>
    )
}
