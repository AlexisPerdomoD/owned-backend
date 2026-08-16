import { createSignal, onMount } from 'solid-js'

import { apiPaginateUsrs } from '@/entities/usrs/api'

export function usePaginateUsrs({ pageSize = 20 } = {}) {
    /** @type {import("solid-js").Signal<import('@/entities/usrs').Usr[]>}*/
    const [usrs, setUsrs] = createSignal([])
    /** @type {import("solid-js").Signal<number>}*/
    const [page, setPage] = createSignal(1)
    /** @type {import("solid-js").Signal<number>}*/
    const [total, setTotal] = createSignal(0)
    /** @type {import("solid-js").Signal<boolean>} */
    const [loading, setLoading] = createSignal(true)
    /** @type {import("solid-js").Signal<string>} */
    const [search, setSearch] = createSignal('')
    /** @type {import("solid-js").Signal<string | null>} */
    const [role, setRole] = createSignal(null)
    /**
     * @param {number} pg
     * @param {string} term
     * @param {string | undefined} roleFilter
     */
    const fetch = async (pg, term, roleFilter) => {
        setLoading(true)
        try {
            const res = await apiPaginateUsrs(pg, pageSize, term, roleFilter)
            setUsrs(res.data)
            setTotal(res.total_count)
            setPage(res.page)
        } finally {
            setLoading(false)
        }
    }

    onMount(() => {
        fetch(1, '', undefined)
    })

    return {
        usrs,
        page,
        total,
        loading,
        search,
        setSearch,
        role,
        setRole,
        goTo: pg => fetch(pg, search(), role()),
        refresh: () => fetch(page(), search(), role()),
        setPage: pg => fetch(pg, search(), role())
    }
}
