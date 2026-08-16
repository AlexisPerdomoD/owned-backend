import { reqJSON } from '@/shared/api/client'

/**
 * @param {string} usrId
 * @returns {Promise<{message: string}>}
 */
export async function apiDeleteUsr(usrId) {
    return await reqJSON(`/api/v1/usrs/${usrId}`, {
        method: 'DELETE'
    })
}
