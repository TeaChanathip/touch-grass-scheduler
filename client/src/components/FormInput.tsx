"use client"

import { ZodType } from "zod"
import ClearRoundedIcon from "@mui/icons-material/ClearRounded"
import { ChangeEvent, useState } from "react"

export default function FormStringInput(props: {
    label?: string
    name: string
    placeholder?: string
    required?: boolean
    type: "number" | "text" | "email" | "password" | "tel" | "search" | "url"
    schema?: ZodType
}) {
    const { label, name, placeholder, type, required, schema } = props
    const [inputValue, setInputValue] = useState("")
    const [warningMsg, setWarningMsg] = useState("")

    const inputChangeHandler = (e: ChangeEvent<HTMLInputElement>) => {
        setInputValue(e.target.value)

        if (schema === undefined) return

        // Perform input validation
        const result = schema.safeParse(inputValue)
        if (result.success) {
            setWarningMsg("")
        } else {
            setWarningMsg(result.error.issues[0].message)
        }
    }

    return (
        <div className="flex flex-col">
            <label className="text-2xl">
                {label}
                {required && <span className="text-prim-red">*</span>}
            </label>
            <span className="flex items-center justify-end">
                <input
                    name={name}
                    placeholder={placeholder}
                    type={type}
                    value={inputValue}
                    required={required}
                    onChange={inputChangeHandler}
                    className="w-full h-11 pl-3 pr-8 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
                />
                <ClearRoundedIcon
                    className="absolute mr-2 hover:text-prim-gray-300 cursor-pointer"
                    onClick={() => setInputValue("")}
                />
            </span>
            {schema && (
                <p className="self-center text-prim-red">
                    {inputValue && warningMsg}&nbsp;
                </p>
            )}
        </div>
    )
}

// TODO: Remove validator, and require warning message from parent instead
