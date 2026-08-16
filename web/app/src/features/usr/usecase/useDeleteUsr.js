import { createSignal } from 'solid-js'

import { apiDeleteUsr } from '@/entities/usrs/api'
import { toast } from '@/shared/ui'

export function useDeleteUsr() {
    const [loading, setLoading] = createSignal(false)

    const remove = async usrId => {
        setLoading(true)
        try {
            await apiDeleteUsr(usrId)
            toast({ type: 'success', message: 'User deleted.' })
            return [true, null]
        } catch (e) {
            return [false, { general: e.message }]
        } finally {
            setLoading(false)
        }
    }

    return { remove, loading }
}
