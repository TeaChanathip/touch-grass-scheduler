import { UseFormRegisterReturn } from "react-hook-form"

export default function FormSelect({
    label,
    options,
    required,
    register,
    warningMsg,
}: {
    label?: string
    options: { value: string; label: string }[]
    required?: boolean
    register?: UseFormRegisterReturn<any>
    warningMsg?: string
}) {
    const optionItems = options.map((option) => {
        return { ...option, id: crypto.randomUUID() }
    })

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
                className="w-full h-11 pl-3 pr-8 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
                style={
                    warningMsg ? { borderColor: "var(--color-prim-red)" } : {}
                }
            >
                <option value="" disabled={required} hidden={required}>
                    Select an option
                </option>
                {optionItems.map((item) => (
                    <option value={item.value} key={item.id}>
                        {item.label}
                    </option>
                ))}
            </select>
            {register && (
                <p className="self-center text-prim-red">{warningMsg}&nbsp;</p>
            )}
        </div>
    )
}
