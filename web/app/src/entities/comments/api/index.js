import z from 'zod'

import { reqJSON } from '@/shared/api/client'

export class CreateCommentDTO {
    static #schema = z.strictObject({
        content: z
            .string('Content must be a string.')
            .min(1, 'Content is required.')
            .max(10000, 'Content must be at most 10000 characters.')
    })

    static get schema() {
        return structuredClone(CreateCommentDTO.#schema)
    }

    constructor({ content }) {
        this.content = content
    }

    toJSON() {
        return { node_id: this.node_id, content: this.content }
    }

    /**
     * Validates the DTO.
     * @returns {[false, import('zod').ZodIssue[]] | [true, CreateCommentDTO]}
     */
    static build(content) {
        const r = CreateCommentDTO.#schema.safeParse({ content })
        if (!r.success) {
            const issues = {}
            r.error.issues.forEach(issue => {
                issues[issue.path.at(-1)?.toString() ?? 'general'] =
                    issue.message
            })
            return [false, issues]
        }
        return [true, new CreateCommentDTO(r.data)]
    }
}

export const buildCreateCommentDTO = CreateCommentDTO.build

/**
 * Get comments for a node.
 * @param {string} nodeId
 * @returns {Promise<import('@/entities/nodes').NodeComment[]>}
 */
export async function apiGetComments(nodeId) {
    return await reqJSON(`/api/v1/nodes/${nodeId}/comments`)
}

/**
 * Create a comment.
 * @param {string} nodeId
 * @param {CreateCommentDTO} dto
 * @returns {Promise<import('@/entities/nodes').NodeComment>}
 */
export async function apiCreateComment(nodeId, dto) {
    return await reqJSON(`/api/v1/nodes/${nodeId}/comments`, {
        method: 'POST',
        body: JSON.stringify(dto)
    })
}

/**
 * Delete a comment.
 * @param {string} commentId
 * @returns {Promise<{message: string}>}
 */
export async function apiDeleteComment(commentId) {
    return await reqJSON(`/api/v1/comments/${commentId}`, {
        method: 'DELETE'
    })
}
