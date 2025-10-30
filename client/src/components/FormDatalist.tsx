import { memo, useId } from "react"
import { useFormContext } from "react-hook-form"
import StatusMessage from "./StatusMessage"

const FormDatalist = ({
    label,
    name,
    optionItems,
    readOnly,
    required,
    warn,
    hideMsg,
}: {
    label?: string
    name: string
    optionItems: { value: string; label: string; id: string }[]
    readOnly?: boolean
    required?: boolean
    warn?: boolean
    hideMsg?: boolean
}) => {
    // Generate ID for referencing datalist
    const datalistID = useId()

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
            <input
                type="text"
                list={datalistID}
                readOnly={readOnly}
                required={required}
                {...register(name)}
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

export default memo(FormDatalist)
