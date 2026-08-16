import { Badge, Table } from '@/shared/ui'

const roleVariant = {
    super_usr_role: 'dark',
    normal_usr_role: 'accent',
    limited_usr_role: 'neutral'
}

const roleLabel = {
    super_usr_role: 'Super',
    normal_usr_role: 'Normal',
    limited_usr_role: 'Limited'
}
/**
 * @param {Object} props
 * @param {import('@/entities/usrs').Usr[]} props.usrs
 * @param {(usr: import('@/entities/usrs').Usr) => void} [props.onUpdate]
 * @param {(usr: import('@/entities/usrs').Usr) => void} [props.onDelete]
 * @param {boolean} [props.loading]
 */
export function UsrsTable(props) {
    const columns = [
        {
            key: 'username',
            header: 'Email',
            class: 'w-fit',
            render: row => (
                <span class="font-medium text-ink-dark">{row.username}</span>
            )
        },
        {
            key: 'firstname',
            header: 'First Name',
            class: 'w-fit',
            render: row => row.firstname ?? '-'
        },
        {
            key: 'lastname',
            header: 'Last Name',
            class: 'w-fit',
            render: row => row.lastname ?? '-'
        },
        {
            key: 'role',
            header: 'Role',
            class: 'w-fit',
            render: row => (
                <Badge variant={roleVariant[row.role] ?? 'neutral'}>
                    {roleLabel[row.role] ?? row.role}
                </Badge>
            )
        },
        {
            key: 'created_at',
            header: 'Created',
            class: 'w-fit',
            render: row => new Date(row.created_at).toLocaleDateString()
        },
        {
            key: 'actions',
            header: '',
            class: 'w-fit',
            align: 'right',
            render: row => (
                <div class="flex gap-2 items-center justify-end">
                    {props.onUpdate && (
                        <button
                            type="button"
                            onClick={e => {
                                e.stopPropagation()
                                props.onUpdate(row)
                            }}
                            class="text-sm text-accent hover:underline cursor-pointer font-bold"
                        >
                            Update
                        </button>
                    )}
                    {props.onDelete && (
                        <button
                            type="button"
                            onClick={e => {
                                e.stopPropagation()
                                props.onDelete(row)
                            }}
                            class="text-sm text-danger hover:underline cursor-pointer font-bold"
                        >
                            Delete
                        </button>
                    )}
                </div>
            )
        }
    ]

    return (
        <Table
            columns={columns}
            rows={props.usrs}
            loading={props.loading}
            emptyTitle="No users yet"
            emptyDescription="Create a user to get started."
        />
    )
}
