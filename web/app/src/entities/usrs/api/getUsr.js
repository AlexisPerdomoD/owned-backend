import { reqJSON } from '@/shared/api/client'

/**
 * @param {string} usrId
 * @returns {Promise<import('@/entities/usrs').Usr>}
 */
export async function apiGetUsr(usrId) {
    return await reqJSON(`/api/v1/usrs/${usrId}`)
}
