import { useId } from "react"
import { UseFormRegisterReturn } from "react-hook-form"
import StatusMessage from "./StatusMessage"

export default function FormDatalist({
    label,
    optionItems,
    readOnly,
    required,
    register,
    warn,
    warningMsg,
    hideMsg,
}: {
    label?: string
    optionItems: { value: string; label: string; id: string }[]
    readOnly?: boolean
    required?: boolean
    register?: UseFormRegisterReturn<any>
    warn?: boolean
    warningMsg?: string
    hideMsg?: boolean
}) {
    const datalistID = useId()

    return (
        <div className="flex flex-col">
            {label && (
                <label className="text-2xl">
                    {label}
                    {required && <span className="text-prim-red">*</span>}
                </label>
            )}
            <input
                type="text"
                list={datalistID}
                readOnly={readOnly}
                required={required}
                {...register}
                className="w-full h-11 pl-3 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
                style={warn ? { borderColor: "var(--color-prim-red)" } : {}}
            />
            <datalist id={datalistID}>
                {optionItems.map((item) => (
                    <option value={item.value} key={item.id}>
                        {item.label}
                    </option>
                ))}
            </datalist>
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
