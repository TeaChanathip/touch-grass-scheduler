import { UseFormRegisterReturn } from "react-hook-form"

export default function FormRadioGroup({
    label,
    options,
    register,
    warningMsg,
}: {
    label?: string
    options: { value: string; label: string }[]
    required?: boolean
    register?: UseFormRegisterReturn<any>
    warningMsg?: string
}) {
    // Assign id to options
    const optionItems = options.map((option) => {
        return { ...option, id: crypto.randomUUID() }
    })

    return (
        <div className="flex flex-col">
            {label && (
                <label className="text-2xl">
                    {label}
                    <span className="text-prim-red">*</span>
                </label>
            )}
            <div
                className="w-full px-10 py-5
                border-prim-green-600 border-solid border-2 rounded-xl 
                flex flex-wrap gap-5 justify-between text-xl"
                style={
                    warningMsg ? { borderColor: "var(--color-prim-red)" } : {}
                }
            >
                {optionItems.map((item) => (
                    <span
                        key={item.id}
                        className="flex flex-row items-center gap-2"
                    >
                        <input
                            id={`radio-${label}-${item.value}`}
                            type="radio"
                            value={item.value}
                            {...register}
                            className="size-6"
                        />
                        <label htmlFor={`radio-${label}-${item.value}`}>
                            {item.label}
                        </label>
                    </span>
                ))}
            </div>
            {register && (
                <p className="self-center text-prim-red">{warningMsg}&nbsp;</p>
            )}
        </div>
    )
}
