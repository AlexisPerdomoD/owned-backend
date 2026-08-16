import { Show, createEffect, createSignal } from 'solid-js'

import { apiDownloadDoc } from '@/entities/docs/api'
import { CommentForm, CommentsList } from '@/features/comments/ui/CommentsList'
import {
    useCreateComment,
    useDeleteComment,
    useGetComments
} from '@/features/comments/usecase'
import { DocCard, NodeList, UploadDropzone } from '@/features/node/ui'
import {
    useCreateDoc,
    useCreateFolder,
    useDeleteDoc,
    useDeleteNode,
    useGetNode,
    useUpdateNode
} from '@/features/node/usecase'
import { Button, PageHeader, Spinner, Tabs, toast } from '@/shared/ui'
import { useNavigate, useParams } from '@solidjs/router'

function NodeOverview(props) {
    return (
        <section class="flex flex-col items-center">
            <section class="max-w-md w-full">
                <section class="mb-4 py-4">
                    <h1 class="text-2xl font-semibold text-center font-serif pb-2">
                        {props.node.name}
                    </h1>
                    <h3 class="text-sm text-center font-serif pb-2">
                        {props.node.type === 'folder' ? 'Folder' : 'File'}
                    </h3>
                    <p class="text-sm text-center">{props.node.description}</p>
                    <p class="text-sm text-center">
                        Created {props.node.created_at}
                    </p>
                    <Show
                        when={props.node.updated_at !== props.node.created_at}
                    >
                        <p class="text-sm text-center">
                            Last updated {props.node.updated_at}
                        </p>
                    </Show>
                </section>
                <section class="flex justify-center gap-2">
                    <Button variant="ghost" size="sm" onClick={props.onEdit}>
                        Edit
                    </Button>
                    <Button variant="danger" size="sm" onClick={props.onDelete}>
                        Delete
                    </Button>
                </section>
            </section>
        </section>
    )
}
/**
 * @param {Object} props
 * @param {boolean} props.open
 * @param {import('@/entities/nodes').Node | null} props.node
 * @param {(name: string, description: string) => void} props.onSubmit
 * @param {() => void} props.onClose
 * @param {boolean} props.loading
 */
