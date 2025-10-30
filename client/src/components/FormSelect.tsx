import { useFormContext } from "react-hook-form"
import StatusMessage from "./StatusMessage"
import { memo } from "react"

const FormSelect = ({
    label,
    optionItems,
    name,
    required,
    disabled,
    warn,
    hideMsg,
}: {
    label?: string
    optionItems: { value: string; label: string; id: string }[]
    name: string
    required?: boolean
    disabled?: boolean
    warn?: boolean
    hideMsg?: boolean
}) => {
    // Form Context
    const {
        register,
        formState: { errors },
    } = useFormContext()
    const warningMsg = errors[name]?.message as string | undefined

    return (
        <div className="flex flex-col">
            {label && (
                <label className="text-2xl">
                    {label}
                    {required && <span className="text-prim-red">*</span>}
                </label>
            )}
            <select
                defaultValue=""
                {...register(name)}
                className="w-full h-11 pl-3 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
                style={warn ? { borderColor: "var(--color-prim-red)" } : {}}
                disabled={disabled}
            >
                <option value="" disabled={required} hidden={required}>
                    -- Select --
                </option>
                {optionItems.map((item) => (
                    <option value={item.value} key={item.id}>
                        {item.label}
                    </option>
                ))}
            </select>
            {!hideMsg && (
                <StatusMessage
                    msg={warningMsg}
                    variant="error"
                    className="text-xl"
                />
            )}
        </div>
    )
}

export default memo(FormSelect)
