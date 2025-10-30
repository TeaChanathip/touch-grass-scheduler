"use client"

import { VisibilityOffRounded, VisibilityRounded } from "@mui/icons-material"
import { CSSProperties, InputHTMLAttributes, useState } from "react"
import { UseFormRegisterReturn } from "react-hook-form"
import StatusMessage from "./StatusMessage"

interface FormStringInputProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string
    type: "number" | "text" | "email" | "password" | "tel" | "search" | "url" // Overwrite
    register?: UseFormRegisterReturn<any>
    warn?: boolean
    warningMsg?: string
    hideMsg?: boolean
}

export default function FormStringInput({
    label,
    placeholder,
    type,
    register,
    warn,
    warningMsg,
    readOnly,
    required,
    hideMsg,
    ...restProps
}: FormStringInputProps) {
    // Hooks
    const [isShowPassword, setShowPassword] = useState(false)

    const toggleShowPassword = () => {
        setShowPassword(!isShowPassword)
    }

    const constructAdditionalStyle = () => {
        const style: CSSProperties = {}
        if (warn || warningMsg) {
            style.borderColor = "var(--color-prim-red)"
        }
        if (type === "password") {
            style.paddingRight = "40px"
        }
        return style
    }

    return (
        <div className="flex flex-col">
            {label && (
                <label className="text-2xl">
                    {label}
                    {required && <span className="text-prim-red">*</span>}
                </label>
            )}
            <span className="flex flex-row items-center justify-end">
                <input
                    placeholder={placeholder}
                    type={
                        type === "password"
                            ? isShowPassword
                                ? "text"
                                : "password"
                            : type
                    }
                    {...register}
                    readOnly={readOnly}
                    required={required}
                    {...restProps}
                    className="w-full h-11 px-3 text-xl bg-white
                    border-prim-green-600 border-solid border-2 rounded-xl"
                    style={constructAdditionalStyle()}
                />
                {type === "password" && (
                    <span
                        className="absolute mr-3 cursor-pointer"
                        onClick={() => toggleShowPassword()}
                    >
                        {isShowPassword ? (
                            <VisibilityOffRounded />
                        ) : (
                            <VisibilityRounded />
                        )}
                    </span>
                )}
            </span>
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
