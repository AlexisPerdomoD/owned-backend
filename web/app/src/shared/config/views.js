import {
    LIMITED_USR_ROLE,
    NORMAL_USR_ROLE,
    SUPER_USR_ROLE
} from '@/entities/usrs'

/**
 * @typedef {Object} Route
 * @property {string} path
 * @property {string} label
 * @property {string} icon
 * @property {string[]} allowRoles
 */

/**
 * @type {Route[]}
 */
export const routes = [
    {
        path: '/nodes',
        label: 'Files',
        icon: '📁',
        allowRoles: [NORMAL_USR_ROLE, SUPER_USR_ROLE, LIMITED_USR_ROLE]
    },
    {
        path: '/groups',
        label: 'Groups',
        icon: '👥',
        allowRoles: [SUPER_USR_ROLE, NORMAL_USR_ROLE]
    },
    {
        path: '/usrs',
        label: 'Users',
        icon: '👤',
        allowRoles: [SUPER_USR_ROLE]
    }
]
/**
 * @param {import('@entities/usrs').UsrRole} role
 * @returns {Map<string, Route>}
 */
export function getRoutesMap(role) {
    const result = new Map()
    for (const route of routes) {
        if (!route.allowRoles.includes(role)) {
            continue
        }

        result.set(route.path, route)
    }

    return result
}