function NodeForm(props) {
    const [name, setName] = createSignal('')
    const [description, setDescription] = createSignal('')

    createEffect(() => {
        if (props.open && props.node) {
            setName(props.node.name ?? '')
            setDescription(props.node.description ?? '')
        } else if (props.open) {
            setName('')
            setDescription('')
        }
    })

    const handleSubmit = e => {
        e.preventDefault()
        if (!name().trim() || !props.onSubmit) {
            return
        }

        props.onSubmit(name().trim(), description().trim())
    }

    return (
        <>
            {props.open && (
                <div
                    class="fixed inset-0 z-50 flex items-center justify-center bg-black/30"
                    onClick={e => props.onClose(e)}
                >
                    <div
                        class="bg-surface border border-border rounded-md w-full max-w-md p-6"
                        onClick={e => e.stopPropagation()}
                    >
                        <h3 class="font-serif text-lg text-ink-dark mb-4">
                            {props.node ? 'Edit Node' : 'New Folder'}
                        </h3>

                        <form
                            onSubmit={handleSubmit}
                            class="flex flex-col gap-4"
                        >
                            <div class="flex flex-col gap-1">
                                <label class="text-xs text-muted uppercase tracking-wide">
                                    Name *
                                </label>
                                <input
                                    type="text"
                                    value={name()}
                                    onInput={e => setName(e.target.value)}
                                    maxLength={255}
                                    class="
                                w-full font-sans font-light text-sm
                                text-ink-dark
                                bg-bg border border-border rounded-xs
                                px-3 py-2
                                focus:outline-none focus:border-ink
                            "
                                    required
                                />
                            </div>

                            <div class="flex flex-col gap-1">
                                <label class="text-xs text-muted uppercase tracking-wide">
                                    Description
                                </label>
                                <textarea
                                    value={description()}
                                    onInput={e =>
                                        setDescription(e.target.value)
                                    }
                                    rows={2}
                                    maxLength={255}
                                    class="
                                w-full font-sans font-light text-sm
                                text-ink-dark
                                bg-bg border border-border rounded-xs
                                px-3 py-2 resize-none
                                focus:outline-none focus:border-ink
                            "
                                />
                            </div>

                            <div class="flex justify-end gap-2">
                                <Button
                                    type="button"
                                    variant="ghost"
                                    onClick={props.onClose}
                                >
                                    Cancel
                                </Button>
                                <Button type="submit" loading={props.loading}>
                                    {props.node ? 'Update' : 'Create'}
                                </Button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </>
    )
}

export function NodeView() {
    const params = useParams()
    const navigate = useNavigate()
    const { node, loading } = useGetNode(() => params.id)
    const { create, loading: creating } = useCreateFolder()
    const { update, loading: updating } = useUpdateNode()
    const { remove, loading: deleting } = useDeleteNode()

    const { upload, loading: uploading } = useCreateDoc(() => params.id)
    const { remove: deleteDoc, loading: deletingDoc } = useDeleteDoc()

    const [showForm, setShowForm] = createSignal(false)
    const [editingNode, setEditingNode] = createSignal(null)
    const [showCreate, setShowCreate] = createSignal(false)
    const [newName, setNewName] = createSignal('')
    const [newDesc, setNewDesc] = createSignal('')
    const [tab, setTab] = createSignal('contents')

    const {
        comments,
        loading: commentsLoading,
        refresh: refreshComments
    } = useGetComments(() => params.id)
    const { create: addComment, loading: addingComment } = useCreateComment()
    const { remove: deleteComment, loading: deletingComment } =
        useDeleteComment()

    const handleEdit = () => {
        setEditingNode(node())
        setShowForm(true)
    }

    const handleDelete = async () => {
        if (confirm(`Delete "${node().name}"?`)) {
            const [success] = await remove(node().id)
            if (success) {
                navigate('/nodes')
            }
        }
    }

    const handleUpdate = async (name, description) => {
        const [success] = await update(node().id, name, description)
        if (success) {
            setShowForm(false)
            setEditingNode(null)
        }
    }

    const handleCreate = async () => {
        const [success] = await create(node().id, newName(), newDesc())
        if (success) {
            setShowCreate(false)
            setNewName('')
            setNewDesc('')
        }
    }

    const handleUpload = async file => {
        const [success, err] = await upload(file)
        if (!success) {
            toast({ type: 'error', message: err?.general ?? 'Upload failed' })
            return
        }
        toast({ type: 'success', message: 'File uploaded successfully' })
    }

    const handleDownload = async doc => {
        try {
            await apiDownloadDoc(doc.id, doc.title)
        } catch (e) {
            toast({ type: 'error', message: e.message ?? 'Download failed' })
        }
    }

    const handleDeleteDoc = async doc => {
        if (!confirm(`Delete "${doc.title}"?`)) return
        const [success, err] = await deleteDoc(doc.id)
        if (!success) {
            toast({ type: 'error', message: err?.general ?? 'Delete failed' })
            return
        }
        toast({ type: 'success', message: 'File deleted successfully' })
    }

    const nodeData = () => node()
    const canEdit = () => nodeData()?.type === 'folder'

    const folderChildren = () =>
        nodeData()?.children?.filter(n => n.type === 'folder') ?? []
    const fileChildren = () =>
        nodeData()?.children?.filter(n => n.type === 'file') ?? []

    return (
        <section class="flex flex-col p-6">
            <Show
                when={!loading()}
                fallback={<Spinner size="lg" class="mx-auto mt-20" />}
            >
                <PageHeader
                    title={nodeData()?.name ?? 'Node'}
                    subtitle={nodeData()?.type}
                    backTo="/nodes"
                />

                <Show when={canEdit()}>
                    <div class="mb-4">
                        <Button onClick={() => setShowCreate(true)}>
                            + New Folder
                        </Button>
                    </div>
                </Show>

                <NodeOverview
                    node={nodeData()}
                    onEdit={handleEdit}
                    onDelete={handleDelete}
                />

                <Show when={nodeData()?.type === 'folder'}>
                    <section class="mt-4">
                        <h3 class="font-serif text-lg mb-2">
                            {nodeData()?.name}'s contents
                        </h3>
                        <NodeList
                            nodes={folderChildren()}
                            loading={loading()}
                            onNodeClick={n => navigate(`/nodes/${n.id}`)}
                        />
                    </section>

                    <section class="mt-4">
                        <h3 class="font-serif text-lg mb-2">Files</h3>
                        <Show
                            when={fileChildren().length > 0}
                            fallback={
                                <p class="text-sm text-[--color-muted]">
                                    No files in this folder
                                </p>
                            }
                        >
                            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                                {fileChildren().map(fileNode => (
                                    <DocCard
                                        doc={fileNode.doc}
                                        onDownload={handleDownload}
                                        onDelete={handleDeleteDoc}
                                    />
                                ))}
                            </div>
                        </Show>
                    </section>

                    <section class="mt-4">
                        <UploadDropzone
                            onUpload={handleUpload}
                            loading={uploading()}
                        />
                    </section>
                </Show>

                <Show when={nodeData()?.type === 'file'}>
                    <DocCard doc={nodeData()?.doc} />
                </Show>

                <Show when={nodeData()?.type === 'folder'}>
                    <div class="mb-6 mt-6">
                        <Tabs
                            items={[
                                { key: 'contents', label: 'Contents' },
                                { key: 'comments', label: 'Comments' }
                            ]}
                            active={tab()}
                            onChange={setTab}
                        />
                    </div>
                </Show>

                <Show
                    when={nodeData()?.type === 'folder' && tab() === 'contents'}
                >
                    <section class="mt-4">
                        <h3 class="font-serif text-lg mb-2">Folders</h3>
                        <NodeList
                            nodes={folderChildren()}
                            loading={loading()}
                            onNodeClick={n => navigate(`/nodes/${n.id}`)}
                        />
                    </section>
                </Show>

                <Show
                    when={nodeData()?.type === 'folder' && tab() === 'comments'}
                >
                    <section class="mt-4">
                        <h3 class="font-serif text-lg mb-2">Comments</h3>
                        <CommentsList
                            comments={comments()}
                            loading={commentsLoading()}
                            onDelete={cId =>
                                deleteComment(cId).then(
                                    isSuccess => isSuccess && refreshComments()
                                )
                            }
                        />
                        <div class="mt-4">
                            <CommentForm
                                loading={addingComment()}
                                onSubmit={content =>
                                    addComment(node().id, content).then(
                                        isSuccess =>
                                            isSuccess && refreshComments()
                                    )
                                }
                            />
                        </div>
                    </section>
                </Show>

                <NodeForm
                    open={showForm()}
                    node={editingNode()}
                    loading={updating() || deleting()}
                    onSubmit={handleUpdate}
                    onClose={() => {
                        setShowForm(false)
                        setEditingNode(null)
                    }}
                />

                <NodeForm
                    open={showCreate()}
                    node={null}
                    loading={creating()}
                    onSubmit={(name, desc) => {
                        setNewName(name)
                        setNewDesc(desc)
                        handleCreate()
                    }}
                    onClose={() => setShowCreate(false)}
                />
            </Show>
        </section>
    )
}
