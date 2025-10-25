import { UseFormRegisterReturn } from "react-hook-form"
import StatusMessage from "./StatusMessage"

export default function FormSelect({
    label,
    optionItems,
    required,
    register,
    warn,
    warningMsg,
}: {
    label?: string
    optionItems: { value: string; label: string; id: string }[]
    required?: boolean
    register?: UseFormRegisterReturn<any>
    warn?: boolean
    warningMsg?: string
}) {
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
                {...register}
                className="w-full h-11 pl-3 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
                style={warn ? { borderColor: "var(--color-prim-red)" } : {}}
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
            <StatusMessage
                msg={warningMsg}
                variant="error"
                className="text-xl"
            />
        </div>
    )
}
