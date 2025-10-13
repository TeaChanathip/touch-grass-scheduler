"use client"

import { UseFormRegisterReturn } from "react-hook-form"

export default function FormStringInput(props: {
    label?: string
    placeholder?: string
    required?: boolean
    type: "number" | "text" | "email" | "password" | "tel" | "search" | "url"
    register?: UseFormRegisterReturn<any>
    warningMsg?: string
}) {
    const { label, placeholder, type, required, register, warningMsg } = props

    return (
        <div className="flex flex-col">
            <label className="text-2xl">
                {label}
                {required && <span className="text-prim-red">*</span>}
            </label>
            <input
                placeholder={placeholder}
                type={type}
                required={required}
                {...register}
                className="w-full h-11 pl-3 pr-8 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
            />
            {register && (
                <p className="self-center text-prim-red">{warningMsg}&nbsp;</p>
            )}
        </div>
    )
}
