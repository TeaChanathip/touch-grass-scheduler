import { useId } from "react"
import { UseFormRegisterReturn } from "react-hook-form"

export default function FormDatalist({
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
            {warningMsg !== undefined && (
                <p className="self-center text-prim-red">{warningMsg}&nbsp;</p>
            )}
        </div>
    )
}
