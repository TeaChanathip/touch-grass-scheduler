type Variant = "info" | "success" | "error"

const variantColorMap = new Map<Variant, string>([
    ["info", "var(--color-prim-gray-200)"],
    ["success", "var(--color-prim-green-600)"],
    ["error", "var(--color-prim-red)"],
])

export default function ResponseMessage({
    msg,
    variant,
}: {
    msg: string | undefined
    variant: Variant
}) {
    return msg ? (
        <p
            className="text-center text-xl"
            style={{ color: variantColorMap.get(variant) }}
        >
            {msg}
        </p>
    ) : (
        <p className="text-xl">&nbsp;</p>
    )
}
