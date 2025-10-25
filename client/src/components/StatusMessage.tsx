import React from "react"

type Variant = "info" | "success" | "error"

const variantColorMap = new Map<Variant, string>([
    ["info", "var(--color-prim-gray-200)"],
    ["success", "var(--color-prim-green-600)"],
    ["error", "var(--color-prim-red)"],
])

interface StatusMessageProps extends React.HTMLProps<HTMLParagraphElement> {
    msg?: string
    variant: Variant
}

export default function StatusMessage(props: StatusMessageProps) {
    const { msg, variant, className, ...restProps } = props

    return (
        <p
            {...restProps}
            className={`text-center text-xl ${className}`}
            style={{
                color: variantColorMap.get(variant),
                visibility: msg ? "visible" : "hidden",
            }}
        >
            {msg ?? <>&nbsp;</>}
        </p>
    )
}
