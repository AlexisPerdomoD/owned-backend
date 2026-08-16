import { createSignal } from 'solid-js'

import { apiCreateUsr, buildCreateUsrDTO } from '@/entities/usrs/api'
import { toast } from '@/shared/ui'

export function useCreateUsr() {
    const [loading, setLoading] = createSignal(false)

    const create = async (
        username,
        password,
        firstname,
        lastname,
        role,
        access = []
    ) => {
        const [valid, dto] = buildCreateUsrDTO(
            username,
            password,
            firstname,
            lastname,
            role,
            access
        )
        if (!valid) {
            return [false, dto]
        }

        setLoading(true)
        try {
            const usr = await apiCreateUsr(dto)
            toast({ type: 'success', message: 'User created.' })
            return [true, usr]
        } catch (e) {
            return [false, { general: e.message }]
        } finally {
            setLoading(false)
        }
    }

    return { create, loading }
}
