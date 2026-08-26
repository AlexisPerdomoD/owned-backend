import { For, createSignal } from 'solid-js'

import { ROLES } from '@/entities/usrs'
import { CreateUsrForm, UsrsTable } from '@/features/usr/ui'
import {
    useCreateUsr,
    useDeleteUsr,
    usePaginateUsrs
} from '@/features/usr/usecase'
import { Button, PageHeader, Pagination, Select, toast } from '@/shared/ui'

export function UsrsView() {
    const {
        usrs,
        page,
        total,
        loading,
        search,
        setSearch,
        role,
        setRole,
        goTo,
        refresh
    } = usePaginateUsrs({ pageSize: 20 })
    const { create, loading: creating } = useCreateUsr()
    const { remove, loading: deleting } = useDeleteUsr()

    const [showForm, setShowForm] = createSignal(false)

    const handleSearch = () => {
        goTo(1)
    }

    const handleCreate = (username, password, firstname, lastname, role) => {
        create(username, password, firstname, lastname, role).then(
            ([success, usr]) => {
                if (!success) {
                    toast({
                        type: 'error',
                        message: usr.general ?? 'Error creating user.'
                    })
                    return
                }

                setShowForm(false)
                refresh()
            }
        )
    }

    const handleDelete = async usr => {
        if (!confirm(`Delete user "${usr.username}"?`)) {
            return
        }

        const [success, issues] = await remove(usr.id)
        if (!success) {
            toast({ type: 'error', message: issues.general ?? 'Delete failed' })
            return
        }

        refresh()
    }

    return (
        <section class="flex flex-col p-6">
            <PageHeader
                title="Manage Users"
                subtitle="Manage users and their roles on the system."
                actions={
                    <Button onClick={() => setShowForm(true)}>
                        + New User
                    </Button>
                }
            />

            <div class="mb-4 flex items-center gap-2">
                <input
                    type="text"
                    placeholder="Search users by name or username"
                    value={search()}
                    onInput={e => setSearch(e.target.value)}
                    onKeyDown={e => e.key === 'Enter' && handleSearch(e)}
                    class="
                        w-64 font-sans font-light text-sm
                        text-ink-dark placeholder:text-muted
                        bg-surface border border-border rounded-xs
                        px-3 py-2
                        focus:outline-none focus:border-ink
                    "
                />

                <Select
                    value={role()}
                    onChange={e => setRole(e.target.value ?? null)}
                >
                    <option value="">All Roles</option>
                    <For each={ROLES}>
                        {opt => <option value={opt.value}>{opt.label}</option>}
                    </For>
                </Select>
                <Button variant="ghost" size="md" onClick={handleSearch}>
                    Search
                </Button>
            </div>

            <UsrsTable
                usrs={usrs()}
                loading={loading()}
                onDelete={handleDelete}
            />

            <div class="mt-4">
                <Pagination
                    page={page()}
                    total={total()}
                    pageSize={20}
                    onChange={goTo}
                />
            </div>

            <CreateUsrForm
                open={showForm()}
                loading={creating()}
                onSubmit={handleCreate}
                onClose={() => setShowForm(false)}
            />
        </section>
    )
}
