/**
 * Surface container with sub-components: `Card`, `Card.Header`, `Card.Body`, `Card.Footer`.
 *
 * @example
 * <Card onClick={handleClick}>
 *   <Card.Header title="Document" action={<Button size="sm">New</Button>} />
 *   <Card.Body>...</Card.Body>
 *   <Card.Footer>...</Card.Footer>
 * </Card>
 */

import { splitProps } from 'solid-js'

/**
 * @param {Object} props
 * @param {boolean} [props.hoverable=false] - Adds subtle hover effect for clickable cards.
 * @param {string} [props.class]
 * @param {import('solid-js').JSX.HTMLAttributes<HTMLDivElement>} props
 * @returns {import('solid-js').JSX.Element}
 */
export function Card(props) {
    const [local, rest] = splitProps(props, ['hoverable', 'class', 'children'])
    return (
        <div
            class={`
                bg-bg border border-border-subtle
                rounease-basition-colors ease-base
                ${local.hoverable ? 'hover:border-accent  hover:bg-bg-2 cursor-pointer' : ''}
                ${local.class}
            `}
            {...rest}
        >
            {local.children}
        </div>
    )
}

Card.Header = function CardHeader(props) {
    const [local, rest] = splitProps(props, ["title", "subtitle", "action", "class"])

    return (
        <div
            class={`
                flex items-start justify-between gap-4
                px-4 py-3 border-b border-border-subtle
                ${local.class}
            `}
            {...rest}
        >
            <div class="flex flex-col gap-0.5">
                {local.title && (
                    <h3 class="font-serif text-base text-ink-dark leading-snug">
                        {local.title}
                    </h3>
                )}
                {local.subtitle && <p class="text-xs text-muted">{local.subtitle}</p>}
            </div>
            {local.action && <div class="shrink-0">{local.action}</div>}
        </div>
    )
}

Card.Body = function CardBody(props) {
    const [local, rest] = splitProps(props, ["class", "children"])

    return (
        <div class={`px-4 py-3 ${local.class}`}  {...rest}>
            {local.children}
        </div>
    )
}

Card.Footer = function CardFooter(props) {
    const [local, rest] = splitProps(props, ["class", "children"])


    return (
        <div
            class={`
                flex items-center justify-between
                px-4 py-3 border-t border-border-subtle
                ${local.class}
            `}
            {...rest}
        >
            {local.children}
        </div>
    )
}
